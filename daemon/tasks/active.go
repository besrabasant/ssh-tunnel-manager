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
	output.WriteString("\nActive Tunnels:\n")
	output.WriteString("---------------\n")
	for _, t := range tunnels {
		output.WriteString(fmt.Sprintf("Config: %s, Local Port: %d, Remote: %s, Server: %s\n", t.ConfigName, t.LocalPort, t.RemoteAddr, t.Server))
	}

	return &rpc.ListActiveTunnelsResponse{Result: output.String()}, nil
}
