package main

import (
	"fmt"
	"os"

	"github.com/besrabasant/ssh-tunnel-manager/client/cmd"
)

func main() {
	cmd.SshtmCmd.AddCommand(cmd.AddConfigurationsCmd)
	cmd.SshtmCmd.AddCommand(cmd.ListConfigurationsCmd)
	cmd.SshtmCmd.AddCommand(cmd.EditConfigurationsCmd)
	cmd.SshtmCmd.AddCommand(cmd.DeleteConfigurationsCmd)

	cmd.SshtmCmd.AddCommand(cmd.StartSshTunnelCmd)
	cmd.SshtmCmd.AddCommand(cmd.ListActiveSshTunnels)
	cmd.SshtmCmd.AddCommand(cmd.KillSshTunnelCmd)
	
	if err := cmd.SshtmCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
