package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func UpdateConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest,  manager *TunnelManager) (*rpc.AddOrUpdateConfigurationResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	openConns := manager.connections.Filter(func(c *ConnectionInfo) bool {
		return req.Name == c.Config.Name
	})

	if len(openConns) > 0 {

		var ports []int
		for port := range openConns {
			ports = append(ports, port)
		}

		output.WriteString(fmt.Sprint("Cannot update configuration as connection is open on port ", ports[0], "\n"))
		return &rpc.AddOrUpdateConfigurationResponse{Result: output.String()}, nil
	}

	configdir, err := utils.ResolveDir(configmanager.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	err = configmanager.NewManager(configdir).UpdateConfiguration(*convertRpcTunnelConfigToConfig(req.Data))
	if err != nil {
		output.WriteString(fmt.Sprintf("Cannot update configuration %s: %v", req.Name, err))
		return &rpc.AddOrUpdateConfigurationResponse{Result: output.String()}, nil
	}

	output.WriteString(fmt.Sprintf("Successfully updated configuration %s", req.Name))

	return &rpc.AddOrUpdateConfigurationResponse{Result: output.String()}, nil
}
