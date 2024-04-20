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
		return &rpc.AddOrUpdateConfigurationResponse{Result: "\nError while adding configurations found\n"}, nil
	}

	cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
		return req.Name == c.Name
	})

	if len(cfgs) > 0 {
		return &rpc.AddOrUpdateConfigurationResponse{Result: fmt.Sprintf("\nA congiguration already exits with name %s \n", req.Name)},nil
	}


	configdir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	err = configmanager.NewManager(configdir).AddConfiguration(*convertRpcTunnelConfigToConfig(req.Data))
	if err != nil {
		output.WriteString(fmt.Sprintf("Cannot add configuration %s: %v", req.Name, err))
		return &rpc.AddOrUpdateConfigurationResponse{Result: output.String()}, nil
	}

	output.WriteString(fmt.Sprintf("Successfully add new configuration %s", req.Name))

	return &rpc.AddOrUpdateConfigurationResponse{Result: output.String()}, nil
}