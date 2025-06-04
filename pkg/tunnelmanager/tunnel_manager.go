package tunnelmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"golang.org/x/crypto/ssh"
)

// TunnelManager manages multiple SSH tunnels
type TunnelManager struct {
	Mutex       sync.Mutex
	Connections SSHConnections
	shutdown    chan struct{}
	ResultChan  chan string
	ErrChan     chan error
}

func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		Connections: make(SSHConnections),
		Mutex:       sync.Mutex{},
		shutdown:    make(chan struct{}),
	}
}

func (m *TunnelManager) CreateResultChannels() {
	m.ResultChan = make(chan string, 1)
	m.ErrChan = make(chan error, 1)
}

func (m *TunnelManager) StartTunneling(ctx context.Context, entry configmanager.Entry, localPort int) {
	// Defer closing of communication channels
	defer close(m.ResultChan)
	defer close(m.ErrChan)

	// Enable lock and ensure it's released on every exit path
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	// Close Existing connections
	if connInfo, exists := m.Connections[localPort]; exists {
		m.ResultChan <- fmt.Sprint("Closing existing connection on port ", localPort, "\n")

		// If there's an existing connection on the same port, close it
		connInfo.Cancel() // Cancel the context of the existing connection
		if connInfo.Listener != nil {
			connInfo.StopListeners()
			connInfo.ClearConnection()
		}
	}

	// Initial status message
	m.ResultChan <- "Starting tunnel setup...\n"

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
	m.Connections[localPort] = &ConnectionInfo{
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
				m.ErrChan <- fmt.Errorf("bad ssh server address: %v", serverReadErr)
				return
			} else {
				m.ResultChan <- fmt.Sprintf("SSH server %q specifies no port. Will use %s\n", sshServer, config.DefaultSSHPort)
				// Use 22 as a default ssh port.
				sshServer = sshServer + ":" + config.DefaultSSHPort
			}
		} else {
			m.ErrChan <- serverReadErr
		}
	}

	// Load the private key file
	privateKey, keyReadErr := readPrivateKeyFile(keyFile)
	if keyReadErr != nil {
		m.ErrChan <- keyReadErr
		return
	}

	key, keyParseErr := ssh.ParsePrivateKey(privateKey)
	if keyParseErr != nil {
		m.ErrChan <- fmt.Errorf("couldn't parse private key %q: %v", keyFile, keyParseErr)
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
	m.ResultChan <- fmt.Sprintf("Connecting to %q with a timeout of %s", sshServer, timeout)

	client, clientErr := ssh.Dial("tcp", sshServer, config)
	if clientErr != nil {
		m.ErrChan <- fmt.Errorf("couldn't connect to SSH server %q: %v", sshServer, clientErr)
		return
	}

	m.ResultChan <- "\nConnected\n"

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
		m.ErrChan <- fmt.Errorf("couldn't set up local listener after retries: %v", listenerErr)
		return
	}

	// Update the connections map with the actual listener
	currentConnInfo := m.Connections[localPort]
	currentConnInfo.Client = client
	currentConnInfo.Listener = listener
	m.Connections[localPort] = currentConnInfo

	m.ResultChan <- "Local listener set up, ready to accept connections.\n"

	// Start accepting connections on the local listener
	m.ResultChan <- fmt.Sprintf("Tunneling %q <==> %q through %q\n", localAddress, remoteAddress, sshServer)

	// Handle incoming connections on local port
	go m.forwardTunnel(ctx, tunnelCtx, remoteAddress, localPort)
}

// forwardTunnel handles incoming connections on the local port and forwards them to the remote server
func (m *TunnelManager) forwardTunnel(ctx context.Context, tunnelCtx context.Context, remoteAddress string, localPort int) {
	m.Mutex.Lock()
	currentConnInfo, exists := m.Connections[localPort]
	m.Mutex.Unlock()

	if !exists || currentConnInfo.Listener == nil {
		log.Println("No listener found, returning")
		return
	}

	for {
		select {
		case <-m.shutdown:
			m.Mutex.Lock()
			for _, connInfo := range m.Connections {
				if connInfo.Client != nil {
					connInfo.ClearConnection()
				}
			}
			m.Connections = make(SSHConnections) // Clear all connections
			m.Mutex.Unlock()
			return

		case <-tunnelCtx.Done():
			m.Mutex.Lock()
			currentConnInfo.Cancel()
			delete(m.Connections, localPort)
			m.Mutex.Unlock()
			return // Exit the loop and goroutine
		case <-ctx.Done():
			// Additional cleanup if needed
			return // Usually signifies the parent context was cancelled
		default:
			m.Mutex.Lock()
			if currentConnInfo.Listener == nil {
				m.Mutex.Unlock()
				log.Printf("Listener is not available or connection info does not exist for port %d", localPort)
				return
			}
			m.Mutex.Unlock()

			localConn, err := currentConnInfo.Listener.Accept()
			if err != nil {
				// If we receive an error due to the listener being closed, gracefully exit the loop
				if errors.Is(err, net.ErrClosed) {
					break
				}
				log.Printf("Failed to accept local connection: %v", err)
				continue
			}

			m.Mutex.Lock()
			currentConnInfo.AddConnection(localConn)
			m.Mutex.Unlock()

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
}
