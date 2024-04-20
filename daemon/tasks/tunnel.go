package tasks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

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

	// Check for open connections

	if len(manager.Connections) > 0 {

		portBusy := false
		for port := range manager.Connections {
			if port == int(req.LocalPort) {
				portBusy = true
				break
			}
		}

		if portBusy {
			output.WriteString(fmt.Sprint("\nCannot start tunnel as connection is open on port ", req.LocalPort, "\n"))
			return &rpc.StartTunnelResponse{Result: output.String()}, nil
		}
	}

	go manager.StartTunneling(context.Background(), cfg, int(req.LocalPort))

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

	return &rpc.StartTunnelResponse{Result: output.String()}, nil
}
