package printer

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

const (
	exportFmtString = "export %s='%v'\n"
)

var exportMap map[string]interface{}

func (p *Printer) printExport(secrets map[string]interface{}) error {
	exportMap = make(map[string]interface{})

	buildExport(secrets)

	for _, k := range utils.SortMapKeys(exportMap) {
		fmt.Fprintf(p.writer, exportFmtString, k, exportMap[k])
	}

	return nil
}

func buildExport(secrets map[string]interface{}) {
	for _, v := range secrets {
		m, ok := v.(map[string]interface{})
		if ok {
			buildExport(m)
		} else {
			for k, v := range secrets {
				exportMap[k] = v
			}
		}
	}
}
