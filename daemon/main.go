package main

import (
	"context"
	"log"
	"net"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDaemonServiceServer
}

func (s *server) ListConfigurations(ctx context.Context, req *pb.ListConfigurationsRequest) (*pb.ListConfigurationsResponse, error) {
	return tasks.ListConfigurationTask(ctx, req)
}

func (s *server) StartTunnel(ctx context.Context, req *pb.StartTunnelRequest) (*pb.StartTunnelResponse, error) {
	return tasks.StartTunnelTask(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDaemonServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
