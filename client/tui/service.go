package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

type Active struct {
	Name      string
	LocalPort int
}

func LoadConfigs(dir string) ([]configmanager.Entry, error) {
	return configmanager.NewManager(dir).GetConfigurations()
}

func LoadActive() ([]Active, error) {
	dir, err := utils.ResolveDir(config.DefaultConfigDir)
	if err != nil {
		return nil, err
	}
	p := filepath.Join(dir, config.ActiveTunnelsFile)
	ts, err := tunnelmanager.LoadActiveTunnels(p)
	if err != nil {
		return nil, err
	}
	out := make([]Active, 0, len(ts))
	for _, t := range ts {
		out = append(out, Active{Name: t.ConfigName, LocalPort: t.LocalPort})
	}
	return out, nil
}

func AddConfig(dir string, e configmanager.Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	return configmanager.NewManager(dir).AddConfiguration(e)
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

func StartTunnel(name string, localPort int) error {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = c.StartTunnel(ctx, &pb.StartTunnelRequest{ConfigName: name, LocalPort: int32(localPort)})
	if err != nil {
		return fmt.Errorf("start failed: %w", err)
	}
	return nil
}

func KillTunnel(name string, localPort int) error {
	c, cleanup, err := lib.CreateDaemonServiceClient()
	if err != nil {
		return fmt.Errorf("rpc connect: %w", err)
	}
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = c.KillTunnel(ctx, &pb.KillTunnelRequest{ConfigName: name, LocalPort: int32(localPort)})
	if err != nil {
		return fmt.Errorf("kill failed: %w", err)
	}
	return nil
}
