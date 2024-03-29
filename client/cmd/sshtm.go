package cmd

import "github.com/spf13/cobra"

var SshtmCmd = &cobra.Command{
	Use:   "sshtm",
	Short: "An SSH tunnel manager tool with port forwarding capability.",
	Long:  "Save SSH tunnel configurations and start a tunnel with port forwarding using one of the saved configurations.",
}