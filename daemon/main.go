package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/pkg/tunnelmanager"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Initialize components
	m := tunnelmanager.NewTunnelManager()
	cf := configmanager.NewManager(config.DefaultConfigDir)
	svc := tunnelmanager.NewTunnelService(m, cf, config.DefaultConfigDir)

	rpServer := &server{service: svc}
	rpc.RegisterDaemonServiceServer(s, rpServer)

	// restore tunnels that were active before restart
	if err := svc.RestoreTunnels(context.Background()); err != nil {
		log.Printf("failed to restore tunnels: %v", err)
	}

	// handle shutdown signals and persist tunnels
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := svc.PersistTunnels(); err != nil {
			log.Printf("failed to persist tunnels during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		if err := svc.PersistTunnels(); err != nil {
			log.Printf("failed to persist tunnels during error: %v", err)
		}
		log.Fatalf("failed to serve: %v", err)
	}
}

// remove unused functions
