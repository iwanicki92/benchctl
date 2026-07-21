// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"fmt"
	"io"
	"os"
)

// ProgressLine renders a single status line that updates as work proceeds. On a
// terminal it rewrites the current line in place; otherwise it prints each
// update on its own line so piped output and logs stay readable.
type ProgressLine struct {
	writer io.Writer
	tty    bool
}

func NewProgressLine(writer io.Writer) *ProgressLine {
	return &ProgressLine{writer: writer, tty: isTerminal(writer)}
}

// TTY reports whether the line writes to a terminal, so callers can throttle
// updates when output is piped.
func (line *ProgressLine) TTY() bool {
	return line.tty
}

// Update shows msg as the current progress. On a terminal it returns to the
// start of the line and clears it first, so successive updates overwrite.
func (line *ProgressLine) Update(msg string) {
	if line.tty {
		fmt.Fprintf(line.writer, "\r\x1b[K%s", msg)
		return
	}
	fmt.Fprintln(line.writer, msg)
}

// Done writes the final message and ends the line.
func (line *ProgressLine) Done(msg string) {
	if line.tty {
		fmt.Fprintf(line.writer, "\r\x1b[K%s\n", msg)
		return
	}
	fmt.Fprintln(line.writer, msg)
}

// isTerminal reports whether writer is a character device, i.e. a terminal.
func isTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
