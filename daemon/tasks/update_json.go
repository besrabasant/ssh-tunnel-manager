package tasks

import (
	"context"
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

func UpdateConfigurationJSON(ctx context.Context, req *pb.AddOrUpdateConfigurationRequest, manager *tunnelmanager.TunnelManager) (*pb.MutationResponse, error) {
	if req == nil || req.Data == nil || req.Data.Name == "" {
		return &pb.MutationResponse{Status: pb.ResponseStatus_Error, Message: "missing config name"}, nil
	}

	// Block update if there's an active connection using this config (same behavior as legacy string RPC)
	openConns := manager.Connections.Filter(func(c *tunnelmanager.ConnectionInfo) bool {
		return req.Name == c.Config.Name
	})
	if len(openConns) > 0 {
		var port int
		for p := range openConns {
			port = p
			break
		}
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: fmt.Sprintf("Cannot update configuration: connection open on port %d", port),
		}, nil
	}

	// Resolve daemon config dir and update
	dir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}
	if err := configmanager.NewManager(dir).UpdateConfiguration(*configmanager.ConvertRpcTunnelConfigToConfig(req.Data)); err != nil {
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: fmt.Sprintf("Cannot update configuration %s: %v", req.Name, err),
		}, nil
	}

	return &pb.MutationResponse{
		Status:  pb.ResponseStatus_Success,
		Message: fmt.Sprintf("Successfully updated configuration %s", req.Name),
		Data:    req.Data,
	}, nil
}
