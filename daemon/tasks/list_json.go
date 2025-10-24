package tasks

import (
	"context"
	"fmt"

	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func ListConfigurationsJSONTask(ctx context.Context, req *pb.ListConfigurationsJSONRequest) (*pb.ListConfigurationsJSONResponse, error) {
	cfgs, err := getConfigs()
	if err != nil {
		return nil, fmt.Errorf("couldn't get saved configurations: %v", err)
	}
	if len(cfgs) == 0 {
		return &pb.ListConfigurationsJSONResponse{Configs: []*pb.TunnelConfig{}}, nil
	}

	resp := &pb.ListConfigurationsJSONResponse{Configs: make([]*pb.TunnelConfig, 0, len(cfgs))}
	for _, e := range cfgs {
		resp.Configs = append(resp.Configs, &pb.TunnelConfig{
			Name:        e.Name,
			Description: e.Description,
			Server:      e.Server,
			User:        e.User,
			KeyFile:     e.KeyFile,
			RemoteHost:  e.RemoteHost,
			RemotePort:  int32(e.RemotePort),
			LocalPort:   int32(e.LocalPort),
		})
	}

	return resp, nil
}
