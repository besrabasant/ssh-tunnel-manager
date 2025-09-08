package tasks

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func ListActiveTunnelsJSONTask(ctx context.Context, _ *pb.ListActiveTunnelsJSONRequest, manager *tunnelmanager.TunnelManager) (*pb.ListActiveTunnelsJSONResponse, error) {
	resp := &pb.ListActiveTunnelsJSONResponse{Tunnels: make([]*pb.ActiveTunnel, 0, len(manager.Connections))}
	for port, conn := range manager.Connections {
		resp.Tunnels = append(resp.Tunnels, &pb.ActiveTunnel{
			Name:       conn.Config.Name,
			LocalPort:  int32(port),
			RemoteAddr: conn.RemoteAddr,
			LocalAddr:  conn.LocalAddr,
			Server:     conn.Config.Server,
		})
	}
	return resp, nil
}
