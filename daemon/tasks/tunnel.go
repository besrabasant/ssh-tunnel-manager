package tasks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultSSHPort = "22"
)

type ConnectionInfo struct {
	Listen *net.Listener
	Conns  []net.Conn
	Cancel func()
}

func (c *ConnectionInfo) AddConnection(conn net.Conn) {
	c.Conns = append(c.Conns, conn)
}

func (c *ConnectionInfo) StopListening() {
	if c.Listen != nil {
		connListener := *c.Listen
		if err := connListener.Close(); err != nil {
			fmt.Println("Error closing listener:", err)
		}
	}
}

func (c *ConnectionInfo) KillAllConns() {
	if len(c.Conns) > 0 {
		for _, conn := range c.Conns {
			conn.Close()
		}
	}
}

var (
	// Map to keep track of connections by their local port
	connections = make(map[int]ConnectionInfo)
	// Mutex to protect access to the connections map
	connMutex = &sync.Mutex{}

	shutdown = make(chan struct{}) // Signal for shutdown
)

func StartTunnelTask(ctx context.Context, req *pb.StartTunnelRequest) (*pb.StartTunnelResponse, error) {
	var output strings.Builder

	dirpath := configmanager.DefaultConfigDir
	if value := os.Getenv(add.ConfigDirFlagName); value != "" {
		dirpath = value
	}

	configdir, err := utils.ResolveDir(dirpath)
	if err != nil {
		return nil, err
	}

	cfg, err := configmanager.NewManager(configdir).GetConfiguration(req.ConfigName)
	if err != nil {
		return nil, fmt.Errorf("couldn't get configuration %q: %v", req.ConfigName, err)
	}

	resultChan, errorChan := startTunneling(ctx, cfg, int(req.LocalPort))
	var errReceived bool = false

loop:
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				if !errReceived {
					output.WriteString("Tunnel setup completed.\n")
				}
				break loop
			}
			output.WriteString(result + "\n")

		case err, ok := <-errorChan:
			if ok && err != nil { // Check if error is not nil
				output.WriteString(fmt.Sprintf("Failed to start tunneling: %v\n", err))
				errReceived = true
			}
		}
	}

	return &pb.StartTunnelResponse{Result: output.String()}, nil
}

func startTunneling(ctx context.Context, entry configmanager.Entry, localPort int) (<-chan string, <-chan error) {
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		connMutex.Lock()
		if connInfo, exists := connections[localPort]; exists {
			resultChan <- fmt.Sprint("Closing existing connection on port ", localPort, "\n")

			// If there's an existing connection on the same port, close it
			connInfo.Cancel() // Cancel the context of the existing connection
			if connInfo.Listen != nil {
				connInfo.StopListening()
				connInfo.KillAllConns()
			}
		}

		// Initial status message
		resultChan <- "Starting tunnel setup...\n"

		// Create new cancellation context for this connection
		tunnelCtx, cancel := context.WithCancel(ctx)
		connections[localPort] = ConnectionInfo{
			Listen: nil, // This will be updated once the listener is set up
			Cancel: cancel,
		}
		connMutex.Unlock()

		// The SSH server to connect to. The address can contain a port.
		sshServer := entry.Server
		// The username to use when connecting
		sshUser := entry.User
		// The private key file to use for authentication
		keyFile := entry.KeyFile
		// The remote host and port to forward traffic to
		remoteAddress := fmt.Sprintf("%s:%d", entry.RemoteHost, entry.RemotePort)
		localAddress := fmt.Sprintf("%s:%d", "127.0.0.1", localPort)

		// Check if the ssh server address specifies a port. And use 22 if not.
		_, _, err := net.SplitHostPort(sshServer)
		if err != nil {
			var addrErr *net.AddrError
			if errors.As(err, &addrErr) {
				hasPort := strings.LastIndex(sshServer, ":") != -1
				if hasPort {
					errorChan <- fmt.Errorf("bad ssh server address: %v", err)
					return
				} else {
					resultChan <- fmt.Sprintf("SSH server %q specifies no port. Will use %s\n", sshServer, DefaultSSHPort)
					// Use 22 as a default ssh port.
					sshServer = sshServer + ":" + DefaultSSHPort
				}
			} else {
				errorChan <- err
			}
		}

		// Load the private key file
		privateKey, err := readPrivateKeyFile(keyFile)
		if err != nil {
			errorChan <- err
		}

		key, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			errorChan <- fmt.Errorf("couldn't parse private key %q: %v", keyFile, err)
		}

		// Set up the SSH client config
		timeout := 10 * time.Second
		config := &ssh.ClientConfig{
			User: sshUser,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         timeout,
		}

		// Connect to the SSH server
		resultChan <- fmt.Sprintf("Connecting to %q with a timeout of %s", sshServer, timeout)
		client, err := ssh.Dial("tcp", sshServer, config)
		if err != nil {
			errorChan <- fmt.Errorf("couldn't connect to SSH server %q: %v", sshServer, err)
		}
		defer client.Close()
		resultChan <- "\nConnected\n"

		// Set up the local listener
		var localListener net.Listener
		var listenerErr error
		for attempts := 0; attempts < 5; attempts++ {
			localListener, listenerErr = net.Listen("tcp", localAddress)
			if listenerErr == nil {
				break // Successfully bound to the port
			}
			time.Sleep(time.Second) // Wait before retrying
		}
		if listenerErr != nil {
			errorChan <- fmt.Errorf("couldn't set up local listener after retries: %v", err)
			return
		}
		defer localListener.Close()

		// Update the connections map with the actual listener
		connMutex.Lock()
		currentConnInfo := connections[localPort]
		currentConnInfo.Listen = &localListener
		connections[localPort] = currentConnInfo
		connMutex.Unlock()

		resultChan <- "Local listener set up, ready to accept connections.\n"

		// Handle server shutdown triggered by signal or context cancellation
		go func() {
			select {
			case <-tunnelCtx.Done():
				if err := localListener.Close(); err != nil {
					fmt.Printf("Error closing local listener: %v\n", err)
				}
				if err := client.Close(); err != nil {
					fmt.Printf("Error closing SSH client: %v\n", err)
				}
			case <-ctx.Done():
				// Additional cleanup if needed
			}
		}()

		// Start accepting connections on the local listener
		resultChan <- fmt.Sprintf("Tunneling %q <==> %q through %q\n", localAddress, remoteAddress, sshServer)
		close(resultChan)
		close(errorChan)
		for {
			select {
			case <-tunnelCtx.Done():
			case <- shutdown:
				// If we're done, ensure we close the listener if it's not nil.
				if localListener != nil {
					connMutex.Lock()
					currentConnInfo := connections[localPort]
					currentConnInfo.StopListening()
					connMutex.Unlock()
				}
			default:
				if localListener == nil {
					fmt.Println("Unexpected error: localListener is nil")
					return
				}
				localConn, err := localListener.Accept()
				if err != nil && localConn != nil {
					fmt.Printf("Error accepting new connection: %v\n", err)
					<-shutdown
				}

				if localConn != nil {

					// Update the connections map with the actual listener
					connMutex.Lock()
					currentConnInfo := connections[localPort]
					currentConnInfo.AddConnection(localConn)
					connections[localPort] = currentConnInfo
					connMutex.Unlock()

					fmt.Printf("Got new connection: %s\n", localConn.RemoteAddr())
					// Handle the connection...

					// Start the SSH tunnel for each incoming connection
					go func(localConn net.Conn) {
						remoteConn, err := client.Dial("tcp", remoteAddress)
						if err != nil {
							fmt.Printf("error dialing remote address %s: %v", remoteAddress, err)
							localConn.Close()
							return
						}

						runTunnel(localConn, remoteConn)
					}(localConn)
				}
			}
		}
	}()

	return resultChan, errorChan
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
	fmt.Printf("\nConnection closed: %s\n\n", local.RemoteAddr())
}
