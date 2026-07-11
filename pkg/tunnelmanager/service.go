package tunnelmanager

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
)

type tunnelService struct {
	manager   TunnelManager
	cfgMgr    configmanager.ConfigManager
	configDir string
}

func NewTunnelService(manager TunnelManager, cfgMgr configmanager.ConfigManager, configDir string) TunnelService {
	return &tunnelService{
		manager:   manager,
		cfgMgr:    cfgMgr,
		configDir: configDir,
	}
}

func (s *tunnelService) StartTunnel(ctx context.Context, configName string, localPort int32) (string, error) {
	cfg, err := s.cfgMgr.GetConfiguration(configName)
	if err != nil {
		return "", fmt.Errorf("couldn't get configuration %q: %v", configName, err)
	}

	actualPort := localPort
	if actualPort == -1 {
		actualPort = int32(cfg.LocalPort)
	}

	if actualPort == 0 {
		var err error
		actualPortInt, err := s.generateRandomPort()
		if err != nil {
			return "", fmt.Errorf("failed to generate random port: %v", err)
		}
		actualPort = int32(actualPortInt)
	}

	// Check for open connections
	for port := range s.manager.GetConnections() {
		if port == int(actualPort) {
			return fmt.Sprintf("\nCannot start tunnel as connection is already open on port %d\n", actualPort), nil
		}
	}

	s.manager.CreateResultChannels()
	go s.manager.StartTunneling(context.Background(), cfg, int(actualPort))

	// Collect results
	var output string
loop:
	for {
		select {
		case result, ok := <-s.manager.GetResultChan():
			if !ok {
				break loop
			}
			output += result + "\n"
		case err, ok := <-s.manager.GetErrChan():
			if ok && err != nil {
				output += fmt.Sprintf("Failed to start tunneling: %v\n", err)
			}
		}
	}

	s.PersistTunnels()
	return output, nil
}

func (s *tunnelService) StopTunnel(ctx context.Context, configName string, localPort int32) (string, error) {
	var output strings.Builder
	var connInfo *ConnectionInfo
	var connPort int

	connections := s.manager.GetConnections()

	if configName != "" {
		for port, ci := range connections {
			if ci.Config.Name == configName {
				connInfo = ci
				connPort = port
				break
			}
		}
	} else {
		connPort = int(localPort)
		connInfo = connections[connPort]
	}

	if connInfo == nil {
		if configName != "" {
			return fmt.Sprintf("Did not find any connection for configuration %s", configName), nil
		}
		return fmt.Sprintf("Did not find any connection on port %d", localPort), nil
	}

	output.WriteString(fmt.Sprintf("\nClosing existing connection on port %d for %s\n", connPort, connInfo.Config.Name))
	connInfo.ClearConnection()

	s.PersistTunnels()
	return output.String(), nil
}

func (s *tunnelService) ListActiveTunnels(ctx context.Context) ([]ActiveTunnel, error) {
	tunnels := make([]ActiveTunnel, 0)
	for port, ci := range s.manager.GetConnections() {
		tunnels = append(tunnels, ActiveTunnel{
			ConfigName: ci.Config.Name,
			LocalPort:  port,
			LocalAddr:  ci.LocalAddr,
			RemoteAddr: ci.RemoteAddr,
			Server:     ci.Config.Server,
			User:       ci.Config.User,
		})
	}
	return tunnels, nil
}

func (s *tunnelService) RestoreTunnels(ctx context.Context) error {
	configdir, err := utils.ResolveDir(s.configDir)
	if err != nil {
		return err
	}
	activeFile := filepath.Join(configdir, config.ActiveTunnelsFile)
	tunnels, err := LoadActiveTunnels(activeFile)
	if err != nil {
		return err
	}
	for _, t := range tunnels {
		// Start in background to avoid blocking daemon startup
		go s.StartTunnel(context.Background(), t.ConfigName, int32(t.LocalPort))
	}
	return nil
}

func (s *tunnelService) PersistTunnels() error {
	configdir, err := utils.ResolveDir(s.configDir)
	if err != nil {
		return err
	}
	return s.manager.SaveActiveTunnels(filepath.Join(configdir, config.ActiveTunnelsFile))
}
func (s *tunnelService) GetManager() TunnelManager {
	return s.manager
}

func (s *tunnelService) generateRandomPort() (int, error) {
	// Listen on port 0 to bind to a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// Extract the port number from the listener address
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, err
	}

	// Convert the port number to an integer
	randomPort, err := net.LookupPort("tcp", port)
	if err != nil {
		return 0, err
	}

	return randomPort, nil
}
