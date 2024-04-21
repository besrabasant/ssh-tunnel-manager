package cmd

import (
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/config"
	"github.com/spf13/cobra"
)


var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SSH Tunnel Manager",
	Long:  `All software has versions. This is SSH Tunnel Manager's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SSH Tunnel Manager Version " + config.AppVersion)
	},
}