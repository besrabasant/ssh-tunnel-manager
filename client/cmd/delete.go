package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/spf13/cobra"
)

var DeleteConfigurationsCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "Delete a configuration",
	Long: `
Delete a configuration.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configName := ""

		if len(args) == 0 {
			fmt.Println("\n<configuration name> needed but not provided")
			return
		}

		if len(args) > 0 {
			configName = args[0]
		}

		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		r, err := c.DeleteConfiguration(ctx, &rpc.DeleteConfigurationRequest{Name: configName})
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}

		fmt.Printf("%s", r.GetResult())
	},
}
