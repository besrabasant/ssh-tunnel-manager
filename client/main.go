package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50051"

var sshtmCmd = &cobra.Command{
	Use:   "sshtm",
	Short: "Client for interacting with the daemon",
	Long:  `A client for sending commands to the daemon and executing them in the background.`,
}

var listConfigurationsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List configurations (You can use a pattern to only list the configurations that fuzzy match that pattern)",
	Long: `When using it like this "sshtm list prod" it will only list configurations
	that fuzzy match the word "prod". If you have these configurations (client-prod, client1-stage, otherclient-prod), 
	only (client-prod, otherclient-prod) will be displayed.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		searchPattern := ""

		if len(args) > 0 {
			// Use the first argument as the searchPattern if provided
			searchPattern = args[0]
		}

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			fmt.Printf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewDaemonServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		r, err := c.ListConfigurations(ctx, &pb.ListConfigurationsRequest{SearchPattern: searchPattern})
		if err != nil {
			log.Fatalf("could not execute command: %v", err)
		}
		fmt.Printf("%s", r.GetResult())
	},
}

func main() {
	sshtmCmd.AddCommand(listConfigurationsCmd)
	if err := sshtmCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
