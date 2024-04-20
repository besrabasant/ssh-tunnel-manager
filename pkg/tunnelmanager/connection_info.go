package tunnelmanager

import (
	"fmt"
	"net"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"golang.org/x/crypto/ssh"
)

// SSH Connection info
type ConnectionInfo struct {
    Client      *ssh.Client     // Pointer to an ssh.Client which manages the SSH connection
    Listener    net.Listener    // Network listener for the local side of the SSH tunnel
    LocalAddr   string          // Local address where the tunnel's listener is bound
    RemoteAddr  string          // Remote address to which the tunnel connects
    Connections []net.Conn      // Slice of net.Conn representing active connections over this tunnel
    Config      configmanager.Entry // Configuration entry, assuming this is a struct from the 'configmanager' package
    Cancel      func()          // Function to call to cleanly shutdown or cancel the tunnel/connection
}

// Add new connection
func (c *ConnectionInfo) AddConnection(conn net.Conn) {
	c.Connections = append(c.Connections, conn)
}

func (c *ConnectionInfo) StopClient() {
	if c.Client != nil {
		if err := c.Client.Close(); err != nil {
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
	c.StopClient() // Ensures all connections are closed first
	c.StopListeners() // Ensures all connections are closed first
	c.KillAllConnections() // Ensures all connections are closed first
	c.Connections = nil    // Clears the slice completely to release resources
	if c.Cancel != nil {
        c.Cancel()
    }
}
