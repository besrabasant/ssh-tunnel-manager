package tasks

import (
	"context"
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

// AddConfigurationJSON is a structured variant of AddConfiguration that returns statusmessagedata.
// It preserves the same validation and duplicate checks as the string-based task.
func AddConfigurationJSON(ctx context.Context, req *pb.AddOrUpdateConfigurationRequest) (*pb.MutationResponse, error) {
	if req == nil || req.Data == nil || req.Data.Name == "" {
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: "missing config name",
		}, nil
	}

	// Reuse the same discovery path used by the string-based list/add tasks.
	cfgs, err := getConfigs()
	if err != nil {
		return nil, fmt.Errorf("read configs: %w", err)
	}

	// Prevent duplicate by name (same behavior as non-JSON path)
	cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
		return req.Name == c.Name
	})
	if len(cfgs) > 0 {
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: fmt.Sprintf("A configuration already exists with name %q", req.Name),
		}, nil
	}

	// Resolve daemon config dir and add
	dir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}

	if err := configmanager.NewManager(dir).AddConfiguration(*configmanager.ConvertRpcTunnelConfigToConfig(req.Data)); err != nil {
		return &pb.MutationResponse{
			Status:  pb.ResponseStatus_Error,
			Message: fmt.Sprintf("Cannot add configuration %s: %v", req.Name, err),
		}, nil
	}

	return &pb.MutationResponse{
		Status:  pb.ResponseStatus_Success,
		Message: fmt.Sprintf("Successfully added configuration %s", req.Name),
		Data:    req.Data,
	}, nil
}
