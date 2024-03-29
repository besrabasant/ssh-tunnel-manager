package cmd

import "github.com/spf13/cobra"

var SshtmCmd = &cobra.Command{
	Use:   "sshtm",
	Short: "An SSH tunnel manager tool with port forwarding capability.",
	Long:  "Save SSH tunnel configurations and manage tunnels with port forwarding using the saved configurations.",
}
