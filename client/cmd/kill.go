package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var KillSshTunnelCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{"l", "ls"},
	Short:   "Kill a ssh tunnel",
	Long: `
Kill a ssh tunnel.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Not  implemented yet ", cmd.Use)
	},
}
