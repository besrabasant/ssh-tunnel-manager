package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DeleteConfigurationsCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"l", "ls"},
	Short:   "Delete a configuration",
	Long: `
Delete a configuration.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Not  implemented yet ", cmd.Short)
	},
}
