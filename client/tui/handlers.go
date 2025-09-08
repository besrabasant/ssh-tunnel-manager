package tui

import (
	"fmt"
	"strings"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
	"github.com/gdamore/tcell/v2"
)

func attachHandlers(s *State) {
	s.Root.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEnter && s.App.GetFocus() == s.List {
			if e, ok := s.SelectedEntry(); ok {
				startOrKillSelected(s, e)
				return nil
			}
			return nil
		}
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
	var msg string
	var err error

	if s.IsActive(e.Name) {
		msg, err = KillTunnel(e.Name, 0) // let daemon find the port
		if err != nil {
			showError(s, err)
			s.LogError("Kill failed for %s: %v", e.Name, err)
			return
		}
		s.LogInfo("Killed tunnel %s", e.Name)
	} else {
		lp := e.LocalPort
		if lp <= 0 {
			lp = -1 // request auto-assign
		}
		msg, err = StartTunnel(e.Name, lp)
		if err != nil {
			showError(s, err)
			s.LogError("Start failed for %s (requested port %d): %v", e.Name, lp, err)
			return
		}
		s.LogInfo("Started tunnel %s (requested port %d)", e.Name, lp)
	}

	// Refresh active list & status as before
	acts, err := LoadActive()
	if err != nil {
		showError(s, err)
		s.LogError("Active refresh error: %v", err)
		return
	}
	m := map[string]int{}
	for _, a := range acts {
		m[a.Name] = a.LocalPort
	}
	s.Active = m
	decorateListActive(s)
	updateStatus(s)

	// Push daemon/TM output to logs
	if strings.TrimSpace(msg) == "" {
		msg = "No output from daemon."
	}
	s.LogInfo("Daemon: %s", strings.TrimSpace(msg))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
