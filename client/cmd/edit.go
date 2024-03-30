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

var EditConfigurationsCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"l", "ls"},
	Short:   "Edit a configuration",
	Long: `
Edit a configuration.
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

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
				fmt.Sprint(data.RemotePort),
				defaultFieldWidth,
				func(textToCheck string, lastChar rune) bool {
					return true
				},
				func(text string) {
					port, err := strconv.Atoi(text)
					if err != nil {
						data.RemotePort = int32(port)
					}
				},
			).
			AddButton("Update", func() {
				
				app.Stop()

				r, err := c.UpdateConfiguration(ctx, &rpc.AddOrUpdateConfigurationRequest{Name: configName, Data: data})
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
