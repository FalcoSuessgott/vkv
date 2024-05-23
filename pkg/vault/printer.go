package vault

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/markdown"
	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/juju/ansiterm"
	"github.com/xlab/treeprint"
)

// CustomPrinter returns a map of all custom printers for that entity.
func (kv *KVSecrets) PrinterFuncs(fOpts *FormatOptions) printer.PrinterFuncMap {
	return printer.PrinterFuncMap{
		"yaml":     kv.PrintYAML(),
		"json":     kv.PrintJSON(),
		"default":  kv.PrintDetailed(fOpts),
		"policy":   kv.PrintPolicy(),
		"template": kv.PrintTemplate(),
		"export":   kv.PrintExport(),
		"markdown": kv.PrintMarkdown(),
	}
}

func (kv *KVSecrets) PrintDetailed(fOpts *FormatOptions) printer.PrinterFunc {
	return func() ([]byte, error) {
		// prepare data before computing print output
		if fOpts.ShowDiff {
			if err := kv.ComputeDiffChangelog(); err != nil {
				return nil, fmt.Errorf("failed to compute diff changelog: %w", err)
			}
		}

		if fOpts.OnlyKeys {
			kv.OnlyKeys()
		}

		tree := treeprint.NewWithRoot(kv.Title())
		secrets := make(map[string]interface{})

		// transform the paths of the secrets by splitting them at "/"
		for path, secret := range kv.Secrets {
			m := utils.UnflattenMap(path, secret)
			secrets = utils.DeepMergeMaps(secrets, m)
		}

		secret, ok := secrets[kv.MountPath].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("should not happen")
		}

		// iterate through the map and create a branch for each secret
		for _, i := range utils.SortMapKeys(secret) {
			tree.AddBranch(kv.branch(i, secret[i], fOpts))
		}

		// idk why but the treeprint lib adds a bunch of empty newlines at the end
		return bytes.TrimSpace(tree.Bytes()), nil
	}
}

func (kv *KVSecrets) branch(p string, m interface{}, fOpts *FormatOptions) treeprint.Tree {
	tree := treeprint.NewWithRoot(kv.SecretName(p))

	// here its just a secret path, so we need to go deeper
	if strings.HasSuffix(p, utils.Delimiter) {
		if m, ok := m.(map[string]interface{}); ok {
			for _, i := range utils.SortMapKeys(m) {
				tree.AddBranch(kv.branch(p+i, m[i], fOpts))
			}
		}
	}

	// here is the actual secrets
	if secret, ok := m.([]*Secret); ok {
		tree = treeprint.NewWithRoot(kv.SecretName(p))

		// if metadata, add it
		if secret[0].Metadata() != "" {
			tree = treeprint.NewWithRoot(fmt.Sprintf("%s {%s}", kv.SecretName(p), secret[0].Metadata()))
		}

		// iterate backwards, so latest secret is first
		for i := len(secret) - 1; i >= 0; i-- {
			s := secret[i]

			// use tabwriter to align the map keys & values and write it to a buffer
			var b bytes.Buffer
			w := ansiterm.NewTabWriter(&b, 0, 0, 1, ' ', 0)
			t := treeprint.NewWithRoot("")

			str := s.String(fOpts.MaskSecrets, fOpts.MaxValueLength)
			if fOpts.ShowDiff {
				str = s.DiffString(fOpts.OnlyKeys, fOpts.MaskSecrets, fOpts.MaxValueLength)
			}

			fmt.Fprintln(w, str)

			// write to buffer
			w.Flush()

			// and then split the content at new line char and add to tree
			for _, i := range strings.Split(b.String(), "\n") {
				if i != "" {
					t.AddNode(i)
				}
			}

			tree.AddMetaBranch(s.Title(), t)
		}
	}

	return tree
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
		var b bytes.Buffer

		for _, secrets := range kv.Secrets {
			// get the latest secret
			for k, v := range secrets[len(secrets)-1].Data {
				b.Write([]byte(fmt.Sprintf("export %s='%v'\n", k, v)))
			}
		}

		return b.Bytes(), nil
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
		header := []string{"secret", "key", "value", "version", "metadata"}
		rows := [][]string{}

		for path, secrets := range kv.Secrets {
			for i := len(secrets) - 1; i >= 0; i-- {
				s := secrets[i]

				for k, v := range s.Data {
					rows = append(rows, []string{path, k, fmt.Sprintf("%v", v), fmt.Sprintf("%d", s.Version), s.Metadata()})
				}

				// append an empty row for better readability
				rows = append(rows, make([]string, len(header)))
			}
		}
		return markdown.Table(header, rows)
	}
}
