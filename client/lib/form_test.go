package lib

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestConfigurationFormInputCaptureConsumesCtrlC(t *testing.T) {
	app := tview.NewApplication()
	capture := configurationFormInputCapture(app)

	if event := capture(tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)); event != nil {
		t.Fatalf("expected Ctrl+C to be consumed, got %#v", event)
	}
}

func TestConfigurationFormInputCapturePassesThroughOtherKeys(t *testing.T) {
	app := tview.NewApplication()
	capture := configurationFormInputCapture(app)
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	if got := capture(event); got != event {
		t.Fatalf("expected non-Ctrl+C event to pass through unchanged")
	}
}
