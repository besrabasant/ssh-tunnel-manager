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

var AddConfigurationsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a SSH tunnel configuration",
	Long: `
	Add a SSH tunnel configuration.
`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c, cleanup, err := lib.CreateDaemonServiceClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer cleanup()

		data := rpc.TunnelConfig{}

		app := tview.NewApplication()

		formConfig := lib.ConfigurationFormData{
			Title:             "Add configuration",
			PrimaryBtnLabel:   "Add",
			SecondaryBtnLabel: "Cancel",
		}

		addForm := lib.ConfigurationForm(formConfig, app, &data, func(data *rpc.TunnelConfig) {
			addCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			r, err := c.AddConfiguration(addCtx, &rpc.AddOrUpdateConfigurationRequest{Name: data.Name, Data: data})
			if err != nil {
				fmt.Printf("could not execute command: %v", err)
				return
			}

			fmt.Printf("%s", r.GetResult())
		})

		
		if err := app.SetRoot(addForm, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
			panic(err)
		}
	},
}
