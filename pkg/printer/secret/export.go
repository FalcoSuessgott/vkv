package secret

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

const (
	exportFmtString = "export %s='%v'\n"
)

var exportMap map[string]interface{}

func (p *Printer) printExport(secrets map[string]interface{}) error {
	exportMap = make(map[string]interface{})

	utils.TransformMap("", secrets, &exportMap)

	for _, path := range utils.SortMapKeys(exportMap) {
		secrets, ok := exportMap[path].(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot convert %v to map[string]interface{}", exportMap[path])
		}

		for _, secretName := range utils.SortMapKeys(secrets) {
			envKey := secretName

			if p.includePath {
				envKey = strings.ReplaceAll(fmt.Sprintf("%s/%s", path, envKey), "/", "_")
			}

			if p.upper {
				envKey = strings.ToUpper(envKey)
			}

			fmt.Fprintf(p.writer, exportFmtString, envKey, secrets[secretName])
		}
	}

	return nil
}
