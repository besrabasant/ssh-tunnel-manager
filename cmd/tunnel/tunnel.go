package tunnel

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultSSHPort = "22"
)

var Cmd cli.Command = cli.Command{

	Name:      "tunnel",
	Usage:     "Start a tunnel using a configuration. The tunnel will forward connections to [local port] if specified or to a random port.",
	UsageText: "ssh-tunnel-manager tunnel <configuration name> [local port]",
	ArgsUsage: "<configuration name>",
	Action: func(cCtx *cli.Context) error {

		entryName := cCtx.Args().First()
		if entryName == "" {
			return fmt.Errorf("<configuration name> needed but not provided")
		}

		localPortStr := cCtx.Args().Get(1)
		var localPort int
		if localPortStr == "" {
			// Generate random port
			randomPort, err := generateRandomPort()
			if err != nil {
				return fmt.Errorf("couldn't generate a random port: %v", err)
			}
			localPort = randomPort
		} else {
			// Parse given port
			localPortInt, err := strconv.Atoi(localPortStr)
			if err != nil {
				return fmt.Errorf("provided local port %q is not a valid port", localPortStr)
			}
			localPort = localPortInt
		}

		configdir, err := utils.ResolveDir(cCtx.String(add.ConfigDirFlagName))
		if err != nil {
			return err
		}

		cfg, err := configmanager.NewManager(configdir).GetConfiguration(entryName)
		if err != nil {
			return fmt.Errorf("couldn't get configuration %q: %v", entryName, err)
		}

		return startTunneling(cfg, localPort)
	},
}

func startTunneling(entry configmanager.Entry, localPort int) error {
	// The SSH server to connect to. The address can contain a port.
	sshServer := entry.Server
	// The username to use when connecting
	sshUser := entry.User
	// The private key file to use for authentication
	keyFile := entry.KeyFile
	// The remote host and port to forward traffic to
	remoteAddress := fmt.Sprintf("%s:%d", entry.RemoteHost, entry.RemotePort)
	localAddress := fmt.Sprintf("%s:%d", "localhost", localPort)

	// Check if the ssh server address specifies a port. And use 22 if not.
	_, _, err := net.SplitHostPort(sshServer)
	if err != nil {
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) {
			hasPort := strings.LastIndex(sshServer, ":") != -1
			if hasPort {
				return fmt.Errorf("bad ssh server address: %v", err)
			} else {
				log.Printf("SSH server %q specifies no port. Will use %s\n", sshServer, DefaultSSHPort)
				// Use 22 as a default ssh port.
				sshServer = sshServer + ":" + DefaultSSHPort
			}
		} else {
			return err
		}
	}

	// Load the private key file
	privateKey, err := readPrivateKeyFile(keyFile)
	if err != nil {
		return err
	}

	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("couldn't parse private key %q: %v", keyFile, err)
	}

	// Set up the SSH client config
	timeout := 8 * time.Second
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	// Connect to the SSH server
	log.Printf("Connecting to %q with a timeout of %s", sshServer, timeout)
	client, err := ssh.Dial("tcp", sshServer, config)
	if err != nil {
		return fmt.Errorf("couldn't connect to SSH server %q: %v", sshServer, err)
	}
	log.Printf("Connected\n")

	// Set up the local listener
	localListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", localPort))
	if err != nil {
		return fmt.Errorf("couldn't set up local listener: %v", err)
	}

	// Handle server shutdown
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, syscall.SIGTERM, os.Interrupt)
	go func(stopper <-chan os.Signal) {
		<-stopper
		log.Printf("Stopping the tunnel server. Closing any existing connections...\n")
		code := 0
		if err := client.Close(); err != nil {
			log.Printf("error while stopping the tunnel server: %v", err)
			code = 1
		}
		os.Exit(code)
	}(stopper)

	// Start accepting connections on the local listener
	log.Printf("Tunneling %q <==> %q through %q\n", localAddress, remoteAddress, sshServer)
	for {
		localConn, err := localListener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting new connection: %v", err)
		}

		log.Printf("Got new connection: %s\n", localConn.RemoteAddr())

		// Start the SSH tunnel for each incoming connection
		go func(localConn net.Conn) {
			remoteConn, err := client.Dial("tcp", remoteAddress)
			if err != nil {
				log.Printf("error dialing remote address %s: %v\n", remoteAddress, err)
				localConn.Close()
				return
			}

			go runTunnel(localConn, remoteConn)
		}(localConn)
	}
}

// Helper function to read the private key file
func readPrivateKeyFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read private key file %s: %v", file, err)
	}
	return data, nil
}

// runTunnel runs a tunnel between two connections; as soon as one connection
// reaches EOF or reports an error, both connections are closed and this
// function returns.
func runTunnel(local, remote net.Conn) {
	// Clean up
	defer local.Close()
	defer remote.Close()

	done := make(chan struct{}, 2)
	go func() {
		io.Copy(local, remote)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(remote, local)
		done <- struct{}{}
	}()

	<-done
	log.Printf("Connection closed: %s\n", local.RemoteAddr())
}

func generateRandomPort() (int, error) {
	// Listen on port 0 to bind to a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	// Extract the port number from the listener address
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, err
	}

	// Convert the port number to an integer
	randomPort, err := net.LookupPort("tcp", port)
	if err != nil {
		return 0, err
	}

	// Close the listener
	err = listener.Close()
	if err != nil {
		return 0, err
	}

	return randomPort, nil
}
