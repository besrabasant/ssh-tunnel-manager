package tasks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultSSHPort = "22"
)

// TunnelManager manages multiple SSH tunnels
type TunnelManager struct {
	mutex       sync.Mutex
	connections SSHConnections
	shutdown    chan struct{}
	resultChan  chan string
	errorChan   chan error
}

func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		connections: make(SSHConnections),
		mutex:       sync.Mutex{},
		shutdown:    make(chan struct{}),
	}
}

func (m *TunnelManager) createResultChannels() {
	m.resultChan = make(chan string, 1)
	m.errorChan = make(chan error, 1)
}

func (m *TunnelManager) startTunneling(ctx context.Context, entry configmanager.Entry, localPort int) {
	// Defer closing of communication channels
	defer close(m.resultChan)
	defer close(m.errorChan)

	// Enable lock
	m.mutex.Lock()

	// Close Existing connections
	if connInfo, exists := m.connections[localPort]; exists {
		m.resultChan <- fmt.Sprint("Closing existing connection on port ", localPort, "\n")

		// If there's an existing connection on the same port, close it
		connInfo.Cancel() // Cancel the context of the existing connection
		if connInfo.Listener != nil {
			connInfo.StopListeners()
			connInfo.KillAllConnections()
		}
	}

	// Initial status message
	m.resultChan <- "Starting tunnel setup...\n"

	// The SSH server to connect to. The address can contain a port.
	sshServer := entry.Server
	// The username to use when connecting
	sshUser := entry.User
	// The private key file to use for authentication
	keyFile := entry.KeyFile
	// The remote host and port to forward traffic to
	remoteAddress := fmt.Sprintf("%s:%d", entry.RemoteHost, entry.RemotePort)
	localAddress := fmt.Sprintf("%s:%d", "127.0.0.1", localPort)

	// Create new cancellation context for this connection
	tunnelCtx, cancel := context.WithCancel(ctx)
	m.connections[localPort] = ConnectionInfo{
		Listener:   nil, // This will be updated once the listener is set up
		Cancel:     cancel,
		RemoteAddr: remoteAddress,
		LocalAddr:  localAddress,
		Config:     entry,
	}

	// Check if the ssh server address specifies a port. And use 22 if not.
	_, _, serverReadErr := net.SplitHostPort(sshServer)
	if serverReadErr != nil {
		var addrErr *net.AddrError
		if errors.As(serverReadErr, &addrErr) {
			hasPort := strings.LastIndex(sshServer, ":") != -1
			if hasPort {
				m.errorChan <- fmt.Errorf("bad ssh server address: %v", serverReadErr)
				return
			} else {
				m.resultChan <- fmt.Sprintf("SSH server %q specifies no port. Will use %s\n", sshServer, DefaultSSHPort)
				// Use 22 as a default ssh port.
				sshServer = sshServer + ":" + DefaultSSHPort
			}
		} else {
			m.errorChan <- serverReadErr
		}
	}

	// Load the private key file
	privateKey, keyReadErr := readPrivateKeyFile(keyFile)
	if keyReadErr != nil {
		m.errorChan <- keyReadErr
	}

	key, keyParseErr := ssh.ParsePrivateKey(privateKey)
	if keyParseErr != nil {
		m.errorChan <- fmt.Errorf("couldn't parse private key %q: %v", keyFile, keyParseErr)
		return
	}

	// Define timeout for SSH connection
	timeout := 10 * time.Second

	// Set up the SSH client config
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	// Connect to the SSH server
	m.resultChan <- fmt.Sprintf("Connecting to %q with a timeout of %s", sshServer, timeout)

	client, clientErr := ssh.Dial("tcp", sshServer, config)
	if clientErr != nil {
		m.errorChan <- fmt.Errorf("couldn't connect to SSH server %q: %v", sshServer, clientErr)
	}

	m.resultChan <- "\nConnected\n"

	// Set up the local listener
	var listener net.Listener
	var listenerErr error

	for attempts := 0; attempts < 5; attempts++ {
		// Forward the local port to the remote address
		listener, listenerErr = net.Listen("tcp", localAddress)
		if listenerErr == nil {
			break // Successfully bound to the port
		}
		time.Sleep(time.Second) // Wait before retrying
	}

	if listenerErr != nil {
		m.errorChan <- fmt.Errorf("couldn't set up local listener after retries: %v", listenerErr)
		return
	}

	// Update the connections map with the actual listener
	currentConnInfo := m.connections[localPort]
	currentConnInfo.Client = client
	currentConnInfo.Listener = listener
	m.connections[localPort] = currentConnInfo

	m.mutex.Unlock()

	m.resultChan <- "Local listener set up, ready to accept connections.\n"

	// Handle server shutdown triggered by signal or context cancellation
	go m.handleServerShutdown(ctx, tunnelCtx, localPort)

	// Start accepting connections on the local listener
	m.resultChan <- fmt.Sprintf("Tunneling %q <==> %q through %q\n", localAddress, remoteAddress, sshServer)

	// Handle incoming connections on local port
	go m.forwardTunnel(tunnelCtx, remoteAddress, localPort)
}

func (m *TunnelManager) handleServerShutdown(ctx context.Context, tunnelCtx context.Context, localPort int) {
	currentConnInfo := m.connections[localPort]

	select {
	case <-tunnelCtx.Done():
	case <-m.shutdown:
		// If we're done, ensure we close the listener if it's not nil.
		if currentConnInfo.Client != nil {
			m.mutex.Lock()
			currentConnInfo.StopListeners()
			currentConnInfo.KillAllConnections()
			currentConnInfo.StopClient()
			delete(m.connections, localPort)
			m.mutex.Unlock()
		}
	case <-ctx.Done():
		// Additional cleanup if needed
	}
}

// forwardTunnel handles incoming connections on the local port and forwards them to the remote server
func (m *TunnelManager) forwardTunnel(tunnelCtx context.Context, remoteAddress string, localPort int) {
	currentConnInfo := m.connections[localPort]

	for {
		if currentConnInfo.Listener == nil {
			log.Printf("Unexpected error: localListener is nil")
			return
		}
		localConn, err := currentConnInfo.Listener.Accept()
		if err != nil {
			log.Printf("Failed to accept local connection: %v", err)
			continue
		}

		m.mutex.Lock()
		currentConnInfo := m.connections[localPort]
		currentConnInfo.AddConnection(localConn)
		m.connections[localPort] = currentConnInfo
		m.mutex.Unlock()

		// Start the SSH tunnel for each incoming connection
		go func(localConn net.Conn) {
			remoteConn, err := currentConnInfo.Client.Dial("tcp", remoteAddress)
			if err != nil {
				log.Printf("error dialing remote address %s: %v", remoteAddress, err)
				localConn.Close()
				return
			}

			runTunnel(localConn, remoteConn)
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
	log.Printf("\nConnection closed: %s\n\n", local.RemoteAddr())
}
