package cmd

import (
	"os"
	"testing"

	"github.com/fatih/color"
)

// TestMain disables ANSI colors so test assertions compare against plain text,
// regardless of whether the test runner's stdout is a terminal.
func TestMain(m *testing.M) {
	color.NoColor = true

	os.Exit(m.Run())
}
