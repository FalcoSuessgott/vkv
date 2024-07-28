package vault

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/SerhiiCho/timeago/v2"
	"github.com/r3labs/diff/v3"
	"github.com/samber/lo"
	"github.com/savioxavier/termlink"
)

func (kv *KVSecrets) Title() string {
	return fmt.Sprintf("%s [%s] %s",
		func() string {
			if _, ok := os.LookupEnv(utils.NoHyperlinksEnv); ok {
				return kv.MountPath
			}

			if termlink.SupportsHyperlinks() {
				addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv", kv.Client.Address(), kv.MountPath)

				return termlink.Link(kv.MountPath, addr, false)
			}

			return kv.MountPath
		}(),
		kv.Type,
		func() string {
			if kv.Description != "" {
				return fmt.Sprintf("(%s)", kv.Description)
			}

			return ""
		}(),
	)
}

func (kv *KVSecrets) ComputeDiffChangelog() error {
	for path, secrets := range kv.Secrets {
		// lets prepend an empty secret version as the first secret
		secretVersions := []*Secret{{}}
		secretVersions = append(secretVersions, secrets...)

		for i := range secretVersions {
			if i+1 < len(secretVersions) {
				log, err := diff.Diff(secretVersions[i].Data, secretVersions[i+1].Data)
				if err != nil {
					return err
				}

				kv.Secrets[path][i].Changelog = log
			}
		}
	}

	return nil
}

func (kv *KVSecrets) SecretName(p string) string {
	name := strings.TrimSuffix(p, utils.Delimiter)

	elems := strings.Split(name, utils.Delimiter)
	if len(elems) > 1 {
		name = path.Base(p)
	}

	if !strings.HasSuffix(p, utils.Delimiter) {
		if _, ok := os.LookupEnv(utils.NoHyperlinksEnv); ok {
			return name
		}

		if termlink.SupportsHyperlinks() {
			addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", kv.Client.Address(), kv.MountPath, p)

			if len(elems) > 1 {
				addr = fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", kv.Client.Address(), kv.MountPath, url.QueryEscape(p))
			}

			name = termlink.Link(name, addr, false)
		}

		return name
	}

	return name
}

func (s Secret) Title() string {
	status := "created"
	tAgo := timeago.Parse(s.VersionCreatedTime)

	if s.Deleted {
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

	for _, k := range utils.SortMapKeys(s.CustomMetadata) {
		metadata += fmt.Sprintf("%s=%s ", k, s.CustomMetadata[k])
	}

	return strings.TrimSuffix(metadata, " ")
}

// String returns a string representation of the secret.
func (s *Secret) String(mask bool) string {
	str := ""

	for _, k := range utils.SortMapKeys(s.Data) {
		v := fmt.Sprintf("%s", s.Data[k])

		if mask {
			v = utils.MaskString(v)
		}

		str += fmt.Sprintf("\"%s\"\t= \"%s\"\n", utils.ColorGreen(k), utils.ColorGreen(v))
	}

	return str
}

// DiffString returns a string representing the changes compared to the previous secrets version.
func (s *Secret) DiffString(mask bool) string {
	// if no changelog, secret and previous version match, output the secret
	if s.Changelog == nil || len(s.Changelog) == 0 {
		return s.String(mask)
	}

	var (
		m    = make(map[string]struct{ op, v string })
		keys = []string{}
	)

	// write all changes colored to a map
	for _, change := range s.Changelog {
		keys = append(keys, change.Path[0])

		switch change.Type {
		case diff.CREATE:
			v := fmt.Sprintf("\"%s\"", utils.ColorGreen(change.To.(string)))

			if mask {
				v = fmt.Sprintf("\"%s\"", utils.ColorGreen(utils.MaskString(change.To)))
			}

			m[change.Path[0]] = struct{ op, v string }{
				op: fmt.Sprintf("%s \"%s\"", utils.ColorGreen("[+]"), utils.ColorGreen(change.Path[0])),
				v:  v,
			}
		case diff.UPDATE:
			v := fmt.Sprintf("\"%s\" -> \"%s\"", utils.ColorYellow(change.From.(string)), utils.ColorYellow(change.To.(string)))

			if mask {
				v = fmt.Sprintf("\"%s\" -> \"%s\"",
					utils.ColorYellow(utils.MaskString(change.From)),
					utils.ColorYellow(utils.MaskString(change.To)))
			}

			m[change.Path[0]] = struct{ op, v string }{
				op: fmt.Sprintf("%s \"%s\"", utils.ColorYellow("[~]"), utils.ColorYellow(change.Path[0])),
				v:  v,
			}

		case diff.DELETE:
			v := fmt.Sprintf("\"%s\"", utils.ColorRed(change.From.(string)))

			if mask {
				v = fmt.Sprintf("\"%s\"", utils.ColorRed(utils.MaskString(change.From)))
			}

			m[change.Path[0]] = struct{ op, v string }{
				op: fmt.Sprintf("%s \"%s\"", utils.ColorRed("[-]"), change.Path[0]),
				v:  v,
			}
		}
	}

	// write all other (untouched) keys to the map
	for k, value := range s.Data {
		if !slices.Contains(keys, k) {
			data := struct{ op, v string }{
				op: fmt.Sprintf("\"%s\"", k),
				v:  fmt.Sprintf("\"%s\"", value),
			}

			if mask {
				data.v = fmt.Sprintf("\"%s\"", utils.MaskString(value))
			}

			m[k] = data
		}
	}

	var (
		mapKeys = lo.Keys(m)
		str     = ""
	)

	// output the map in alphabetical order
	slices.Sort(mapKeys)

	for _, k := range mapKeys {
		str += fmt.Sprintf("%s\t= %s\n", m[k].op, m[k].v)
	}

	return str
}
