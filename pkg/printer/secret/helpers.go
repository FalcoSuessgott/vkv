package secret

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (p *Printer) printOnlykeys(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if ok {
			res[k] = p.printOnlykeys(m)
		} else {
			res[k] = ""
		}
	}

	return res
}

func (p *Printer) printOnlyPaths(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if ok {
			res[k] = p.printOnlyPaths(m)
		} else {
			res[k] = nil
		}
	}

	return res
}

func (p *Printer) printMergePaths(secrets map[string]interface{}, root string) map[string]interface{} {
	m := make(map[string]interface{})

	utils.FlattenMap(secrets, m, root)

	return m
}

func (p *Printer) maskValues(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if ok {
			res[k] = p.maskValues(m)
		} else {
			n := fmt.Sprintf("%v", v)
			if len(n) > p.valueLength && p.valueLength != -1 {
				secrets[k] = strings.Repeat(maskChar, p.valueLength)
			} else {
				secrets[k] = strings.Repeat(maskChar, len(n))
			}
		}
	}

	return secrets
}
