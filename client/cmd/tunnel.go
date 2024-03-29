package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	pb "github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/spf13/cobra"
)

/* Release 1.0 */

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

		if localPortStr == "" {
			// Generate random port
			randomPort, err := generateRandomPort()
			if err != nil {
				fmt.Printf("couldn't generate a random port: %v", err)
				return
			}
			localPort = randomPort
		} else {
			// Parse given port
			localPortInt, err := strconv.Atoi(localPortStr)
			if err != nil {
				fmt.Printf("provided local port %q is not a valid port", localPortStr)
				return
			}
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
		
		r, err := c.StartTunnel(ctx, &pb.StartTunnelRequest{ConfigName: configName, LocalPort: int32(localPort)})
		
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}
		
		fmt.Printf("%s", r.GetResult())
	},
}

func generateRandomPort() (int, error) {
	// Listen on port 0 to bind to a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	// Extract the port number from the listener address
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, err
	}

	// Convert the port number to an integer
	randomPort, err := net.LookupPort("tcp", port)
	if err != nil {
		return 0, err
	}

	// Close the listener
	err = listener.Close()
	if err != nil {
		return 0, err
	}

	return randomPort, nil
}
