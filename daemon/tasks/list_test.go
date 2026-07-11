package tasks

import (
	"context"
	"strings"
	"testing"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func TestListConfigurationTaskPopulatesLegacyAndStructuredOutput(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.ConfigDirFlagName, dir)

	entry := configmanager.Entry{
		Name:        "dev",
		Description: "development",
		Server:      "example.com:22",
		User:        "root",
		KeyFile:     "/tmp/id_ed25519",
		RemoteHost:  "127.0.0.1",
		RemotePort:  6443,
		LocalPort:   6444,
	}
	if err := configmanager.NewManager(dir).AddConfiguration(entry); err != nil {
		t.Fatal(err)
	}

	response, err := ListConfigurationTask(context.Background(), &rpc.ListConfigurationsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.GetConfigs()) != 1 {
		t.Fatalf("expected one structured config, got %d", len(response.GetConfigs()))
	}
	if !strings.Contains(response.GetResult(), "dev") {
		t.Fatalf("expected legacy result to contain config name, got %q", response.GetResult())
	}
}

func TestListConfigurationTaskEmptyResultKeepsLegacyMessage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.ConfigDirFlagName, dir)

	response, err := ListConfigurationTask(context.Background(), &rpc.ListConfigurationsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.GetConfigs()) != 0 {
		t.Fatalf("expected no structured configs, got %d", len(response.GetConfigs()))
	}
	if !strings.Contains(response.GetResult(), "No configurations found") {
		t.Fatalf("expected empty legacy message, got %q", response.GetResult())
	}
}
