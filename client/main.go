package main

import (
	"fmt"
	"os"

	"github.com/besrabasant/ssh-tunnel-manager/client/cmd"
)

func main() {
	cmd.SshtmCmd.AddCommand(cmd.ListConfigurationsCmd)
	if err := cmd.SshtmCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
