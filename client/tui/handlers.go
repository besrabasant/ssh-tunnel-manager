package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func attachHandlers(s *State) {
	s.Root.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {

		switch ev.Key() {
		case tcell.KeyCtrlC:
			s.Close()
			s.App.Stop()
			return nil
		}
		// If the filter is focused, treat keys as text (do not run global binds).
		// Navigation + ESC are already handled in filterInput.SetInputCapture above.
		if s.App.GetFocus() == s.FilterInput {
			return ev
		}
		switch ev.Rune() {
		case 'q':
			s.Close()
			s.App.Stop()
			return nil
		case 'r':
			if err := s.ReloadData(); err != nil {
				showError(s, fmt.Errorf("reload failed: %w", err))
				return nil
			}
			populateList(s)
			updateStatus(s)
			return nil
		case '/':
			s.App.SetFocus(s.FilterInput)
			s.FilterFocused = true
			return nil
		case 'a':
			showAddForm(s)
			return nil
		case 'e':
			if e, ok := s.SelectedEntry(); ok {
				showEditForm(s, e)
			}
			return nil
		case 'd':
			if e, ok := s.SelectedEntry(); ok {
				showDeleteConfirm(s, e.Name)
			}
			return nil
		case 'g':
			s.List.SetCurrentItem(0)
			updateDetail(s)
			return nil
		case 'G':
			s.List.SetCurrentItem(max(0, s.List.GetItemCount()-1))
			updateDetail(s)
			return nil
		}
		return ev
	})
}

func startOrKillSelected(s *State, name string) {
	if s.IsActive(name) {
		if err := KillTunnel(name); err != nil {
			showError(s, err)
		}
	} else {
		if err := StartTunnel(name); err != nil {
			showError(s, err)
		}
	}
	acts, err := LoadActive()
	if err != nil {
		showError(s, err)
		return
	}
	m := map[string]int{}
	for _, a := range acts {
		m[a.Name] = a.LocalPort
	}
	s.Active = m
	decorateListActive(s)
	updateStatus(s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
