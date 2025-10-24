package cmd

import (
	"github.com/besrabasant/ssh-tunnel-manager/client/tui"
	"github.com/spf13/cobra"
)

// SshtmCmd launches the tview-based TUI.
var SshtmCmd = &cobra.Command{
	Use:   "sshtm",
	Short: "Manage SSH tunnels with ease, including saving configurations and setting up port forwarding.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}
