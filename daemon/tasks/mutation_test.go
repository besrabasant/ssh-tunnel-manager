package tasks

import (
	"context"
	"strings"
	"testing"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func TestAddConfigurationPopulatesLegacyAndStructuredOutput(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(config.ConfigDirFlagName, t.TempDir())

	data := &rpc.TunnelConfig{
		Name:       "dev",
		Server:     "example.com:22",
		User:       "root",
		KeyFile:    "/tmp/id_ed25519",
		RemoteHost: "127.0.0.1",
		RemotePort: 6443,
		LocalPort:  6444,
	}
	response, err := AddConfiguration(context.Background(), &rpc.AddOrUpdateConfigurationRequest{
		Name: "dev",
		Data: data,
	})
	if err != nil {
		t.Fatal(err)
	}

	if response.GetStatus() != rpc.ResponseStatus_Success {
		t.Fatalf("expected success status, got %v", response.GetStatus())
	}
	if response.GetData().GetName() != "dev" {
		t.Fatalf("expected structured config data, got %#v", response.GetData())
	}
	if !strings.Contains(response.GetResult(), "Successfully add new configuration dev") {
		t.Fatalf("expected legacy result, got %q", response.GetResult())
	}
	if response.GetMessage() != "Successfully added new configuration dev" {
		t.Fatalf("expected structured message, got %q", response.GetMessage())
	}
}
