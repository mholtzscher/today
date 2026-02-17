package output

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

func Configure(disableStyling bool) {
	if disableStyling || !isTTY(os.Stdout) {
		pterm.DisableStyling()
	}
}

func Stdoutln(message string) {
	pterm.Fprintln(os.Stdout, message)
}

func Stderrln(message string) {
	pterm.Fprintln(os.Stderr, message)
}

func Confirm(prompt string) (bool, error) {
	confirmed, err := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show(prompt)
	if err != nil {
		return false, fmt.Errorf("show confirmation prompt: %w", err)
	}

	return confirmed, nil
}

func IsInputTTY() bool {
	return isTTY(os.Stdin)
}

func isTTY(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}

	return stat.Mode()&os.ModeCharDevice != 0
}
