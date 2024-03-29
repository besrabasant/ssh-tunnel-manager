package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"

	"github.com/spf13/cobra"
)

var ListConfigurationsCmd = &cobra.Command{
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

		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		r, err := c.ListConfigurations(ctx, &pb.ListConfigurationsRequest{SearchPattern: searchPattern})
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}
		fmt.Printf("%s", r.GetResult())
	},
}
