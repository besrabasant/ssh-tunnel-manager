package remove

import (
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/urfave/cli/v2"
)

var Cmd cli.Command = cli.Command{
	Name:      "remove",
	Usage:     "Remove a configuration",
	Aliases:   []string{"rm"},
	UsageText: "ssh-tunnel-manager remove <configuration name>",
	Action: func(cCtx *cli.Context) error {
		entryName := cCtx.Args().First()
		if entryName == "" {
			return fmt.Errorf("<configuration name> needed but not provided")
		}

		configdir, err := utils.ResolveDir(cCtx.String(add.ConfigDirFlagName))
		if err != nil {
			return err
		}

		err = configmanager.NewManager(configdir).RemoveConfiguration(entryName)
		if err != nil {
			return fmt.Errorf("couldn't remove configuration %s: %v", entryName, err)
		}

		return nil
	},
}
