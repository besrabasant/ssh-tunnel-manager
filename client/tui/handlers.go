package tui

import (
	"fmt"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
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

func startOrKillSelected(s *State, e configmanager.Entry) {
	if s.IsActive(e.Name) {
		lp, ok := s.Active[e.Name]
		if !ok || lp <= 0 {
			showError(s, fmt.Errorf("cannot determine active local port for %q", e.Name))
			return
		}
		if err := KillTunnel(e.Name, lp); err != nil {
			showError(s, err)
			return
		}
	} else {
		lp := e.LocalPort
		if lp <= 0 {
			lp = -1 // request auto-assign
		}
		if err := StartTunnel(e.Name, lp); err != nil {
			showError(s, err)
			return
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
