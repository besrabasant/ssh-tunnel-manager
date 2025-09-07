package tui

import (
	"fmt"
	"time"
)

func Run() error {
	s := NewState()
	applyTheme(s)

	dir, err := ResolveConfigDir()
	if err != nil {
		return fmt.Errorf("resolve config dir: %w", err)
	}
	s.ConfigDir = dir

	if err := s.ReloadData(); err != nil {
		return fmt.Errorf("load data: %w", err)
	}

	buildUI(s)
	populateList(s)
	updateStatus(s)
	attachHandlers(s)

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-s.quitCh:
				return
			case <-t.C:
				if acts, err := LoadActive(); err == nil {
					m := map[string]int{}
					for _, a := range acts {
						m[a.Name] = a.LocalPort
					}
					s.App.QueueUpdateDraw(func() { s.Active = m; decorateListActive(s) })
				}
			}
		}
	}()

	s.App.SetRoot(s.Pages, true).SetFocus(s.List)

	return s.App.Run()
}
