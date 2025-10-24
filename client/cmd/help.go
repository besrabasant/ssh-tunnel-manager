package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const sshtmPreface = `
SSH Tunnel Manager (sshtm) is a comprehensive tool designed to simplify the management of SSH tunnels. It allows you to save configurations for SSH tunnels and utilize those configurations to establish tunnels with port forwarding capabilities. This tool is ideal for users who frequently need to set up secure SSH connections for accessing remote services or securing data transfers.

Features:
- Save Configurations: Permanently store SSH tunnel setups for quick reuse.
- Port Forwarding: Easily set up port forwarding through SSH tunnels, specifying local and remote ports as needed.
- Manage Tunnels: Start, stop, and list active tunnels based on your saved configurations.

Usage:
- 'sshtm' to Run TUI Client.
- 'sshtm add' to save a new tunnel configuration.
- 'sshtm tunnel <configuration name>' to start a tunnel using a saved configuration.
- 'sshtm list' to view all saved configurations.

This tool streamlines the process of setting up and managing secure SSH connections, making it an invaluable resource for developers, system administrators, and any users who regularly work with remote servers.
`

var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Get Help for SSHTM",
	Long:  sshtmPreface,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Then delegate to Cobra's normal help
		root := cmd.Root()
		if len(args) == 0 {
			return root.Help()
		}

		target, _, err := root.Find(args)
		if err != nil || target == nil {
			return fmt.Errorf("unknown help topic: %q", strings.Join(args, " "))
		}
		return target.Help()
	},
}
