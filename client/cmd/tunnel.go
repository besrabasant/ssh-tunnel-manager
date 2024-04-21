package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/spf13/cobra"
)

var StartSshTunnelCmd = &cobra.Command{
	Use:   "tunnel <configuration name> [local port]",
	Short: "Start an SSH tunnel using a saved configuration, optionally specifying a local port.",
	Long: `
Start an SSH tunnel using a predefined configuration, with the option to specify a local port for forwarding.

This command initiates an SSH tunnel based on a saved configuration you specify by name. It's designed to forward connections from a local port on your machine to a remote destination defined in the configuration. If you do not specify a local port, the system will automatically allocate a random port for forwarding.

When specifying a local port, ensure it is not in use to avoid binding errors. The command provides a seamless way to establish secure SSH connections for various purposes like secure remote access or port forwarding for services.

Usage:
- sshtm tunnel my_configuration: Starts a tunnel using the my_configuration setup. The tunnel will use the local port saved with the configuration if not specified.
- sshtm tunnel my_configuration 8080: Starts a tunnel using the my_configuration setup with local port 8080 explicitly defined for forwarding.
`,
	Args:          cobra.MinimumNArgs(0),
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		configName := ""
		localPortStr := ""

		var localPort int

		if len(args) == 0 {
			fmt.Println("\n<configuration name> needed but not provided")
			return
		}

		if len(args) > 0 {
			configName = args[0]
		}

		if len(args) > 1 {
			localPortStr = args[1]
		}

		// Parse given port
		localPortInt, err := strconv.Atoi(localPortStr)
		if err != nil {
			// If no local port provided set as -1
			localPortInt = -1
		}
		localPort = localPortInt

		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r, err := c.StartTunnel(ctx, &pb.StartTunnelRequest{ConfigName: configName, LocalPort: int32(localPort)})
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}

		fmt.Printf("%s", r.GetResult())
	},
}
