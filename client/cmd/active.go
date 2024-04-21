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
	Short: "List all active ssh tunnels",
	Long: `
List active SSH tunnels.

This command queries the SSH Tunnel Manager daemon to retrieve a list of all currently active SSH tunnels. An active SSH tunnel is one that is currently established and facilitating data transmission. This is useful for monitoring or managing ongoing SSH connections, providing insights into which tunnels are active and potentially consuming system resources.

Use this command to ensure that your SSH tunnels are running as expected, or to diagnose issues related to network connections established through SSH tunneling.	
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
