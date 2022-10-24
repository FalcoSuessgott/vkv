package printer

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printJSON(secrets map[string]interface{}) error {
	out, err := utils.ToJSON(secrets)
	if err != nil {
		return err
	}

	fmt.Fprintf(p.writer, "%s\n", string(out))

	return nil
}
