package tui

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelError
	LevelCritical
)

func (s *State) log(level LogLevel, format string, args ...any) {
	if s == nil || s.Logs == nil {
		return
	}
	ts := time.Now().Format("15:04:05")
	msg := strings.TrimSpace(fmt.Sprintf(format, args...))

	var tag string
	switch level {
	case LevelInfo:
		tag = "[green]INFO[-]"
	case LevelError:
		tag = "[yellow]ERROR[-]"
	case LevelCritical:
		tag = "[red]CRITICAL[-]"
	default:
		tag = "[white]LOG[-]"
	}

	prev := s.Logs.GetText(false)
	if prev != "" && !strings.HasSuffix(prev, "\n") {
		prev += "\n"
	}
	s.Logs.SetText(prev + fmt.Sprintf("%s %s %s\n", ts, tag, msg))

}

// Public helpers
func (s *State) LogInfo(format string, args ...any)     { s.log(LevelInfo, format, args...) }
func (s *State) LogError(format string, args ...any)    { s.log(LevelError, format, args...) }
func (s *State) LogCritical(format string, args ...any) { s.log(LevelCritical, format, args...) }

// Optional: clear logs
func (s *State) ClearLogs() {
	if s == nil || s.Logs == nil {
		return
	}
	s.Logs.SetText("")
}
