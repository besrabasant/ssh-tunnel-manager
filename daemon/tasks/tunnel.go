package tasks

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func StartTunnelTask(ctx context.Context, req *rpc.StartTunnelRequest, service tunnelmanager.TunnelService) (*rpc.StartTunnelResponse, error) {
	result, err := service.StartTunnel(ctx, req.ConfigName, req.LocalPort)
	if err != nil {
		return nil, err
	}
	status, events := operationParts(result)
	return &rpc.StartTunnelResponse{Result: result, Status: status, Events: events}, nil
}
