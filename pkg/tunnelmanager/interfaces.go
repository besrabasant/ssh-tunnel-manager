package tunnelmanager

import (
	"context"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
)

// TunnelManager defines the low-level operations for managing individual SSH tunnels.
type TunnelManager interface {
	CreateResultChannels() (chan string, chan error)
	StartTunneling(ctx context.Context, entry configmanager.Entry, localPort int, resultChan chan<- string, errChan chan<- error)
	SaveActiveTunnels(path string) error
	GetResultChan() <-chan string
	GetErrChan() <-chan error
	GetConnections() SSHConnections
	Shutdown()
}

// TunnelService defines the high-level operations for managing tunnels and their persistence.
type TunnelService interface {
	StartTunnel(ctx context.Context, configName string, localPort int32) (string, error)
	StopTunnel(ctx context.Context, configName string, localPort int32) (string, error)
	ListActiveTunnels(ctx context.Context) ([]ActiveTunnel, error)
	RestoreTunnels(ctx context.Context) error
	PersistTunnels() error
	GetManager() TunnelManager
}
