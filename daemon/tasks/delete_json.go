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

func DeleteConfigurationJSON(ctx context.Context, req *pb.DeleteConfigurationRequest, service tunnelmanager.TunnelService) (*pb.MutationResponse, error) {
	if req == nil || req.Name == "" {
		return &pb.MutationResponse{Status: pb.ResponseStatus_Error, Message: "missing config name"}, nil
	}

	// Block delete if active connection exists
	connections := service.GetManager().GetConnections()
	var port int = -1
	for p, ci := range connections {
		if ci.Config.Name == req.Name {
			port = p
			break
		}
	}

	if port != -1 {
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
