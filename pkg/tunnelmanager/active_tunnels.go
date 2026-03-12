package tunnelmanager

import (
	"encoding/json"
	"os"
	"time"
)

// ActiveTunnel represents a tunnel that should be restored on restart.
type ActiveTunnel struct {
	ConfigName string `json:"config_name"`
	LocalPort  int    `json:"local_port"`
	LocalAddr  string `json:"local_addr,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
	Server     string `json:"server,omitempty"`
	User       string `json:"user,omitempty"`
	SavedAt    string `json:"saved_at,omitempty"`
}

// SaveActiveTunnels persists current active tunnels to the given file.
func (m *TunnelManager) SaveActiveTunnels(path string) error {
	tunnels := make([]ActiveTunnel, 0)
	savedAt := time.Now().UTC().Format(time.RFC3339)
	m.Mutex.Lock()
	for port, ci := range m.Connections {
		tunnels = append(tunnels, ActiveTunnel{
			ConfigName: ci.Config.Name,
			LocalPort:  port,
			LocalAddr:  ci.LocalAddr,
			RemoteAddr: ci.RemoteAddr,
			Server:     ci.Config.Server,
			User:       ci.Config.User,
			SavedAt:    savedAt,
		})
	}
	m.Mutex.Unlock()

	if len(tunnels) == 0 {
		// remove file if no active tunnels
		os.Remove(path)
		return nil
	}

	data, err := json.MarshalIndent(tunnels, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadActiveTunnels reads persisted tunnels from the given file.
func LoadActiveTunnels(path string) ([]ActiveTunnel, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ActiveTunnel{}, nil
		}
		return nil, err
	}
	var tunnels []ActiveTunnel
	if err := json.Unmarshal(b, &tunnels); err != nil {
		return nil, err
	}
	return tunnels, nil
}
