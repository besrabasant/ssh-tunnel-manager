package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ListActiveSshTunnels = &cobra.Command{
	Use:     "active",
	Aliases: []string{"l", "ls"},
	Short:   "List active ssh tunnels",
	Long: `
List active ssh tunnels.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Not  implemented yet ", cmd.Use)
	},
}
