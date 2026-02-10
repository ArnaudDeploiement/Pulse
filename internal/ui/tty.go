package ui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsTTY returns true if stdout is connected to a terminal.
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
