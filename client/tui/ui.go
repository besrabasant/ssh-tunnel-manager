package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func applyTheme(s *State) {
	th := tview.Styles
	th.PrimitiveBackgroundColor = tcell.ColorDefault
	th.ContrastBackgroundColor = tcell.ColorDefault
	th.MoreContrastBackgroundColor = tcell.ColorDefault
	th.BorderColor = tcell.ColorGray
	th.TitleColor = tcell.ColorLightSlateGray
	th.PrimaryTextColor = tcell.ColorWhite
	th.SecondaryTextColor = tcell.ColorDarkGray
	th.TertiaryTextColor = tcell.ColorGray
	th.InverseTextColor = tcell.ColorBlack
	th.ContrastSecondaryTextColor = tcell.ColorGreen
}

func buildUI(s *State) {
	list := tview.NewList().ShowSecondaryText(true)
	list.SetBorder(true).SetTitle(" Configurations ")
	list.SetSelectedFunc(func(i int, main, secondary string, r rune) {
		if e, ok := s.SelectedEntry(); ok {
			startOrKillSelected(s, e)
		}
	})
	list.SetChangedFunc(func(i int, main, secondary string, r rune) {
		if s.Rebuilding {
			return
		}
		updateDetail(s)
	})
	s.List = list

	filterInput := tview.NewInputField().SetLabel(" Filter: ")
	filterInput.SetFieldWidth(0)
	filterInput.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape:
			if strings.TrimSpace(filterInput.GetText()) != "" {
				filterInput.SetText("")
				s.Filter = ""
				populateList(s)
				updateStatus(s)
				return nil
			}
			s.App.SetFocus(s.List)
			return nil

		// Make the filter "live": use arrows/Page/Tab to jump to the list
		// (the same key press will then act on the list)
		case tcell.KeyDown, tcell.KeyUp, tcell.KeyPgDn, tcell.KeyPgUp, tcell.KeyTab:
			s.App.SetFocus(s.List)
			return ev

		case tcell.KeyEnter:
			s.App.SetFocus(s.List)
			return nil
		}
		return ev
	})

	filterInput.SetChangedFunc(func(text string) {
		s.Filter = text
		populateList(s)
		updateStatus(s)
	})

	filterBox := tview.NewFlex().SetDirection(tview.FlexRow)
	filterBox.AddItem(filterInput, 1, 0, false)
	filterBox.SetBorder(true).SetTitle(" Search ")
	s.FilterInput = filterInput
	s.FilterBox = filterBox

	left := tview.NewFlex().SetDirection(tview.FlexRow)
	left.AddItem(list, 0, 1, true)
	left.AddItem(filterBox, 3, 0, false)
	s.Left = left

	infoTable := tview.NewTable().
		SetBorders(true).
		SetSelectable(false, false)
	s.DetailInfo = infoTable

	sshView := tview.NewTextView().SetDynamicColors(true).SetWrap(true)
	sshView.SetText("SSH session will appear here…")
	s.DetailSSH = sshView

	detailPages := tview.NewPages()
	detailPages.AddPage("info", infoTable, true, true)
	// To-DO: add SSH page later
	detailPages.AddPage("ssh", sshView, true, false)
	s.DetailPages = detailPages

	right := tview.NewFlex().SetDirection(tview.FlexRow)
	right.AddItem(detailPages, 0, 1, false)
	right.SetBorder(true).SetTitle(" Details ")
	s.Right = right

	status := tview.NewTextView().SetDynamicColors(true)
	status.SetBorder(false)
	status.SetTextAlign(tview.AlignLeft)
	status.SetText("     Hints: Up/Down to navigate | Enter: Start/Kill Tunnel | a: add, e: edit, d: delete, r: reload, /: filter, q: quit")
	s.Status = status
	s.Footer = status

	content := tview.NewFlex().SetDirection(tview.FlexColumn)
	content.AddItem(left, 40, 0, true)
	content.AddItem(right, 0, 1, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(content, 0, 1, true)
	root.AddItem(status, 1, 0, false)

	pages := tview.NewPages()
	pages.AddPage("main", root, true, true)

	s.Root = root
	s.Pages = pages
}

func populateList(s *State) {
	s.Rebuilding = true
	s.List.Clear()
	items := s.Filtered()
	for _, e := range items {
		marker := "[gray]○ [-]"
		if s.IsActive(e.Name) {
			marker = "[green]● [-]"
		}
		main := marker + " " + e.Name
		secondary := fmt.Sprintf("%d->%s:%d", e.LocalPort, e.RemoteHost, e.RemotePort)
		cfg := e
		s.List.AddItem(main, secondary, 0, func() { startOrKillSelected(s, cfg) })
	}
	s.Rebuilding = false

	if s.List.GetItemCount() > 0 {
		s.List.SetCurrentItem(0)
	}
	updateDetail(s)
}

func decorateListActive(s *State) {
	selName := ""
	if e, ok := s.SelectedEntry(); ok {
		selName = e.Name
	}

	s.Rebuilding = true
	s.List.Clear()

	items := s.Filtered()
	for _, e := range items {
		marker := "[gray]○ [-]"
		if s.IsActive(e.Name) {
			marker = "[green]● [-]"
		}
		main := marker + " " + e.Name
		secondary := fmt.Sprintf("%d->%s:%d", e.LocalPort, e.RemoteHost, e.RemotePort)
		cfg := e
		s.List.AddItem(main, secondary, 0, func() { startOrKillSelected(s, cfg) })
	}
	s.Rebuilding = false

	if selName != "" {
		selectByName(s, selName)
	}

	updateDetail(s)
}

func updateDetail(s *State) {
	e, ok := s.SelectedEntry()
	if !ok {
		s.DetailInfo.Clear()
		s.DetailInfo.SetCell(0, 0, tview.NewTableCell("No configuration selected.").
			SetAttributes(tcell.AttrBold))
		return
	}

	var statusVal string
	if active := s.IsActive(e.Name); active {
		statusVal = "[green::b]ACTIVE[::-]"
	} else {
		statusVal = "[gray]idle[-]"
	}

	// Collect rows
	rows := [][]string{
		{"Name", e.Name},
		{"Description", strings.TrimSpace(e.Description)},
		{"Status", statusVal},
		{"Server", e.Server},
		{"User", e.User},
		{"KeyFile", e.KeyFile},
		{"RemoteHost", e.RemoteHost},
		{"RemotePort", fmt.Sprintf("%d", e.RemotePort)},
		{"Default LocalPort", fmt.Sprintf("%d", e.LocalPort)},
	}

	// Compute max width for each column
	maxLabel := 0
	maxValue := 0
	for _, r := range rows {
		if len(r[0]) > maxLabel {
			maxLabel = len(r[0])
		}
		if len(r[1]) > maxValue {
			maxValue = len(r[1])
		}
	}

	// Clear and repopulate
	s.DetailInfo.Clear()
	for i, r := range rows {
		s.DetailInfo.SetCell(i, 0, tview.NewTableCell(r[0]).
			SetAlign(tview.AlignRight).
			SetExpansion(0).
			SetMaxWidth(maxLabel).
			SetAttributes(tcell.AttrBold))
		s.DetailInfo.SetCell(i, 1, tview.NewTableCell(r[1]).
			SetAlign(tview.AlignLeft).
			SetExpansion(0).
			SetMaxWidth(maxValue))
	}

	// Make table as tight as possible
	s.DetailInfo.SetBorders(true).
		SetFixed(2, 2).
		SetSelectable(false, false)
}

func updateStatus(s *State) {
	// no-op; footer text is static hints now
}
