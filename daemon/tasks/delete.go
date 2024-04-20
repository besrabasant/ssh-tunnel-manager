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

func DeleteTunnelConfigTask(ctx context.Context, req *rpc.DeleteConfigurationRequest,  manager *tunnelmanager.TunnelManager) (*rpc.DeleteConfigurationResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	cfgs, err := getConfigs()
	if err != nil {
		return &rpc.DeleteConfigurationResponse{Result: "\nError while reading configurations found\n"}, nil
	}

	if len(cfgs) == 0 {
		return &rpc.DeleteConfigurationResponse{Result: "\nNo configurations found\n"}, nil
	}

	cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
		return req.Name == c.Name
	})

	if len(cfgs) == 0 {
		return &rpc.DeleteConfigurationResponse{Result: fmt.Sprintf("\nNo configurations found with name %s \n", req.Name)}, nil
	}

	// Check for open connections
	openConns := manager.Connections.Filter(func(c *tunnelmanager.ConnectionInfo) bool {
		return req.Name == c.Config.Name
	})

	if len(openConns) > 0 {

		var ports []int
		for port := range openConns {
			ports = append(ports, port)
		}

		output.WriteString(fmt.Sprint("Cannot delete configuration as connection is open on port ", ports[0], "\n"))
		return &rpc.DeleteConfigurationResponse{Result: output.String()}, nil
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

	output.WriteString(fmt.Sprintf("Successfully deleted configuration %s", req.Name))

	return &rpc.DeleteConfigurationResponse{Result: output.String()}, nil
}
