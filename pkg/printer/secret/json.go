package secret

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printJSON(secrets map[string]interface{}) error {
	out, err := utils.ToJSON(secrets)
	if err != nil {
		return err
	}

	fmt.Fprint(p.writer, string(out))

	return nil
}
