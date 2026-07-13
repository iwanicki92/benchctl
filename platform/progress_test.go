package platform

import (
	"bytes"
	"testing"
)

func TestProgressLineTerminalUpdatesInPlace(t *testing.T) {
	var buf bytes.Buffer
	line := &ProgressLine{writer: &buf, tty: true}

	line.Update("50%")
	line.Update("75%")
	line.Done("100%")

	want := "\r\x1b[K50%\r\x1b[K75%\r\x1b[K100%\n"
	if got := buf.String(); got != want {
		t.Errorf("progress on a terminal = %q, want %q", got, want)
	}
}

func TestProgressLineNonTerminalPrintsLines(t *testing.T) {
	var buf bytes.Buffer
	line := &ProgressLine{writer: &buf, tty: false}

	line.Update("50%")
	line.Update("75%")
	line.Done("100%")

	want := "50%\n75%\n100%\n"
	if got := buf.String(); got != want {
		t.Errorf("progress without a terminal = %q, want %q", got, want)
	}
}
