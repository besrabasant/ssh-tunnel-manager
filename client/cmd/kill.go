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
	Use:     "kill <configuration name or local port>",
	Aliases: []string{"terminate"},
	Short:   "Terminate an active SSH tunnel.",
	Long: `
Terminate an active SSH tunnel either by specifying its configuration name or the local port it uses.

This command allows you to forcefully close an active SSH tunnel. You can specify the tunnel by its configuration name or the local port number that the tunnel uses. This is particularly useful for managing resources or ending tunnels that are no longer required, are malfunctioning, or for security purposes.

The command requires either a configuration name or a local port number as an argument. If neither is provided, the command will prompt you to enter one of them. Ensure you correctly identify the tunnel to avoid accidentally terminating the wrong connection.

Example Usage:
- sshtm kill my_configuration
- sshtm kill 8080
- sshtm terminate my_configuration
- sshtm terminate 8080
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
