package secret

import (
	"fmt"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/disiqueira/gotree/v3"
)

func (p *Printer) printBase(enginePath string, secrets map[string]interface{}) error {
	var tree gotree.Tree

	m := make(map[string]interface{})

	for _, k := range utils.SortMapKeys(secrets) {
		tree = gotree.New(enginePath + utils.Delimiter)
		m = utils.ToMapStringInterface(secrets[k])
	}

	for _, i := range utils.SortMapKeys(m) {
		//nolint: forcetypeassert
		tree.AddTree(p.printTree(enginePath, i, m[i].(map[string]interface{})))
	}

	fmt.Fprint(p.writer, tree.Print())

	return nil
}

func (p *Printer) printTree(rootPath, subPath string, m map[string]interface{}) gotree.Tree {
	tree := gotree.New(p.buildTreeName(rootPath, subPath))

	if strings.HasSuffix(subPath, utils.Delimiter) {
		for _, i := range utils.SortMapKeys(m) {
			//nolint: forcetypeassert
			tree.AddTree(p.printTree(rootPath, subPath+i, m[i].(map[string]interface{})))
		}
	} else {
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

func (p *Printer) buildTreeName(rootPath, subPath string) string {
	name := subPath

	subPathParts := strings.Split(strings.TrimSuffix(subPath, utils.Delimiter), utils.Delimiter)
	if len(subPathParts) > 1 {
		name = path.Base(subPath)
	}

	if p.showVersion {
		if v, err := p.vaultClient.ReadSecretVersion(rootPath, subPath); err == nil {
			name = fmt.Sprintf("v%v: %s", v, name)
		}
	}

	if p.showMetadata {
		if v, err := p.vaultClient.ReadSecretMetadata(rootPath, subPath); err == nil {
			md := ""
			metadata, ok := v.(map[string]interface{})

			if ok {
				for k, v := range metadata {
					md = fmt.Sprintf("%v %s=%v", md, k, v)
				}

				name = fmt.Sprintf("%s [%v]", name, strings.TrimPrefix(md, " "))
			}
		}
	}

	return name
}
