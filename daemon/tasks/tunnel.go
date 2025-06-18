package tasks

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

var randomPortGenerator = generateRandomPort

func StartTunnelTask(ctx context.Context, req *rpc.StartTunnelRequest, manager *tunnelmanager.TunnelManager) (*rpc.StartTunnelResponse, error) {
	var output strings.Builder
	manager.CreateResultChannels()

	dirpath := config.DefaultConfigDir
	if value := os.Getenv(config.ConfigDirFlagName); value != "" {
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

	localPort := req.LocalPort

	if localPort == -1 {

		localPort = int32(cfg.LocalPort)
	}

	if localPort == 0 {
		// Generate random port
		randomPort, err := randomPortGenerator()
		if err != nil {
			return nil, fmt.Errorf("couldn't generate a random port: %v", err)
		}
		localPort = int32(randomPort)
	}

	// Check for open connections
	if len(manager.Connections) > 0 {

		portBusy := false
		for port := range manager.Connections {
			if port == int(localPort) {
				portBusy = true
				break
			}
		}

		if portBusy {
			output.WriteString(fmt.Sprintf("\nCannot start tunnel as connection is already open on port %d (requested %d)\n", localPort, req.LocalPort))
			return &rpc.StartTunnelResponse{Result: output.String()}, nil
		}
	}

	go manager.StartTunneling(context.Background(), cfg, int(localPort))

	var errReceived bool = false

loop:
	for {
		select {
		case result, ok := <-manager.ResultChan:
			if !ok {
				if !errReceived {
					output.WriteString("Tunnel setup completed.\n")
				}
				break loop
			}
			output.WriteString(result + "\n")

		case err, ok := <-manager.ErrChan:
			if ok && err != nil { // Check if error is not nil
				output.WriteString(fmt.Sprintf("Failed to start tunneling: %v\n", err))
				errReceived = true
			}
		}
	}

	// persist running tunnels
	if err := manager.SaveActiveTunnels(filepath.Join(configdir, config.ActiveTunnelsFile)); err != nil {
		output.WriteString(fmt.Sprintf("Failed to persist active tunnels: %v\n", err))
	}

	return &rpc.StartTunnelResponse{Result: output.String()}, nil
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
