package printer

import (
	"fmt"
	"log"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/disiqueira/gotree/v3"
)

func (p *Printer) printBase(secrets map[string]interface{}) error {
	for _, k := range utils.SortMapKeys(secrets) {
		tree := gotree.New(k)
		m := utils.ToMapStringInterface(secrets[k])

		for _, i := range utils.SortMapKeys(m) {
			subMap, ok := m[i].(map[string]interface{})
			if !ok {
				log.Fatalf("cannot convert %T to map[string]interface", m[i])
			}

			tree.AddTree(p.printTree(i, subMap))
		}

		fmt.Fprint(p.writer, tree.Print())
	}

	return nil
}

func (p *Printer) printTree(path string, m map[string]interface{}) gotree.Tree {
	var tree gotree.Tree

	if strings.HasSuffix(path, utils.Delimiter) {
		tree = gotree.New(path)

		for _, i := range utils.SortMapKeys(m) {
			subMap, ok := m[i].(map[string]interface{})
			if !ok {
				log.Fatalf("cannot convert %T to map[string]interface", m[i])
			}

			tree.AddTree(p.printTree(i, subMap))
		}
	} else {
		tree = gotree.New(path)
		for _, k := range utils.SortMapKeys(m) {
			if p.onlyKeys {
				tree.Add(k)
			}

			if !p.onlyKeys && !p.onlyPaths {
				tree.Add(fmt.Sprintf("%s=%v", k, m[k]))
			}
		}
	}

	return tree
}
