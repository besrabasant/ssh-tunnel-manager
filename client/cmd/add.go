package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AddConfigurationsCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"l", "ls"},
	Short:   "Add a configuration",
	Long: `
Add a configuration.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Not  implemented yet ", cmd.Use)
	},
}
