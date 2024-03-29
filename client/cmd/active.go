package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/spf13/cobra"
)

var ListActiveSshTunnels = &cobra.Command{
	Use:   "active",
	Short: "List active ssh tunnels",
	Long: `
List active ssh tunnels.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		r, err := c.ListActiveTunnels(ctx, &rpc.ListActiveTunnelsRequest{})
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}
		fmt.Printf("%s", r.GetResult())
	},
}
