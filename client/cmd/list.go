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
	Short:   "List all SSH tunnel configurations, optionally filtering by a search pattern. (You can use a pattern to only list the configurations that fuzzy match that pattern)",
	Long: `
List all SSH tunnel configurations, with optional fuzzy matching based on a provided search pattern.

This command displays all saved SSH tunnel configurations. If a search pattern is provided as an argument, the command will perform a fuzzy search and only list configurations that match the search term. This feature is useful for quickly finding specific configurations among many, especially when you remember only part of the configuration name or details.

For instance, using 'sshtm list prod' will filter and display configurations that contain the substring 'prod' anywhere in their names or related fields, such as 'client-prod' and 'otherclient-prod'. This makes it easy to manage and access configurations in environments with numerous setups.

Example Usage:
- sshtm list (Lists all configurations)
- sshtm list prod (Lists configurations that contain 'prod')
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
