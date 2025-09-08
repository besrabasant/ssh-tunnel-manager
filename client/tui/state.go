package tui

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/besrabasant/ssh-tunnel-manager/pkg/configmanager"
)

type State struct {
	App *tview.Application

	ConfigDir string
	Configs   []configmanager.Entry
	Active    map[string]int
	Filter    string

	Pages         *tview.Pages
	Root          *tview.Flex
	Left          *tview.Flex
	Title         *tview.TextView
	FilterBox     *tview.Flex
	FilterInput   *tview.InputField
	List          *tview.List
	Right         *tview.Flex
	DetailPages   *tview.Pages
	DetailInfo    *tview.Table
	DetailSSH     *tview.TextView
	Logs          *tview.TextView
	Status        *tview.TextView
	FilterFocused bool
	Footer        *tview.TextView
	Rebuilding    bool
	quitCh        chan struct{}
}

func NewState() *State {
	s := &State{
		App:    tview.NewApplication(),
		Active: map[string]int{},
		quitCh: make(chan struct{}),
	}
	s.App.SetTitle("SSH Tunnel Manager")
	s.App.EnableMouse(true)
	return s
}

func (s *State) Close() { close(s.quitCh) }

func (s *State) IsActive(name string) bool {
	_, ok := s.Active[name]
	return ok
}

func (s *State) SelectedEntry() (configmanager.Entry, bool) {
	idx := s.List.GetCurrentItem()
	items := s.Filtered()
	if idx >= 0 && idx < len(items) {
		return items[idx], true
	}
	return configmanager.Entry{}, false
}

func (s *State) Filtered() []configmanager.Entry {
	var out []configmanager.Entry
	if strings.TrimSpace(s.Filter) == "" {
		out = append(out, s.Configs...)
	} else {
		q := strings.ToLower(strings.TrimSpace(s.Filter))
		for _, e := range s.Configs {
			h := strings.ToLower(e.Name + " " + e.Server + " " + e.RemoteHost)
			if strings.Contains(h, q) {
				out = append(out, e)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out
}

func (s *State) ReloadData() error {
	cfgs, err := LoadConfigs(s.ConfigDir)
	if err != nil {
		return err
	}
	s.Configs = cfgs
	acts, err := LoadActive()
	if err != nil {
		return err
	}
	m := map[string]int{}
	for _, a := range acts {
		m[a.Name] = a.LocalPort
	}
	s.Active = m
	return nil
}

func IsEscape(ev *tcell.EventKey) bool { return ev.Key() == tcell.KeyEscape }
