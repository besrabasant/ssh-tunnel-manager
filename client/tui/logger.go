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

	// Build level tag + timestamp
	tag := "[white]LOG[-]"
	switch level {
	case LevelInfo:
		tag = "[green]INFO[-]"
	case LevelError:
		tag = "[yellow]ERROR[-]"
	case LevelCritical:
		tag = "[red]CRITICAL[-]"
	}
	ts := time.Now().Format("15:04:05")
	prefix := ts + " " + tag + " "

	// Normalize and split incoming text
	raw := fmt.Sprintf(format, args...)
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")

	// Keep only non-empty lines (no blank lines)
	outLines := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		outLines = append(outLines, prefix+ln)
	}
	if len(outLines) == 0 {
		return
	}
	newBlock := strings.Join(outLines, "\n")

	// Append to existing buffer with exactly one newline separation and one trailing newline
	prev := strings.TrimRight(s.Logs.GetText(false), "\n")
	if prev == "" {
		s.Logs.SetText(newBlock + "\n")
	} else {
		s.Logs.SetText(prev + "\n" + newBlock + "\n")
	}
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
