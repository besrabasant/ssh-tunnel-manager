package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func ListActiveTunnelsTask(ctx context.Context, req *rpc.ListActiveTunnelsRequest, service tunnelmanager.TunnelService) (*rpc.ListActiveTunnelsResponse, error) {
	tunnels, err := service.ListActiveTunnels(ctx)
	if err != nil {
		return nil, err
	}

	var output strings.Builder
	output.WriteString("\n")
	for i, t := range tunnels {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(fmt.Sprintf(":%d\n", t.LocalPort))
		output.WriteString(fmt.Sprintf("- Connection:                  %s\n", t.ConfigName))
		output.WriteString(fmt.Sprintf("- Remote Address:              %s\n", t.RemoteAddr))
		output.WriteString(fmt.Sprintf("- Local Address:               %s\n", t.LocalAddr))
		output.WriteString(fmt.Sprintf("- SSH server:                  %s\n", t.Server))
	}

	active := make([]*rpc.ActiveTunnel, 0, len(tunnels))
	for _, t := range tunnels {
		active = append(active, &rpc.ActiveTunnel{
			Name:       t.ConfigName,
			LocalPort:  int32(t.LocalPort),
			RemoteAddr: t.RemoteAddr,
			LocalAddr:  t.LocalAddr,
			Server:     t.Server,
		})
	}

	return &rpc.ListActiveTunnelsResponse{Result: output.String(), Tunnels: active}, nil
}
