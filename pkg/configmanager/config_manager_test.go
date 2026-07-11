package configmanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManagerResolvesAndCreatesHomeRelativeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	manager := NewManager("~/.ssh-tunnel-manager").(*manager)
	expected := filepath.Join(home, ".ssh-tunnel-manager")

	if manager.dir != expected {
		t.Fatalf("expected manager directory %q, got %q", expected, manager.dir)
	}
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("expected manager directory to exist: %v", err)
	}
}
