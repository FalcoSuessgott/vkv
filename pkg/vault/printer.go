package vault

import (
	"bytes"
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/render"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/juju/ansiterm"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/xlab/treeprint"
)

// CustomPrinter returns a map of all custom printers for that entity.
func (kv *KVSecrets) PrinterFuncs(fOpts *FormatOptions) printer.PrinterFuncMap {
	return printer.PrinterFuncMap{
		"yaml":     kv.PrintYAML(),
		"json":     kv.PrintJSON(),
		"default":  kv.PrintDefault(fOpts),
		"full":     kv.PrintDetailed(fOpts),
		"policy":   kv.PrintPolicy(),
		"template": kv.PrintTemplate(fOpts),
		"export":   kv.PrintExport(),
		"markdown": kv.PrintMarkdown(fOpts),
	}
}

func (kv *KVSecrets) PrintDefault(fOpts *FormatOptions) printer.PrinterFunc {
	return kv.PrintDetailed(&FormatOptions{
		maskSecrets: fOpts.maskSecrets,
		showDiff:    false,
	})
}

func (kv *KVSecrets) PrintDetailed(fOpts *FormatOptions) printer.PrinterFunc {
	return func() ([]byte, error) {
		if err := kv.ComputeDiffChangelog(); err != nil {
			return nil, fmt.Errorf("failed to compute diff changelog: %w", err)
		}

		tree := treeprint.NewWithRoot(utils.ColorBold(kv.Title()))
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
	tree := treeprint.NewWithRoot(utils.ColorBold(kv.SecretName(p)))

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
		tree = treeprint.NewWithRoot(utils.ColorBold(kv.SecretName(p)))

		// if metadata, add it
		if secret[0].Metadata() != "" {
			tree = treeprint.NewWithRoot(utils.ColorBold(fmt.Sprintf("%s {%s}", kv.SecretName(p), secret[0].Metadata())))
		}

		// iterate backwards, so latest secret is first
		for i := len(secret) - 1; i >= 0; i-- {
			s := secret[i]

			// use tabwriter to align the map keys & values and write it to a buffer
			var b bytes.Buffer
			w := ansiterm.NewTabWriter(&b, 0, 0, 1, ' ', 0)
			t := treeprint.NewWithRoot("")

			str := s.String(fOpts.maskSecrets)

			if fOpts.showDiff {
				str = s.DiffString(fOpts.maskSecrets)
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

			tree.AddMetaBranch(utils.ColorBold(s.Title()), t)
		}
	}

	return tree
}

func (kv *KVSecrets) PrintPolicy() printer.PrinterFunc {
	return func() ([]byte, error) {
		var b bytes.Buffer
		w := ansiterm.NewTabWriter(&b, 0, 0, 1, ' ', 0)

		fmt.Fprint(w, utils.ColorBold("PATH\tCREATE\tREAD\tUPDATE\tDELETE\tLIST\tROOT\n"))

		for path := range kv.Secrets {
			cap, err := kv.GetCapabilities(path)
			if err != nil {
				return nil, err
			}

			fmt.Fprintf(w, "%s\t%s", utils.ColorBold(path), cap)
		}

		w.Flush()

		return bytes.TrimSpace(b.Bytes()), nil
	}
}

func (kv *KVSecrets) PrintTemplate(fOpts *FormatOptions) printer.PrinterFunc {
	return func() ([]byte, error) {
		data := make(map[string]interface{})

		for path, secrets := range kv.Secrets {
			// find the latest version that contains secrets
			for i := len(secrets) - 1; i >= 0; i-- {
				if len(secrets[i].Data) > 0 {
					data[path] = secrets[i].Data

					break
				}
			}
		}

		out, err := render.Apply([]byte(fOpts.template), data)
		if err != nil {
			return nil, err
		}

		return bytes.TrimSpace(out), nil
	}
}

func (kv *KVSecrets) PrintExport() printer.PrinterFunc {
	return func() ([]byte, error) {
		var b bytes.Buffer

		for _, secrets := range kv.Secrets {
			data := make(map[string]interface{})

			// find the latest version that contains secrets
			for i := len(secrets) - 1; i >= 0; i-- {
				if len(secrets[i].Data) > 0 {
					data = secrets[i].Data

					break
				}
			}

			if len(data) != 0 {
				// output in export format depending on OS
				for _, k := range utils.SortMapKeys(data) {
					switch os := runtime.GOOS; os {
					case "windows":
						fmt.Fprintf(&b, "set %s='%v'\n", k, data[k])
					case "linux", "darwin":
						fmt.Fprintf(&b, "export %s='%v'\n", k, data[k])
					default:
						return nil, fmt.Errorf("unsupported OS: %s", os)
					}
				}
			}
		}

		return bytes.TrimSpace(b.Bytes()), nil
	}
}

func (kv *KVSecrets) PrintJSON() printer.PrinterFunc {
	return func() ([]byte, error) {
		out, err := utils.ToJSON(kv)

		return bytes.TrimSpace(out), err
	}
}

func (kv *KVSecrets) PrintYAML() printer.PrinterFunc {
	return func() ([]byte, error) {
		out, err := utils.ToYAML(kv)

		return bytes.TrimSpace(out), err
	}
}

func (kv *KVSecrets) PrintMarkdown(fOpts *FormatOptions) printer.PrinterFunc {
	return func() ([]byte, error) {
		var b bytes.Buffer
		header := []string{"path", "key", "value", "version", "metadata", "last update"}

		table := tablewriter.NewWriter(&b)
		table.SetHeader(header)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetAutoMergeCellsByColumnIndex([]int{0})
		table.SetCaption(true, kv.Title())

		paths := lo.Keys(kv.Secrets)

		slices.Sort(paths)

		for _, path := range paths {
			secret := kv.Secrets[path][len(kv.Secrets[path])-1]

			version := fmt.Sprintf("%d", secret.Version)

			// secret is deleted
			if secret.Deleted {
				version = fmt.Sprintf("%s (deleted)", version)
				table.Append([]string{path, "", "", version, secret.Metadata(), secret.DeletionTime.Format(dateFormat)})
			}

			// secret is destroyed
			if secret.Destroyed {
				version = fmt.Sprintf("%s (destroyed)", version)
				table.Append([]string{path, "", "", version, secret.Metadata(), ""})
			}

			// secrets exists
			for _, k := range utils.SortMapKeys(secret.Data) {
				v := fmt.Sprintf("%s", secret.Data[k])

				if fOpts.maskSecrets {
					v = utils.MaskString(v)
				}

				table.Append([]string{path, k, v, version, secret.Metadata(), secret.VersionCreatedTime.Format(dateFormat)})
			}
		}

		table.Render()

		return b.Bytes(), nil
	}
}
