package main

import (
	"log"
	"net"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/besrabasant/ssh-tunnel-manager/daemon/tasks"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
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
	
	m := tasks.NewTunnelManager() 
	rpServer.RegisterTunnelManger(m)
	

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
