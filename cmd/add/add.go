package add

import (
	"os"
	"path/filepath"

	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/urfave/cli/v2"
)

const ConfigDirFlagName = "config-dir"

var Cmd cli.Command = cli.Command{
	Name:  "add",
	Usage: "Add a configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "Name of the configuration",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "description",
			Usage: "Short description of the configuration (optional)",
		},
		&cli.StringFlag{
			Name: "server",
			Usage: "Address of the SSH server to connect to. Can contain the port. " +
				"If not, port 22 will be used by default",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "user",
			Usage:    "Username to use when connecting",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keyFile",
			Usage:    "Path to the private key file for SSH authentication",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "remoteHost",
			Usage:    "The remote host to forward traffic to",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "remotePort",
			Usage:    "The remote port to forward traffic to",
			Required: true,
		},
	},
	Action: func(cCtx *cli.Context) error {
		configdir, err := utils.ResolveDir(cCtx.String(ConfigDirFlagName))
		if err != nil {
			return err
		}
		e := configmanager.Entry{
			Name:        cCtx.String("name"),
			Description: cCtx.String("description"),
			Server:      cCtx.String("server"),
			User:        cCtx.String("user"),
			KeyFile:     cCtx.String("keyFile"),
			RemoteHost:  cCtx.String("remoteHost"),
			RemotePort:  cCtx.Int("remotePort"),
		}
		return configmanager.NewManager(configdir).AddConfiguration(e)
	},
}

var FlagsPredictor = map[string]complete.Predictor{
	"name":        predict.Nothing,
	"description": predict.Nothing,
	"server":      predict.Nothing,
	"user":        predict.Nothing,
	"keyFile":     complete.PredictFunc(getKeyFiles),
	"remoteHost":  predict.Nothing,
	"remotePort":  predict.Nothing,
}

func getKeyFiles(prefix string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	entries, err := os.ReadDir(filepath.Join(home, ".ssh"))
	if err != nil {
		return nil
	}

	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(home, ".ssh", entry.Name()))
		}
	}
	return files
}
