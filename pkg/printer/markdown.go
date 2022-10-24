package printer

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/olekukonko/tablewriter"
)

func (p *Printer) printMarkdownTable(secrets map[string]interface{}) error {
	headers, data := p.buildMarkdownTable(secrets)

	table := tablewriter.NewWriter(p.writer)
	table.SetHeader(headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1}) // merge mounts and paths columns
	table.Render()

	return nil
}

func (p *Printer) buildMarkdownTable(secrets map[string]interface{}) ([]string, [][]string) {
	data := [][]string{}
	headers := []string{}

	m := make(map[string]interface{})
	utils.TransformMap("", secrets, &m)

	for _, k := range utils.SortMapKeys(m) {
		v := utils.ToMapStringInterface(m[k])

		switch {
		case p.onlyPaths:
			headers = []string{"path"}

			data = append(data, []string{k})
		case p.onlyKeys:
			headers = []string{"path", "key"}

			for _, j := range utils.SortMapKeys(v) {
				data = append(data, []string{k, j})
			}
		default:
			headers = []string{"path", "key", "value"}

			for _, j := range utils.SortMapKeys(v) {
				data = append(data, []string{k, j, fmt.Sprintf("%v", v[j])})
			}
		}
	}

	return headers, data
}
