package tunnelmanager

import (
	"fmt"
	"net"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"golang.org/x/crypto/ssh"
)

// SSH Connection info
type ConnectionInfo struct {
	Client      *ssh.Client
	Listener    net.Listener
	LocalAddr   string
	RemoteAddr  string
	Connections []net.Conn
	Config      configmanager.Entry
	Cancel      func()
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
	if len(c.Connections) > 0 {
		for _, conn := range c.Connections {
			conn.Close()
		}
	}
}
