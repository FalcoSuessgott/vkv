package vault

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/markdown"
	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/SerhiiCho/timeago/v2"
	"github.com/savioxavier/termlink"
	"github.com/xlab/treeprint"
)

// CustomPrinter returns a map of all custom printers for that entity.
func (kv *KVSecrets) PrinterFuncs() printer.PrinterFuncMap {
	return printer.PrinterFuncMap{
		"yaml": kv.PrintYAML(),
		"json": kv.PrintJSON(),
		// "default" kv.PrintDefault(),
		"detailed": kv.PrintDetailed(),
		"policy":   kv.PrintPolicy(),
		"template": kv.PrintTemplate(),
		"export":   kv.PrintExport(),
		"markdown": kv.PrintMarkdown(),
	}
}

func (kv *KVSecrets) PrintPolicy() printer.PrinterFunc {
	return func() ([]byte, error) {
		return []byte("policy printer"), nil
	}
}

func (kv *KVSecrets) PrintTemplate() printer.PrinterFunc {
	return func() ([]byte, error) {
		return []byte("template printer"), nil
	}
}

func (kv *KVSecrets) PrintExport() printer.PrinterFunc {
	return func() ([]byte, error) {
		return []byte("export printer"), nil
	}
}

func (kv *KVSecrets) PrintJSON() printer.PrinterFunc {
	return func() ([]byte, error) {
		return utils.ToJSON(kv)
	}
}

func (kv *KVSecrets) PrintYAML() printer.PrinterFunc {
	return func() ([]byte, error) {
		return utils.ToYAML(kv)
	}
}

func (kv *KVSecrets) PrintMarkdown() printer.PrinterFunc {
	return func() ([]byte, error) {
		return markdown.Table([]string{"test"}, [][]string{
			{"ok"},
			{"test"},
		})
	}
}

func (kv *KVSecrets) Title() string {
	addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv", kv.Client.Address(), kv.MountPath)

	return fmt.Sprintf("%s/ [%s] (%s)", termlink.Link(kv.MountPath, addr, false), kv.Type, kv.Description)
}

func (kv *KVSecrets) PrintDetailed() printer.PrinterFunc {
	return func() ([]byte, error) {
		root := kv.MountPath + utils.Delimiter
		tree := treeprint.NewWithRoot(kv.Title())
		secrets := make(map[string]interface{})

		// transform the paths of the secrets by splitting them at "/"
		for path, secret := range kv.Secrets {
			m := utils.PathMap[Secret](path, secret)
			secrets = utils.DeepMergeMaps(secrets, m)
		}

		secret, ok := secrets[root].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("should not happen")
		}

		// iterate through the map and create a branch for each secret
		for _, i := range utils.SortMapKeys(secret) {
			tree.AddBranch(kv.branch(i, secret[i]))
		}

		return tree.Bytes(), nil
	}
}

func (kv *KVSecrets) branch(p string, m interface{}) treeprint.Tree {
	tree := treeprint.NewWithRoot(kv.SecretName(p))

	// here its just a secret path, so we need to go deeper
	if strings.HasSuffix(p, utils.Delimiter) {
		if m, ok := m.(map[string]interface{}); ok {
			for _, i := range utils.SortMapKeys(m) {
				tree.AddBranch(kv.branch(p+i, m[i]))
			}
		}
	}

	// here is the actual secrets
	if secret, ok := m.([]*Secret); ok {
		tree = treeprint.NewWithRoot(kv.SecretName(p))

		// if metadata, add it
		if secret[0].Metadata() != "" {
			tree = treeprint.NewWithRoot(fmt.Sprintf("%s [%s]", kv.SecretName(p), secret[0].Metadata()))
		}

		// iterate backwards, so latest secret is first
		for i := len(secret) - 1; i >= 0; i-- {
			s := secret[i]

			t := treeprint.NewWithRoot("")

			for k, v := range s.Data {
				// handle --only-keys
				if v == "" {
					t.AddNode(k)

					continue
				}

				t.AddNode(fmt.Sprintf("%s=%v", k, v))
			}

			// dont add empty secrets (e.g --only-paths)
			if len(s.Data) == 0 {
				break
			}

			tree.AddMetaBranch(s.Title(), t)
		}
	}

	return tree
}

func (kv *KVSecrets) SecretName(p string) string {
	name := strings.TrimSuffix(p, utils.Delimiter)

	elems := strings.Split(name, utils.Delimiter)
	if len(elems) > 1 {
		name = path.Base(p)
	}

	if !strings.HasSuffix(p, utils.Delimiter) {
		addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", kv.Client.Address(), kv.MountPath, p)

		if len(elems) > 1 {
			addr = fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", kv.Client.Address(), kv.MountPath, url.QueryEscape(p))
		}

		name = termlink.Link(name, addr, false)
	}

	return name
}

func (s Secret) Title() string {
	status := "created"
	tAgo := timeago.Parse(s.VersionCreatedTime)

	if s.DeletionTime.Format("20060102150405") != defaultTimestamp {
		status = "deleted"
		tAgo = timeago.Parse(s.DeletionTime)
	}

	if s.Destroyed {
		status = "destroyed"
	}

	return fmt.Sprintf("Version %d %s %s", s.Version, status, tAgo)
}

func (s Secret) Metadata() string {
	metadata := ""

	for k, v := range s.CustomMetadata {
		metadata += fmt.Sprintf("%s=%s ", k, v)
	}

	return strings.TrimSuffix(metadata, " ")
}

func (kv *KVSecrets) MaskSecrets(length int) printer.SanitizerFunc {
	return func() error {
		if length == -1 {
			return nil
		}

		for _, secrets := range kv.Secrets {
			for _, s := range secrets {
				for k := range s.Data {
					if len(k) > length {
						s.Data[k] = strings.Repeat("*", length)
					} else {
						s.Data[k] = strings.Repeat("*", len(k))
					}
				}
			}
		}

		return nil
	}
}

func (kv *KVSecrets) OnlyKeys() printer.SanitizerFunc {
	return func() error {
		for _, secrets := range kv.Secrets {
			for _, s := range secrets {
				for k := range s.Data {
					s.Data[k] = ""
				}
			}
		}

		return nil
	}
}

func (kv *KVSecrets) OnlyPaths() printer.SanitizerFunc {
	return func() error {
		for _, secrets := range kv.Secrets {
			for _, s := range secrets {
				s.Data = map[string]interface{}{}
			}
		}

		return nil
	}
}
