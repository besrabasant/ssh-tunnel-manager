package main

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

type server struct {
	rpc.UnimplementedDaemonServiceServer
	manager  *tunnelmanager.TunnelManager
}

func (s *server) RegisterTunnelManger(manager *tunnelmanager.TunnelManager) {
	s.manager = manager
}

func (s *server) ListConfigurations(ctx context.Context, req *rpc.ListConfigurationsRequest) (*rpc.ListConfigurationsResponse, error) {
	return tasks.ListConfigurationTask(ctx, req)
}

func (s *server) FetchConfiguration(ctx context.Context, req *rpc.FetchConfigurationRequest) (*rpc.FetchConfigurationResponse, error) {
	return tasks.FetchTunnelConfigTask(ctx, req)
}

func (s *server) UpdateConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.AddOrUpdateConfigurationResponse, error) {
	return tasks.UpdateConfiguration(ctx, req, s.manager)
}

func (s *server) AddConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.AddOrUpdateConfigurationResponse, error) {
	return tasks.AddConfiguration(ctx, req)
}

func (s *server) DeleteConfiguration(ctx context.Context, req *rpc.DeleteConfigurationRequest) (*rpc.DeleteConfigurationResponse, error) {
	return tasks.DeleteTunnelConfigTask(ctx, req, s.manager)
}

// Tunneling
func (s *server) StartTunnel(ctx context.Context, req *rpc.StartTunnelRequest) (*rpc.StartTunnelResponse, error) {
	return tasks.StartTunnelTask(ctx, req, s.manager)
}

func (s *server) KillTunnel(ctx context.Context, req *rpc.KillTunnelRequest) (*rpc.KillTunnelResponse, error) {
	return tasks.KillTunnelTask(ctx, req, s.manager)
}

func (s *server) ListActiveTunnels(ctx context.Context, req *rpc.ListActiveTunnelsRequest) (*rpc.ListActiveTunnelsResponse, error) {
	return tasks.ListActiveTunnelsTask(ctx, req, s.manager)
}
