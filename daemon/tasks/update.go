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

func UpdateConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest, service tunnelmanager.TunnelService) (*rpc.AddOrUpdateConfigurationResponse, error) {
	var output strings.Builder

	output.WriteString("\n")

	// Check for open connections
	connections := service.GetManager().GetConnections()
	var openConns tunnelmanager.SSHConnections = make(map[int]*tunnelmanager.ConnectionInfo)
	for port, ci := range connections {
		if ci.Config.Name == req.Name {
			openConns[port] = ci
		}
	}

	if len(openConns) > 0 {
		var ports []int
		for port := range openConns {
			ports = append(ports, port)
		}

		message := fmt.Sprint("Cannot update configuration as connection is open on port ", ports[0])
		output.WriteString(message + "\n")
		return mutationResponse(output.String(), rpc.ResponseStatus_Error, message), nil
	}

	configdir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	err = configmanager.NewManager(configdir).UpdateConfiguration(*configmanager.ConvertRpcTunnelConfigToConfig(req.Data))
	if err != nil {
		message := fmt.Sprintf("Cannot update configuration %s: %v", req.Name, err)
		output.WriteString(message)
		return mutationResponse(output.String(), rpc.ResponseStatus_Error, message), nil
	}

	message := fmt.Sprintf("Successfully updated configuration %s", req.Name)
	output.WriteString(message)

	return &rpc.AddOrUpdateConfigurationResponse{
		Result:  output.String(),
		Status:  rpc.ResponseStatus_Success,
		Message: message,
		Data:    req.Data,
	}, nil
}
