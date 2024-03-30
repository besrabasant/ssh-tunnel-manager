package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/besrabasant/ssh-tunnel-manager/client/lib"
	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var AddConfigurationsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a configuration",
	Long: `
Add a configuration.
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
		defaultFieldWidth := 40
		app := tview.NewApplication()

		editform := tview.NewForm().
			SetFieldBackgroundColor(tcell.ColorDefault).
			SetItemPadding(1).
			AddInputField(
				"Name",
				data.Name,
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					data.Name = text
				},
			).
			AddTextArea("Description", data.Description, 40, 0, 0, func(text string) {
				data.Description = text
			}).
			AddInputField(
				"Server",
				data.Server,
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					data.Server = text
				},
			).
			AddInputField(
				"User",
				data.User,
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					data.User = text
				},
			).
			AddInputField(
				"Key File",
				data.KeyFile,
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					data.KeyFile = text
				},
			).
			AddInputField(
				"Remote Host",
				data.RemoteHost,
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					data.RemoteHost = text
				},
			).
			AddInputField(
				"Remote Port",
				"",
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					port, err := strconv.Atoi(text)
					if err == nil {
						data.RemotePort = int32(port)
					}
				},
			).
			AddButton("Add", func() {
				app.Stop()

				addCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				r, err := c.AddConfiguration(addCtx, &rpc.AddOrUpdateConfigurationRequest{Name: data.Name, Data: &data})
				if err != nil {
					fmt.Printf("could not execute command: %v", err)
					return
				}

				fmt.Printf("%s", r.GetResult())
			}).
			AddButton("Cancel", func() {
				app.Stop()
			})

		editform.SetBorder(true).SetTitle("Edit Configuration").SetTitleAlign(tview.AlignLeft)

		if err := app.SetRoot(editform, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
			panic(err)
		}
	},
}
