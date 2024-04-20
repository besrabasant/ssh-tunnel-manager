package tasks

import (
	"context"
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)


func FetchTunnelConfigTask(ctx context.Context, req *rpc.FetchConfigurationRequest) (*rpc.FetchConfigurationResponse, error) {
	cfgs, err := getConfigs()
	if err != nil {
		return &rpc.FetchConfigurationResponse{Status: rpc.ResponseStatus_Error, Message: "\nError while reading configurations found\n"}, nil
	}

	if len(cfgs) == 0 {
		return &rpc.FetchConfigurationResponse{Status: rpc.ResponseStatus_Error, Message: "\nNo configurations found\n"}, nil
	}

	cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
		return req.Name == c.Name
	})

	if len(cfgs) == 0 {
		return &rpc.FetchConfigurationResponse{Status: rpc.ResponseStatus_Error, Message: fmt.Sprintf("\nNo configurations found with name %s \n", req.Name)}, nil
	}

	return &rpc.FetchConfigurationResponse{Status: rpc.ResponseStatus_Success, Message: "", Data: configmanager.ConvertConfigToRpcTunnelConfig(&cfgs[0])}, nil
}
