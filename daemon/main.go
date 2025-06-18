package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	rpServer := &server{}
	rpc.RegisterDaemonServiceServer(s, rpServer)

	m := tunnelmanager.NewTunnelManager()
	rpServer.RegisterTunnelManger(m)

	// restore tunnels that were active before restart
	restoreTunnels(m)

	// handle shutdown signals and persist tunnels
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		persistTunnels(m)
		os.Exit(0)
	}()

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		persistTunnels(m)
		log.Fatalf("failed to serve: %v", err)
	}
}

func restoreTunnels(m *tunnelmanager.TunnelManager) {
	dirpath := config.DefaultConfigDir
	if v := os.Getenv(config.ConfigDirFlagName); v != "" {
		dirpath = v
	}
	configdir, err := utils.ResolveDir(dirpath)
	if err != nil {
		log.Printf("failed to resolve config dir: %v", err)
		return
	}
	activeFile := filepath.Join(configdir, config.ActiveTunnelsFile)
	tunnels, err := tunnelmanager.LoadActiveTunnels(activeFile)
	if err != nil {
		log.Printf("failed to load active tunnels: %v", err)
		return
	}
	for _, t := range tunnels {
		if _, err := tasks.StartTunnelTask(context.Background(), &rpc.StartTunnelRequest{ConfigName: t.ConfigName, LocalPort: int32(t.LocalPort)}, m); err != nil {
			log.Printf("failed to restore tunnel %s: %v", t.ConfigName, err)
		}
	}
}

func persistTunnels(m *tunnelmanager.TunnelManager) {
	dirpath := config.DefaultConfigDir
	if v := os.Getenv(config.ConfigDirFlagName); v != "" {
		dirpath = v
	}
	configdir, err := utils.ResolveDir(dirpath)
	if err != nil {
		log.Printf("failed to resolve config dir: %v", err)
		return
	}
	if err := m.SaveActiveTunnels(filepath.Join(configdir, config.ActiveTunnelsFile)); err != nil {
		log.Printf("failed to save active tunnels: %v", err)
	}
}
