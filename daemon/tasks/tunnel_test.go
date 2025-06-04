package tasks

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func createConfig(t *testing.T, dir string, entry configmanager.Entry) {
	t.Helper()
	// ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file, err := json.MarshalIndent(entry, "", " ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, entry.Name+".json"), file, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func TestStartTunnelTask_ConflictWithConfigPort(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv(config.ConfigDirFlagName, tmp)
	defer os.Unsetenv(config.ConfigDirFlagName)

	entry := configmanager.Entry{
		Name:       "test",
		Server:     "server:22",
		User:       "user",
		KeyFile:    "key",
		RemoteHost: "remote",
		RemotePort: 22,
		LocalPort:  5432,
	}
	createConfig(t, tmp, entry)

	manager := tunnelmanager.NewTunnelManager()
	manager.Connections[5432] = &tunnelmanager.ConnectionInfo{}

	req := &rpc.StartTunnelRequest{ConfigName: "test", LocalPort: -1}
	resp, err := StartTunnelTask(context.Background(), req, manager)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "already open on port 5432") {
		t.Fatalf("expected conflict message, got %q", resp.Result)
	}
}

func TestStartTunnelTask_ConflictWithRandomPort(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv(config.ConfigDirFlagName, tmp)
	defer os.Unsetenv(config.ConfigDirFlagName)

	entry := configmanager.Entry{
		Name:       "test",
		Server:     "server:22",
		User:       "user",
		KeyFile:    "key",
		RemoteHost: "remote",
		RemotePort: 22,
		LocalPort:  1111,
	}
	createConfig(t, tmp, entry)

	manager := tunnelmanager.NewTunnelManager()
	manager.Connections[5432] = &tunnelmanager.ConnectionInfo{}

	// override randomPortGenerator
	original := randomPortGenerator
	randomPortGenerator = func() (int, error) { return 5432, nil }
	defer func() { randomPortGenerator = original }()

	req := &rpc.StartTunnelRequest{ConfigName: "test", LocalPort: 0}
	resp, err := StartTunnelTask(context.Background(), req, manager)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Result, "already open on port 5432") {
		t.Fatalf("expected conflict message, got %q", resp.Result)
	}
}
