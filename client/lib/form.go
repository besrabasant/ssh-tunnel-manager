package lib

import (
	"strconv"

	"github.com/besrabasant/ssh-tunnel-manager/rpc"
	"github.com/besrabasant/ssh-tunnel-manager/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfigurationFormData struct {
	Title             string
	PrimaryBtnLabel   string
	SecondaryBtnLabel string
}

func ConfigurationForm(formConfig ConfigurationFormData, app *tview.Application, data *rpc.TunnelConfig, updateFn func(data *rpc.TunnelConfig)) *tview.Form {
	defaultFieldWidth := 40
	form := tview.NewForm().
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
			"SSH Server",
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
			"SSH User",
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
			"Key file Absolute Path",
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
			utils.IntToString(int(data.RemotePort)),
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
		AddInputField(
			"Local Port",
			utils.IntToString(int(data.LocalPort)),
			defaultFieldWidth,
			func(textToCheck string, lastChar rune) bool {
				return true
			},
			func(text string) {
				port, err := strconv.Atoi(text)
				if err == nil {
					data.LocalPort = int32(port)
				}
			},
		).
		AddButton(formConfig.PrimaryBtnLabel, func() {
			app.Stop()

			updateFn(data)
		}).
		AddButton(formConfig.SecondaryBtnLabel, func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle(formConfig.Title).SetTitleAlign(tview.AlignLeft)
	return form
}
