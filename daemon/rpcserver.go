package main

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

type server struct {
	rpc.UnimplementedDaemonServiceServer
	service tunnelmanager.TunnelService
}

func (s *server) RegisterTunnelService(service tunnelmanager.TunnelService) {
	s.service = service
}

func (s *server) ListConfigurations(ctx context.Context, req *rpc.ListConfigurationsRequest) (*rpc.ListConfigurationsResponse, error) {
	return tasks.ListConfigurationTask(ctx, req)
}

func (s *server) ListConfigurationsJSON(ctx context.Context, req *rpc.ListConfigurationsJSONRequest) (*rpc.ListConfigurationsJSONResponse, error) {
	return tasks.ListConfigurationsJSONTask(ctx, req)
}

func (s *server) FetchConfiguration(ctx context.Context, req *rpc.FetchConfigurationRequest) (*rpc.FetchConfigurationResponse, error) {
	return tasks.FetchTunnelConfigTask(ctx, req)
}

func (s *server) UpdateConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.AddOrUpdateConfigurationResponse, error) {
	return tasks.UpdateConfiguration(ctx, req, s.service)
}

func (s *server) UpdateConfigurationJSON(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.MutationResponse, error) {
	return tasks.UpdateConfigurationJSON(ctx, req, s.service)
}

func (s *server) AddConfiguration(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.AddOrUpdateConfigurationResponse, error) {
	return tasks.AddConfiguration(ctx, req)
}

func (s *server) AddConfigurationJSON(ctx context.Context, req *rpc.AddOrUpdateConfigurationRequest) (*rpc.MutationResponse, error) {
	return tasks.AddConfigurationJSON(ctx, req)
}

func (s *server) DeleteConfiguration(ctx context.Context, req *rpc.DeleteConfigurationRequest) (*rpc.DeleteConfigurationResponse, error) {
	return tasks.DeleteTunnelConfigTask(ctx, req, s.service)
}

func (s *server) DeleteConfigurationJSON(ctx context.Context, req *rpc.DeleteConfigurationRequest) (*rpc.MutationResponse, error) {
	return tasks.DeleteConfigurationJSON(ctx, req, s.service)
}

// Tunneling
func (s *server) StartTunnel(ctx context.Context, req *rpc.StartTunnelRequest) (*rpc.StartTunnelResponse, error) {
	return tasks.StartTunnelTask(ctx, req, s.service)
}

func (s *server) KillTunnel(ctx context.Context, req *rpc.KillTunnelRequest) (*rpc.KillTunnelResponse, error) {
	return tasks.KillTunnelTask(ctx, req, s.service)
}

func (s *server) ListActiveTunnels(ctx context.Context, req *rpc.ListActiveTunnelsRequest) (*rpc.ListActiveTunnelsResponse, error) {
	return tasks.ListActiveTunnelsTask(ctx, req, s.service)
}

func (s *server) ListActiveTunnelsJSON(ctx context.Context, req *rpc.ListActiveTunnelsJSONRequest) (*rpc.ListActiveTunnelsJSONResponse, error) {
	return tasks.ListActiveTunnelsJSONTask(ctx, req, s.service)
}
