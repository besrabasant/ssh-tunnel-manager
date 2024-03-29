package list

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/besrabasant/ssh-tunnel-manager/cmd/add"
	"github.com/besrabasant/ssh-tunnel-manager/configmanager"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/posener/complete/v2"
	"github.com/urfave/cli/v2"
)

var Predictor complete.Predictor = complete.PredictFunc(predictConfigurations)

var Cmd cli.Command = cli.Command{
	Name:    "list",
	Aliases: []string{"l", "ls"},
	Usage:   "List configurations (You can use a pattern to only list the configurations that fuzzy match that pattern)",
	Description: `When using it like this "ssh-tunnel-manager list prod" it will only list configurations
that fuzzy match the word "prod". If you have these configurations (client-prod, client1-stage, otherclient-prod), 
only (client-prod, otherclient-prod) will be displayed.`,

	Action: func(cCtx *cli.Context) error {
		cfgs, err := getConfigs()
		if err != nil {
			return err
		}

		output := os.Stdout

		if len(cfgs) == 0 {
			fmt.Fprintf(output, "No configurations found\n")
			return nil
		}

		// The user can filter for certain entries using Fuzzy matching.
		entryPattern := cCtx.Args().First()

		if entryPattern != "" {
			cfgs = configmanager.Entries(cfgs).Filter(func(c *configmanager.Entry) bool {
				return fuzzy.Match(strings.ToLower(entryPattern), strings.ToLower(c.Name))
			})
		}

		for i := range cfgs {

			if i != 0 {
				fmt.Fprintf(output, "\n")
			}

			// config is prented without a new line at its end.
			printConfig(output, cfgs[i])
			fmt.Fprintf(output, "\n")
		}

		return nil
	},
}

func getConfigs() ([]configmanager.Entry, error) {
	dirpath := configmanager.DefaultConfigDir
	if value := os.Getenv(add.ConfigDirFlagName); value != "" {
		dirpath = value
	}

	configdir, err := utils.ResolveDir(dirpath)
	if err != nil {
		return nil, err
	}

	cfgs, err := configmanager.NewManager(configdir).GetConfigurations()
	if err != nil {
		return nil, fmt.Errorf("couldn't get saved configurations: %v", err)
	}

	if len(cfgs) == 0 {
		return nil, nil
	}

	return cfgs, nil
}

func predictConfigurations(prefix string) []string {
	cfgs, err := getConfigs()
	if err != nil {
		return nil
	}
	configs := make([]string, len(cfgs))
	for _, cfg := range cfgs {
		configs = append(configs, cfg.Name)
	}
	return configs
}

func printConfig(w io.Writer, entry configmanager.Entry) {
	template := `%s
  - SSH server:  %s
  - User:        %s
  - Private key: %s
  - Remote:      %s:%d`
	nameAndDesc := bold(entry.Name)
	if strings.TrimSpace(entry.Description) != "" {
		nameAndDesc += " " + "(" + entry.Description + ")"
	}
	w.Write([]byte(
		fmt.Sprintf(
			template,
			nameAndDesc,
			entry.Server,
			entry.User,
			entry.KeyFile,
			entry.RemoteHost,
			entry.RemotePort,
		)))
}

func bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}
