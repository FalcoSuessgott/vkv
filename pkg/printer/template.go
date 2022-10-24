package printer

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/render"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printTemplate(secrets map[string]interface{}) error {
	type entry struct {
		Key   string
		Value interface{}
	}

	data := map[string][]entry{}

	for _, k := range utils.SortMapKeys(secrets) {
		m := utils.ToMapStringInterface(secrets[k])

		for _, i := range utils.SortMapKeys(m) {
			subMap, ok := m[i].(map[string]interface{})
			if !ok {
				return fmt.Errorf("cannot convert %T to map[string]interface", m[i])
			}

			entries := []entry{}
			for _, j := range utils.SortMapKeys(subMap) {
				entries = append(entries, entry{Key: j, Value: subMap[j]})
			}

			data[i] = entries
		}
	}

	output, err := render.String([]byte(p.template), data)
	if err != nil {
		return err
	}

	fmt.Fprintln(p.writer, output.String())

	return nil
}
