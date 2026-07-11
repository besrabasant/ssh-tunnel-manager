package tasks

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func KillTunnelTask(ctx context.Context, req *rpc.KillTunnelRequest, service tunnelmanager.TunnelService) (*rpc.KillTunnelResponse, error) {
	result, err := service.StopTunnel(ctx, req.ConfigName, req.LocalPort)
	if err != nil {
		return nil, err
	}
	return &rpc.KillTunnelResponse{Result: result}, nil
}
