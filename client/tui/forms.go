package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
)

func showAddForm(s *State)                             { showConfigForm(s, configmanager.Entry{}, false) }
func showEditForm(s *State, entry configmanager.Entry) { showConfigForm(s, entry, true) }

func showConfigForm(s *State, initial configmanager.Entry, editing bool) {
	form, collect := buildConfigForm(initial)
	form.AddButton("Save", func() {
		e, err := collect()
		if err != nil {
			showError(s, err)
			return
		}
		if editing {
			if err := UpdateConfig(s.ConfigDir, e); err != nil {
				showError(s, err)
				return
			}
		} else {
			if err := AddConfig(s.ConfigDir, e); err != nil {
				showError(s, err)
				return
			}
		}
		if err := s.ReloadData(); err != nil {
			showError(s, err)
		}
		populateList(s)
		selectByName(s, e.Name)
		updateStatus(s)
		s.Pages.RemovePage("modal")
	})
	form.AddButton("Cancel", func() { s.Pages.RemovePage("modal") })
	form.SetButtonsAlign(tview.AlignCenter)
	var title string
	if editing {
		title = " Edit configuration "
	} else {
		title = " Add configuration "
	}
	form.SetBorder(true).SetTitle(title)
	addModalPage(s, form)
}

func buildConfigForm(initial configmanager.Entry) (*tview.Form, func() (configmanager.Entry, error)) {
	f := tview.NewForm()
	f.SetItemPadding(1)
	name := tview.NewInputField().SetLabel("Name").SetText(initial.Name)
	desc := tview.NewInputField().SetLabel("Description").SetText(strings.TrimSpace(initial.Description))
	server := tview.NewInputField().SetLabel("Server (host[:port])").SetText(initial.Server)
	user := tview.NewInputField().SetLabel("User").SetText(initial.User)
	key := tview.NewInputField().SetLabel("KeyFile").SetText(initial.KeyFile)
	rhost := tview.NewInputField().SetLabel("RemoteHost").SetText(initial.RemoteHost)
	rport := tview.NewInputField().SetLabel("RemotePort").SetText(intToStr(initial.RemotePort))
	lport := tview.NewInputField().SetLabel("LocalPort (0=auto)").SetText(intToStr(initial.LocalPort))
	for _, it := range []tview.FormItem{name, desc, server, user, key, rhost, rport, lport} {
		f.AddFormItem(it)
	}
	f.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, ev.Modifiers())
		case tcell.KeyDown:
			return tcell.NewEventKey(tcell.KeyTab, 0, ev.Modifiers())
		case tcell.KeyPgUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, ev.Modifiers())
		case tcell.KeyPgDn:
			return tcell.NewEventKey(tcell.KeyTab, 0, ev.Modifiers())
		}
		return ev
	})
	collect := func() (configmanager.Entry, error) {
		rp, _ := strconv.Atoi(strings.TrimSpace(rport.GetText()))
		lp := 0
		if t := strings.TrimSpace(lport.GetText()); t != "" {
			lp, _ = strconv.Atoi(t)
		}
		e := configmanager.Entry{Name: strings.TrimSpace(name.GetText()), Description: strings.TrimSpace(desc.GetText()), Server: strings.TrimSpace(server.GetText()), User: strings.TrimSpace(user.GetText()), KeyFile: strings.TrimSpace(key.GetText()), RemoteHost: strings.TrimSpace(rhost.GetText()), RemotePort: rp, LocalPort: lp}
		if err := e.Validate(); err != nil {
			return e, err
		}
		return e, nil
	}
	return f, collect
}

func showDeleteConfirm(s *State, name string) {
	modal := tview.NewModal().SetText(fmt.Sprintf("Delete configuration '%s'?", name)).AddButtons([]string{"Yes", "No"}).SetDoneFunc(func(_ int, label string) {
		if label == "Yes" {
			if err := DeleteConfig(s.ConfigDir, name); err != nil {
				showError(s, err)
			} else {
				_ = s.ReloadData()
				populateList(s)
				updateStatus(s)
			}
		}
		s.Pages.RemovePage("modal")
	})
	addModalPage(s, modal)
}

func showError(s *State, err error) {
	modal := tview.NewModal().SetText(fmt.Sprintf("Error: %v", err)).AddButtons([]string{"OK"}).SetDoneFunc(func(int, string) { s.Pages.RemovePage("modal") })
	addModalPage(s, modal)
}

func addModalPage(s *State, primitive tview.Primitive) {
	// Center the modal content horizontally & vertically.
	col := tview.NewFlex().SetDirection(tview.FlexColumn)
	col.AddItem(nil, 0, 1, false)
	col.AddItem(primitive, 0, 2, true)
	col.AddItem(nil, 0, 1, false)

	wrapper := tview.NewFlex().SetDirection(tview.FlexRow)
	wrapper.AddItem(nil, 0, 1, false)
	wrapper.AddItem(col, 0, 1, true)
	wrapper.AddItem(nil, 0, 1, false)

	// Close the modal with Escape from anywhere inside it.
	wrapper.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEscape {
			s.Pages.RemovePage("modal")
			return nil
		}
		return ev
	})

	s.Pages.AddPage("modal", wrapper, true, true)
	s.App.SetFocus(primitive)
}

func selectByName(s *State, name string) {
	for i, e := range s.Filtered() {
		if e.Name == name {
			if i < s.List.GetItemCount() {
				s.List.SetCurrentItem(i)
			}
			return
		}
	}
}

func intToStr(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}
