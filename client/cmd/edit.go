package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var EditConfigurationsCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"l", "ls"},
	Short:   "Edit a configuration",
	Long: `
Edit a configuration.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Not  implemented yet ", cmd.Use)
	},
}
