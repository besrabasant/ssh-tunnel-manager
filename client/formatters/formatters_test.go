package formatters

import (
	"bytes"
	"strings"
	"testing"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
)

func TestActiveTunnelsFormatterPlainText(t *testing.T) {
	output := NewActiveTunnelsFormatter(&bytes.Buffer{}).Format(&rpc.ListActiveTunnelsResponse{
		Tunnels: []*rpc.ActiveTunnel{{
			Name:       "nesecure-dev-probe",
			LocalPort:  6444,
			RemoteAddr: "127.0.0.1:6443",
			Server:     "160.30.146.243",
		}},
	})

	if strings.Contains(output, "\x1b[") {
		t.Fatalf("expected plain output for a non-terminal writer, got ANSI escapes: %q", output)
	}
	for _, expected := range []string{
		":6444",
		"- Connection:                  nesecure-dev-probe",
		"- Remote Address:              127.0.0.1:6443",
		"- Local Address:               127.0.0.1:6444",
		"- SSH server:                  160.30.146.243",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got %q", expected, output)
		}
	}
}

func TestConfigurationListFormatterFallsBackToLegacyResult(t *testing.T) {
	legacy := "\nNo configurations found\n"
	output := NewConfigurationListFormatter(&bytes.Buffer{}).Format(&rpc.ListConfigurationsResponse{Result: legacy})
	if output != legacy {
		t.Fatalf("expected legacy result %q, got %q", legacy, output)
	}
}

func TestConfigurationListFormatterUsesConnectionLayout(t *testing.T) {
	output := NewConfigurationListFormatter(&bytes.Buffer{}).Format(&rpc.ListConfigurationsResponse{
		Configs: []*rpc.TunnelConfig{{
			Name:       "nesecure-dev-probe",
			Server:     "160.30.146.243",
			User:       "root",
			KeyFile:    "/home/user/.ssh/id_ed25519",
			RemoteHost: "127.0.0.1",
			RemotePort: 6443,
			LocalPort:  6444,
		}},
	})

	for _, expected := range []string{
		":6444",
		"- Connection:                  nesecure-dev-probe",
		"- Remote Address:              127.0.0.1:6443",
		"- Local Address:               127.0.0.1:6444",
		"- SSH server:                  160.30.146.243",
		"- User:                        root",
		"- Private key:                 /home/user/.ssh/id_ed25519",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got %q", expected, output)
		}
	}
}

func TestMutationFormatterSuccessAndError(t *testing.T) {
	formatter := NewMutationFormatter(&bytes.Buffer{})

	success := formatter.Format(MutationResult{
		Status:     rpc.ResponseStatus_Success,
		Message:    "Successfully added configuration dev",
		Structured: true,
	})
	if !strings.Contains(success, "Successfully added configuration dev") {
		t.Fatalf("expected success message, got %q", success)
	}

	errorMessage := "Cannot delete configuration dev"
	errorOutput := formatter.Format(MutationResult{
		Status:     rpc.ResponseStatus_Error,
		Message:    errorMessage,
		Structured: true,
	})
	if !strings.Contains(errorOutput, errorMessage) {
		t.Fatalf("expected error message, got %q", errorOutput)
	}
}

func TestOperationFormatterEvents(t *testing.T) {
	output := NewOperationFormatter(&bytes.Buffer{}).Format(OperationResult{
		Status:     rpc.ResponseStatus_Success,
		Events:     []string{"Connecting", "Connected"},
		Structured: true,
	})

	if !strings.Contains(output, "Connecting\nConnected\n") {
		t.Fatalf("expected operation events, got %q", output)
	}
}
