package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/spf13/cobra"
)

var KillSshTunnelCmd = &cobra.Command{
	Use:     "kill",
	Short:   "Kill a ssh tunnel",
	Long: `
Kill a ssh tunnel.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tunnelIdentifier := ""
		var localPort int

		if len(args) == 0 {
			fmt.Println("\n<configuration name> or <local port> needed but not provided")
			return
		}

		if len(args) > 0 {
			tunnelIdentifier = args[0]
		}

		localPortInt, err := strconv.Atoi(tunnelIdentifier)
		if err == nil {
			tunnelIdentifier = ""
			localPort = localPortInt
		}

		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r, err := c.KillTunnel(ctx, &rpc.KillTunnelRequest{ConfigName: tunnelIdentifier, LocalPort: int32(localPort)})
		
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}
		
		fmt.Printf("%s", r.GetResult())
	},
}
