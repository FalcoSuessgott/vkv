package printer

import (
	"fmt"
	"strings"

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
	table.SetAutoMergeCellsByColumnIndex([]int{0}) // merge mounts and paths columns
	table.Render()

	return nil
}

//nolint: gocognit, nestif, cyclop
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

			rootPath := p.enginePath
			subPath := strings.ReplaceAll(k, rootPath, "")

			for i, j := range utils.SortMapKeys(v) {
				d := []string{k, j, fmt.Sprintf("%v", v[j])} // path, key, value

				if i == 0 {
					if p.showVersion {
						headers = append(headers, "version")

						if v, err := p.vaultClient.ReadSecretVersion(rootPath, subPath); err == nil {
							d = append(d, fmt.Sprintf("%v", v)) // version
						}
					}

					if p.showMetadata {
						headers = append(headers, "metadata")

						if v, err := p.vaultClient.ReadSecretMetadata(rootPath, subPath); err == nil {
							m := ""

							md, ok := v.(map[string]interface{})
							if ok {
								for k, v := range md {
									m = fmt.Sprintf("%v %s=%v", m, k, v)
								}
							}

							d = append(d, fmt.Sprintf("%v", strings.TrimPrefix(m, " "))) // metadata
						}
					}
				} else {
					// add empty cells if metadata or versions enabled
					if p.showVersion {
						d = append(d, "")
					}
					if p.showMetadata {
						d = append(d, "")
					}
				}

				data = append(data, d)
			}
		}
	}

	return headers, data
}
