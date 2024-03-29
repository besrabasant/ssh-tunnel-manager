package edit

import (
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/editor"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/urfave/cli/v2"
)

var Cmd cli.Command = cli.Command{
	Name:      "edit",
	Usage:     "Edit a configuration in your editor (JSON editing)",
	Aliases:   []string{"e"},
	UsageText: "ssh-tunnel-manager edit <configuration name>",
	Action: func(cCtx *cli.Context) error {
		entryName := cCtx.Args().First()
		if entryName == "" {
			return fmt.Errorf("<configuration name> needed but not provided")
		}

		configdir, err := utils.ResolveDir(cCtx.String(add.ConfigDirFlagName))
		if err != nil {
			return err
		}

		editor := editor.New()

		err = configmanager.NewManager(configdir).EditConfiguration(entryName, editor)
		if err != nil {
			return fmt.Errorf("couldn't edit configuration %s: %v", entryName, err)
		}

		return nil
	},
}
