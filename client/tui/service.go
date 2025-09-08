package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
)

type Active struct {
	Name      string
	LocalPort int
}

// LoadConfigs now uses daemon RPC (JSON) instead of reading files directly.
// The dir parameter is ignored on purpose; the daemon is the source of truth.
func LoadConfigs(_ string) ([]configmanager.Entry, error) {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return nil, fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	rj, err := c.ListConfigurationsJSON(ctx, &pb.ListConfigurationsJSONRequest{})
	if err != nil {
		return nil, fmt.Errorf("list configurations (json) rpc: %w", err)
	}

	outs := make([]configmanager.Entry, 0, len(rj.GetConfigs()))
	for _, tc := range rj.GetConfigs() {
		outs = append(outs, configmanager.Entry{
			Name:        tc.GetName(),
			Description: tc.GetDescription(),
			Server:      tc.GetServer(),
			User:        tc.GetUser(),
			KeyFile:     tc.GetKeyFile(),
			RemoteHost:  tc.GetRemoteHost(),
			RemotePort:  int(tc.GetRemotePort()),
			LocalPort:   int(tc.GetLocalPort()),
		})
	}
	return outs, nil
}
func LoadActive() ([]Active, error) {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return nil, fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	rj, err := c.ListActiveTunnelsJSON(ctx, &pb.ListActiveTunnelsJSONRequest{})
	if err != nil {
		return nil, fmt.Errorf("list active tunnels (json) rpc: %w", err)
	}

	out := make([]Active, 0, len(rj.GetTunnels()))
	for _, t := range rj.GetTunnels() {
		out = append(out, Active{
			Name:      t.GetName(),
			LocalPort: int(t.GetLocalPort()),
		})
	}
	return out, nil
}

func AddConfig(_ string, e configmanager.Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	resp, err := c.AddConfigurationJSON(ctx, &pb.AddOrUpdateConfigurationRequest{
		Name: e.Name,
		Data: &pb.TunnelConfig{
			Name:        e.Name,
			Description: e.Description,
			Server:      e.Server,
			User:        e.User,
			KeyFile:     e.KeyFile,
			RemoteHost:  e.RemoteHost,
			RemotePort:  int32(e.RemotePort),
			LocalPort:   int32(e.LocalPort),
		},
	})
	if err != nil {
		return fmt.Errorf("add config (json) rpc: %w", err)
	}
	if resp.GetStatus() == pb.ResponseStatus_Error {
		return fmt.Errorf(resp.GetMessage())
	}
	return nil
}

func UpdateConfig(dir string, e configmanager.Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	return configmanager.NewManager(dir).UpdateConfiguration(e)
}
func DeleteConfig(dir, name string) error {
	return configmanager.NewManager(dir).RemoveConfiguration(name)
}

func StartTunnel(name string, localPort int) (string, error) {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return "", fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := c.StartTunnel(ctx, &pb.StartTunnelRequest{ConfigName: name, LocalPort: int32(localPort)})
	if err != nil {
		return "", fmt.Errorf("start failed: %w", err)
	}
	return resp.GetResult(), nil
}

func KillTunnel(name string, localPort int) (string, error) {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return "", fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := c.KillTunnel(ctx, &pb.KillTunnelRequest{ConfigName: name, LocalPort: int32(localPort)})
	if err != nil {
		return "", fmt.Errorf("kill failed: %w", err)
	}
	return resp.GetResult(), nil
}
