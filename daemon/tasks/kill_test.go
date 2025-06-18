package tasks

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func TestKillTunnelTask_NoDeadlock(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv(config.ConfigDirFlagName, tmp)
	defer os.Unsetenv(config.ConfigDirFlagName)

	manager := tunnelmanager.NewTunnelManager()
	manager.Connections[1234] = &tunnelmanager.ConnectionInfo{
		Config: configmanager.Entry{Name: "test"},
		Cancel: func() {},
	}

	done := make(chan struct{})
	go func() {
		KillTunnelTask(context.Background(), &rpc.KillTunnelRequest{LocalPort: 1234}, manager)
		close(done)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("KillTunnelTask did not return")
	case <-done:
		// success
	}
}
