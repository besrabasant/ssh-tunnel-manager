package tunnelmanager

import (
	"fmt"
	"net"
	"sync"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"golang.org/x/crypto/ssh"
)

// SSH Connection info
type ConnectionInfo struct {
	Client      *ssh.Client         // Pointer to an ssh.Client which manages the SSH connection
	Listener    net.Listener        // Network listener for the local side of the SSH tunnel
	LocalAddr   string              // Local address where the tunnel's listener is bound
	RemoteAddr  string              // Remote address to which the tunnel connects
	Connections []net.Conn          // Slice of net.Conn representing active connections over this tunnel
	Config      configmanager.Entry // Configuration entry, assuming this is a struct from the 'configmanager' package
	Cancel      func()              // Function to call to cleanly shutdown or cancel the tunnel/connection
	clientMu    sync.RWMutex
	reconnectMu sync.Mutex
}

func (c *ConnectionInfo) getClient() *ssh.Client {
	c.clientMu.RLock()
	defer c.clientMu.RUnlock()
	return c.Client
}

// replaceClient swaps the SSH client without racing with tunnel goroutines
// that are using the current client. The old client is closed after the swap
// so connections that survived the network transition cannot linger.
func (c *ConnectionInfo) replaceClient(client *ssh.Client) {
	c.clientMu.Lock()
	old := c.Client
	c.Client = client
	c.clientMu.Unlock()

	if old != nil && old != client {
		_ = old.Close()
	}
}

// Add new connection
func (c *ConnectionInfo) AddConnection(conn net.Conn) {
	c.Connections = append(c.Connections, conn)
}

func (c *ConnectionInfo) StopClient() {
	c.clientMu.Lock()
	client := c.Client
	c.Client = nil
	c.clientMu.Unlock()

	if client != nil {
		if err := client.Close(); err != nil {
			fmt.Printf("Error closing SSH client: %v\n", err)
		}
	}
}

func (c *ConnectionInfo) StopListeners() {
	if c.Listener != nil {
		connListener := c.Listener
		if err := connListener.Close(); err != nil {
			fmt.Println("Error closing listener:", err)
		}
	}
}

func (c *ConnectionInfo) KillAllConnections() {
	for _, conn := range c.Connections {
		if err := conn.Close(); err != nil {
			fmt.Println("Failed to close connection:", err)
		}
	}
	c.Connections = nil // Reset the slice to clear references and help with garbage collection
}

// Method to clear all connections safely and effectively
func (c *ConnectionInfo) ClearConnection() {
	if c.Cancel != nil {
		c.Cancel()
	}
	c.StopClient()         // Ensures all connections are closed first
	c.StopListeners()      // Ensures all connections are closed first
	c.KillAllConnections() // Ensures all connections are closed first
	c.Connections = nil    // Clears the slice completely to release resources
}
