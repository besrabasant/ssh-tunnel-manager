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

func DeleteConfigurationJSON(ctx context.Context, req *pb.DeleteConfigurationRequest, manager *tunnelmanager.TunnelManager) (*pb.MutationResponse, error) {
	if req == nil || req.Name == "" {
		return &pb.MutationResponse{Status: pb.ResponseStatus_Error, Message: "missing config name"}, nil
	}

	// Block delete if active connection exists
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
			Message: fmt.Sprintf("Cannot delete configuration: connection open on port %d", port),
		}, nil
	}

	dir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}
	if err := configmanager.NewManager(dir).RemoveConfiguration(req.Name); err != nil {
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: fmt.Sprintf("Cannot delete configuration %s: %v", req.Name, err),
		}, nil
	}

	return &pb.MutationResponse{
		Status:  pb.ResponseStatus_Success,
		Message: fmt.Sprintf("Successfully deleted configuration %s", req.Name),
	}, nil
}
