package secret

import "github.com/fatih/color"

// shared styles for tree output. fatih/color automatically disables these when
// the output is not a terminal (pipes, files, tests), so plain text is emitted there.
var (
	// boldStyle highlights path elements (engine, directories, secret names).
	boldStyle = color.New(color.Bold).SprintFunc()
	// versionStyle colors version annotations ("[v=N]", "[Version N created ...]").
	versionStyle = color.New(color.FgCyan).SprintFunc()
	// annotationStyle dims secondary annotations (metadata, engine type/description).
	annotationStyle = color.New(color.FgHiBlack).SprintFunc()
)
