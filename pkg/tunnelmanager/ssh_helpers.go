package tunnelmanager

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

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
	log.Printf("Connection closed: %s", local.RemoteAddr())
}
