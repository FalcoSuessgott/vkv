package secret

import (
	"fmt"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/xlab/treeprint"
)

func (p *Printer) printBase(secrets map[string]interface{}) error {
	var tree treeprint.Tree

	m := make(map[string]interface{})

	for _, k := range utils.SortMapKeys(secrets) {
		baseName := p.enginePath + utils.Delimiter

		// if p.withHyperLinks {
		// 	addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv", p.vaultClient.Client.Address(), p.enginePath)

		// 	baseName = termlink.Link(p.enginePath+utils.Delimiter, addr, false)
		// }

		// if p.vaultClient != nil {
		// 	// append description
		// 	desc, err := p.vaultClient.GetEngineDescription(p.enginePath)
		// 	if err == nil && desc != "" {
		// 		baseName = fmt.Sprintf("%s [desc=%s]", baseName, desc)
		// 	}

		// 	// append type + version
		// 	engineType, version, err := p.vaultClient.GetEngineTypeVersion(p.enginePath)
		// 	if err == nil {
		// 		baseName = fmt.Sprintf("%s [type=%s]", baseName, engineType+version)
		// 	}
		// }

		tree = treeprint.NewWithRoot(baseName)

		m = utils.ToMapStringInterface(secrets[k])
	}

	for _, i := range utils.SortMapKeys(m) {
		//nolint: forcetypeassert
		tree.AddBranch(p.printTree(p.enginePath, i, m[i]))
	}

	fmt.Fprintln(p.writer, strings.TrimSpace(tree.String()))

	return nil
}

func (p *Printer) printTree(rootPath, subPath string, m interface{}) treeprint.Tree {
	tree := treeprint.NewWithRoot(buildTreeName(rootPath, subPath))

	// here its just a secret path, so we need to go deeper
	if strings.HasSuffix(subPath, utils.Delimiter) {
		if m, ok := m.(map[string]interface{}); ok {
			for _, i := range utils.SortMapKeys(m) {
				cm, ok := m[i].(map[string]interface{})

				if ok && cm["CustomMetadata"] != nil {
					for k, v := range cm["CustomMetadata"].(map[string]interface{}) {
						fmt.Println(k, v)
					}
				}

				tree.AddBranch(p.printTree(rootPath, subPath+i, m[i]))
			}
		}
	}

	// here is the actual secrets
	if _, ok := m.(map[string]interface{}); !ok {
		list := m.([]interface{})

		for i := len(list) - 1; i >= 0; i-- {
			s, ok := list[i].(map[string]interface{})
			if !ok {
				continue
			}

			status := ""

			if _, ok := s["Destroyed"].(bool); ok {
				status = "DESTROYED"
			}

			if _, ok := s["Deleted"].(bool); ok {
				status = "Deleted"
			}

			f := treeprint.NewWithRoot(fmt.Sprintf("%s destroyed %s (%s ago)", status, s["VersionCreatedTime"], utils.TimeAgo(s["VersionCreatedTime"].(string))))

			if s["Data"] != nil {
				for k, v := range s["Data"].(map[string]interface{}) {
					f.AddNode(fmt.Sprintf("%s=%v", k, v))
				}
			}

			tree.AddMetaBranch(fmt.Sprintf("v%v", i+1), f)
		}
	}

	return tree
}

// nolint: cyclop
func buildTreeName(rootPath, subPath string) string {
	name := strings.TrimSuffix(subPath, utils.Delimiter)

	subPathParts := strings.Split(strings.TrimSuffix(subPath, utils.Delimiter), utils.Delimiter)
	if len(subPathParts) > 1 {
		name = path.Base(subPath)
	}

	// if p.withHyperLinks && !strings.HasSuffix(subPath, utils.Delimiter) {
	// 	// {{ vault-addr }}/ui/vault/secrets/{{ root path }}/kv/{{ sub path (url encoded if contains "/" )}}
	// 	addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), rootPath, subPath)

	// 	if len(subPathParts) > 1 {
	// 		addr = fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), rootPath, url.QueryEscape(subPath))
	// 	}

	// 	name = termlink.Link(name, addr, false)
	// }

	// if p.showVersion {
	// 	if v, err := p.vaultClient.ReadSecretVersion(rootPath, subPath); err == nil {
	// 		name = fmt.Sprintf("%s [v=%v]", name, v)
	// 	}
	// }

	// if p.showMetadata {
	// 	if v, err := p.vaultClient.ReadSecretMetadata(rootPath, subPath); err == nil {
	// 		md := ""
	// 		metadata, ok := v.(map[string]interface{})

	// 		if ok {
	// 			for k, v := range metadata {
	// 				md = fmt.Sprintf("%s %s=%v", md, k, v)
	// 			}

	// 			name = fmt.Sprintf("%s [%v]", name, strings.TrimPrefix(md, " "))
	// 		}
	// 	}
	// }

	return name
}
