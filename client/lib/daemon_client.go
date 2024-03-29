package lib

import (
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateDaemonServiceClient() (pb.DaemonServiceClient, func(), error) {
	conn, err := grpc.Dial(config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}

	cleanup := func() { conn.Close() }
	c := pb.NewDaemonServiceClient(conn)
	return c, cleanup, nil
}
