package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/config"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ListConfigurationsCmd  = &cobra.Command{
	Use:     "list [search pattern]",
	Aliases: []string{"l", "ls"},
	Short:   "List configurations (You can use a pattern to only list the configurations that fuzzy match that pattern)",
	Long: `
When using it like this "sshtm list prod" it will only list configurations that fuzzy match the word "prod". If you have these configurations (client-prod, client1-stage, otherclient-prod), only (client-prod, otherclient-prod) will be displayed.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		searchPattern := ""

		if len(args) > 0 {
			searchPattern = args[0]
		}

		conn, err := grpc.Dial(config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
