package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var EditConfigurationsCmd = &cobra.Command{
	Use:     "edit <configuration name>",
	Aliases: []string{"e"},
	Short:   "Edit an existing SSH tunnel configuration.",
	Long: `
Edit an existing SSH tunnel configuration interactively.

Use this command to modify the details of a saved SSH tunnel configuration, such as changing the local or remote ports, the SSH server, or any other parameter defined in the configuration. This command provides an interactive interface where you can select the configuration you wish to edit and make changes as required.

The command requires the name of the configuration as an argument. If the configuration name is not provided or is incorrect, the command will prompt for the correct name. After selecting a configuration, you will be guided through a series of prompts to update the desired fields.

Example Usage:
- sshtm edit my_configuration
- sshtm e my_configuration
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configName := ""

		if len(args) == 0 {
			fmt.Println("\n<configuration name> needed but not provided")
			return
		}

		if len(args) > 0 {
			configName = args[0]
		}

		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		r, err := c.FetchConfiguration(ctx, &rpc.FetchConfigurationRequest{Name: configName})
		if err != nil {
			fmt.Printf("could not execute command: %v", err)
			return
		}

		status := r.GetStatus()

		if status == rpc.ResponseStatus_Error {
			fmt.Printf("%s", r.GetMessage())
			return
		}

		data := r.GetData()

		cancel()

		app := tview.NewApplication()

		formConfig := lib.ConfigurationFormData{
			Title:             "Edit configuration",
			PrimaryBtnLabel:   "Update",
			SecondaryBtnLabel: "Cancel",
		}

		editForm := lib.ConfigurationForm(formConfig, app, data, func(data *rpc.TunnelConfig) {
			updateCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			r, err := c.UpdateConfiguration(updateCtx, &rpc.AddOrUpdateConfigurationRequest{Name: configName, Data: data})
			if err != nil {
				fmt.Printf("could not execute command: %v", err)
				return
			}

			fmt.Printf("%s", r.GetResult())
		})

		if err := app.SetRoot(editForm, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
			panic(err)
		}
	},
}
