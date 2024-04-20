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
	Use:           "tunnel <configuration name>",
	Short:         "Start a tunnel using a configuration. The tunnel will forward connections to [local port] if specified or to a random port.",
	Long:          "ssh-tunnel-manager tunnel <configuration name> [local port]",
	Args:          cobra.MinimumNArgs(1),
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

