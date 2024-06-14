package secret

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/render"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printTemplate(secrets map[string]interface{}) error {
	m := make(map[string]interface{})
	utils.FlattenMap(secrets, m, "")

	output, err := render.Apply([]byte(p.template), m)
	if err != nil {
		return err
	}

	fmt.Fprintln(p.writer, strings.TrimSpace(output.String()))

	return nil
}
