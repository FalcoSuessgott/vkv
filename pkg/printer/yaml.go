package printer

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printYAML(secrets map[string]interface{}) error {
	out, err := utils.ToYAML(secrets)
	if err != nil {
		return err
	}

	fmt.Fprintf(p.writer, "%s\n", string(out))

	return nil
}
