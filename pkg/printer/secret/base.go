package secret

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/savioxavier/termlink"
	"github.com/xlab/treeprint"
)

func (p *Printer) printBase(secrets map[string]interface{}) error {
	var tree treeprint.Tree

	m := make(map[string]interface{})

	for _, k := range utils.SortMapKeys(secrets) {
		baseName := p.enginePath

		if p.withHyperLinks {
			addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv", p.vaultClient.Client.Address(), p.enginePath)

			baseName = termlink.Link(p.enginePath, addr, false)
		}

		if p.vaultClient != nil {
			// append description
			desc, err := p.vaultClient.GetEngineDescription(p.enginePath)
			if err == nil && desc != "" {
				baseName = fmt.Sprintf("%s [desc=%s]", baseName, desc)
			}

			// append type + version
			engineType, version, err := p.vaultClient.GetEngineTypeVersion(p.enginePath)
			if err == nil {
				baseName = fmt.Sprintf("%s [type=%s]", baseName, engineType+version)
			}
		}

		tree = treeprint.NewWithRoot(baseName)

		m = utils.ToMapStringInterface(secrets[k])
	}

	for _, i := range utils.SortMapKeys(m) {
		//nolint: forcetypeassert
		tree.AddBranch(p.printTree(p.enginePath, i, m[i].(map[string]interface{})))
	}

	fmt.Fprintln(p.writer, strings.TrimSpace(tree.String()))

	return nil
}

func (p *Printer) printTree(rootPath, subPath string, m map[string]interface{}) treeprint.Tree {
	tree := treeprint.NewWithRoot(p.buildTreeName(rootPath, subPath))

	//nolint: nestif
	if strings.HasSuffix(subPath, utils.Delimiter) {
		for _, i := range utils.SortMapKeys(m) {
			if data, ok := m[i].(map[string]interface{}); ok {
				tree.AddBranch(p.printTree(rootPath, subPath+i, data))
			}
		}
	} else {
		for _, k := range utils.SortMapKeys(m) {
			if p.onlyKeys {
				tree.AddNode(k)
			}

			if !p.onlyKeys && !p.onlyPaths {
				tree.AddNode(fmt.Sprintf("%s=%v", k, m[k]))
			}
		}
	}

	return tree
}

// nolint: cyclop
func (p *Printer) buildTreeName(rootPath, subPath string) string {
	name := strings.TrimSuffix(subPath, utils.Delimiter)

	subPathParts := strings.Split(strings.TrimSuffix(subPath, utils.Delimiter), utils.Delimiter)
	if len(subPathParts) > 1 {
		name = path.Base(subPath)
	}

	if p.withHyperLinks && !strings.HasSuffix(subPath, utils.Delimiter) {
		// {{ vault-addr }}/ui/vault/secrets/{{ root path }}/kv/{{ sub path (url encoded if contains "/" )}}
		addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), rootPath, subPath)

		if len(subPathParts) > 1 {
			addr = fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), rootPath, url.QueryEscape(subPath))
		}

		name = termlink.Link(name, addr, false)
	}

	if p.showVersion {
		if v, err := p.vaultClient.ReadSecretVersion(p.ctx, rootPath, subPath); err == nil {
			name = fmt.Sprintf("%s [v=%v]", name, v)
		}
	}

	if p.showMetadata {
		if v, err := p.vaultClient.ReadSecretMetadata(p.ctx, rootPath, subPath); err == nil {
			md := ""
			metadata, ok := v.(map[string]interface{})

			if ok {
				for k, v := range metadata {
					md = fmt.Sprintf("%s %s=%v", md, k, v)
				}

				name = fmt.Sprintf("%s [%v]", name, strings.TrimPrefix(md, " "))
			}
		}
	}

	return name
}
