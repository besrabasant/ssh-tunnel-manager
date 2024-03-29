package main

import (
	"log"
	"os"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/cmd/autocomplete"
	"github.com/besrabasant/ssh-tunnel-manager/cmd/edit"
	"github.com/besrabasant/ssh-tunnel-manager/cmd/list"
	"github.com/besrabasant/ssh-tunnel-manager/cmd/remove"
	"github.com/besrabasant/ssh-tunnel-manager/cmd/tunnel"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/posener/complete/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	// App
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "ssh-tunnel-manager",
		Usage:                "An SSH tunnel manager tool with port forwarding capability.",
		Description:          "Save SSH tunnel configurations and start a tunnel with port forwarding using one of the saved configurations.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  add.ConfigDirFlagName,
				Usage: "Path to a directory where configurations are saved",
				Value: configmanager.DefaultConfigDir,
			},
		},
		Commands: []*cli.Command{
			&list.Cmd,
			&add.Cmd,
			&remove.Cmd,
			&edit.Cmd,
			&tunnel.Cmd,
			&autocomplete.Cmd,
		},
	}

	// Autocomplete
	cmd := &complete.Command{
		Sub: map[string]*complete.Command{
			list.Cmd.Name:         {Args: list.Predictor},
			add.Cmd.Name:          {Args: list.Predictor, Flags: add.FlagsPredictor},
			remove.Cmd.Name:       {Args: list.Predictor},
			edit.Cmd.Name:         {Args: list.Predictor},
			tunnel.Cmd.Name:       {Args: list.Predictor},
			autocomplete.Cmd.Name: {},
		},
	}
	cmd.Complete("ssh-tunnel-manager")

	// Run
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
