package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func AddConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.AddOrUpdateConfigurationResponse, error) {
	var output strings.Builder
	output.WriteString("\n")

	cfgs, err := getConfigs()
	if err != nil {
		return mutationResponse("\nError while adding configurations found\n", rpc.ResponseStatus_Error, "Error while reading configurations"), nil
	}

	cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
		return req.Name == c.Name
	})

	if len(cfgs) > 0 {
		message := fmt.Sprintf("A configuration already exists with name %s", req.Name)
		return mutationResponse(fmt.Sprintf("\nA congiguration already exits with name %s \n", req.Name), rpc.ResponseStatus_Error, message), nil
	}

	configdir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	err = configmanager.NewManager(configdir).AddConfiguration(*configmanager.ConvertRpcTunnelConfigToConfig(req.Data))
	if err != nil {
		message := fmt.Sprintf("Cannot add configuration %s: %v", req.Name, err)
		output.WriteString(message)
		return mutationResponse(output.String(), rpc.ResponseStatus_Error, message), nil
	}

	message := fmt.Sprintf("Successfully added new configuration %s", req.Name)
	output.WriteString(fmt.Sprintf("Successfully add new configuration %s", req.Name))

	return &rpc.AddOrUpdateConfigurationResponse{
		Result:  output.String(),
		Status:  rpc.ResponseStatus_Success,
		Message: message,
		Data:    req.Data,
	}, nil
}
