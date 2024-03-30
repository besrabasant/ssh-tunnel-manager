package main

import (
	"context"
	"log"
	"net"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"google.golang.org/grpc"
)

type server struct {
	rpc.UnimplementedDaemonServiceServer
}

func (s *server) ListConfigurations(ctx context.Context, req *rpc.ListConfigurationsRequest) (*rpc.ListConfigurationsResponse, error) {
	return tasks.ListConfigurationTask(ctx, req)
}

func (s *server) FetchConfiguration(ctx context.Context, req *rpc.FetchConfigurationRequest) (*rpc.FetchConfigurationResponse, error) {
	return tasks.FetchTunnelConfigTask(ctx, req)
}

func (s *server) StartTunnel(ctx context.Context, req *rpc.StartTunnelRequest) (*rpc.StartTunnelResponse, error) {
	return tasks.StartTunnelTask(ctx, req)
}

func (s *server) KillTunnel(ctx context.Context, req *rpc.KillTunnelRequest) (*rpc.KillTunnelResponse, error) {
	return tasks.KillTunnelTask(ctx, req)
}

func (s *server) ListActiveTunnels(ctx context.Context, req *rpc.ListActiveTunnelsRequest) (*rpc.ListActiveTunnelsResponse, error) {
	return tasks.ListActiveTunnelsTask(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	rpc.RegisterDaemonServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
