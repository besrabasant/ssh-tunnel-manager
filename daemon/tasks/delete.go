package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func DeleteTunnelConfigTask(ctx context.Context, req *rpc.DeleteConfigurationRequest, service tunnelmanager.TunnelService) (*rpc.DeleteConfigurationResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	// Check for open connections
	for port, ci := range service.GetManager().GetConnections() {
		if ci.Config.Name == req.Name {
			return nil, fmt.Errorf("cannot delete configuration %q as a connection is open on port %d", req.Name, port)
		}
	}

	cfgs, err := getConfigs()
	if err != nil {
		return &rpc.DeleteConfigurationResponse{Result: "\nError while reading configurations found\n"}, nil
	}

	if len(cfgs) == 0 {
		return &rpc.DeleteConfigurationResponse{Result: "\nNo configurations found\n"}, nil
	}

	configdir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	err = configmanager.NewManager(configdir).RemoveConfiguration(req.GetName())
	if err != nil {
		output.WriteString(fmt.Sprintf("Cannot delete configuration %s: %v", req.Name, err))
		return &rpc.DeleteConfigurationResponse{Result: output.String()}, nil
	}

	if err := service.PersistTunnels(); err != nil {
		output.WriteString(fmt.Sprintf("\nFailed to persist active tunnels: %v\n", err))
	}

	output.WriteString(fmt.Sprintf("Successfully deleted configuration %s", req.Name))

	return &rpc.DeleteConfigurationResponse{Result: output.String()}, nil
}
