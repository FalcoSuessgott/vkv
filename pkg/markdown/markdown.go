package markdown

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// TablePrinter prints the specified data in a markdown table.
func Table(header []string, rows [][]string) ([]byte, error) {
	var w bytes.Buffer

	table := tablewriter.NewWriter(&w)

	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(rows)
	table.SetAutoMergeCellsByColumnIndex([]int{0}) // merge mounts and paths columns

	table.Render()

	return w.Bytes(), nil
}
