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
	Use:     "delete <configuration name>",
	Aliases: []string{"d", "del"},
	Short:   "Delete an existing SSH tunnel configuration.",
	Long: `
Delete an existing SSH tunnel configuration by specifying its configuration name.

This command removes a saved SSH tunnel configuration from your system. Use it when you need to clean up unused configurations or make space for new ones. Deleting a configuration is irreversible, so ensure you have selected the correct one before executing this command.

Simply provide the configuration name of the configuration you wish to delete as an argument to this command. If the configuration name is omitted or incorrect, the command will prompt you to provide a valid configuration name.

Examples:
- sshtm delete my_configuration
- sshtm del old_configuration
- sshtm d my_configuration
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
