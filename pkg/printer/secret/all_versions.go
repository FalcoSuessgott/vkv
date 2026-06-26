package secret

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/savioxavier/termlink"
	"github.com/xlab/treeprint"
)

// printBaseAllVersions renders all versions of all secrets in a tree, listing the
// key-value pairs contained in each individual version.
func (p *Printer) printBaseAllVersions(vs vault.VersionedSecrets) error {
	tree := treeprint.NewWithRoot(p.allVersionsRootName())

	// tracks already created directory branches, keyed by their path prefix
	branches := map[string]treeprint.Tree{"": tree}

	for _, fullPath := range sortedKeys(vs) {
		secret := vs[fullPath]
		parts := strings.Split(fullPath, utils.Delimiter)

		// walk/create the intermediate directory branches
		parent := tree
		prefix := ""

		for _, dir := range parts[:len(parts)-1] {
			if prefix == "" {
				prefix = dir
			} else {
				prefix = prefix + utils.Delimiter + dir
			}

			branch, ok := branches[prefix]
			if !ok {
				branch = parent.AddBranch(dir)
				branches[prefix] = branch
			}

			parent = branch
		}

		// secret branch, optionally annotated with custom metadata
		secretBranch := parent.AddBranch(p.allVersionsSecretName(parts[len(parts)-1], fullPath, secret.CustomMetadata))

		for _, sv := range secret.Versions {
			versionBranch := secretBranch.AddBranch(p.versionTitle(sv))

			// deleted/destroyed versions have no retrievable data
			if sv.Data == nil {
				continue
			}

			for _, k := range utils.SortMapKeys(sv.Data) {
				if p.onlyKeys {
					versionBranch.AddNode(k)

					continue
				}

				versionBranch.AddNode(fmt.Sprintf("%s=%s", k, p.formatVersionValue(sv.Data[k])))
			}
		}
	}

	fmt.Fprintln(p.writer, strings.TrimSpace(tree.String()))

	return nil
}

// printJSONAllVersions serializes all secret versions (with real values) as JSON.
func (p *Printer) printJSONAllVersions(vs vault.VersionedSecrets) error {
	out, err := utils.ToJSON(vs)
	if err != nil {
		return err
	}

	fmt.Fprint(p.writer, string(out))

	return nil
}

// printYAMLAllVersions serializes all secret versions (with real values) as YAML.
func (p *Printer) printYAMLAllVersions(vs vault.VersionedSecrets) error {
	out, err := utils.ToYAML(vs)
	if err != nil {
		return err
	}

	fmt.Fprint(p.writer, string(out))

	return nil
}

// allVersionsRootName builds the tree root label, e.g. "secret/ [kv2] (key/value secret storage)".
func (p *Printer) allVersionsRootName() string {
	display := p.enginePath

	if p.withHyperLinks && p.vaultClient != nil {
		addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv", p.vaultClient.Client.Address(), p.enginePath)
		display = termlink.Link(p.enginePath, addr, false)
	}

	baseName := boldStyle(display)

	if p.vaultClient != nil {
		if engineType, version, err := p.vaultClient.GetEngineTypeVersion(p.ctx, p.enginePath); err == nil {
			baseName = fmt.Sprintf("%s %s", baseName, annotationStyle(fmt.Sprintf("[%s]", engineType+version)))
		}

		if desc, err := p.vaultClient.GetEngineDescription(p.ctx, p.enginePath); err == nil && desc != "" {
			baseName = fmt.Sprintf("%s %s", baseName, annotationStyle(fmt.Sprintf("(%s)", desc)))
		}
	}

	return baseName
}

// allVersionsSecretName builds the secret node label, with custom metadata and an optional hyperlink.
func (p *Printer) allVersionsSecretName(name, fullPath string, customMetadata map[string]interface{}) string {
	if p.withHyperLinks && p.vaultClient != nil {
		addr := fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), p.enginePath, fullPath)
		if strings.Contains(fullPath, utils.Delimiter) {
			addr = fmt.Sprintf("%s/ui/vault/secrets/%s/kv/%s", p.vaultClient.Client.Address(), p.enginePath, url.QueryEscape(fullPath))
		}

		name = termlink.Link(name, addr, false)
	}

	name = boldStyle(name)

	if len(customMetadata) > 0 {
		parts := make([]string, 0, len(customMetadata))
		for _, k := range utils.SortMapKeys(customMetadata) {
			parts = append(parts, fmt.Sprintf("%s=%v", k, customMetadata[k]))
		}

		name = fmt.Sprintf("%s %s", name, annotationStyle(fmt.Sprintf("{%s}", strings.Join(parts, " "))))
	}

	return name
}

// formatVersionValue masks or trims a value according to the printer options.
func (p *Printer) formatVersionValue(v interface{}) string {
	s := fmt.Sprintf("%v", v)

	if p.showValues {
		return s
	}

	if p.valueLength != -1 && len(s) > p.valueLength {
		return strings.Repeat(maskChar, p.valueLength)
	}

	return strings.Repeat(maskChar, len(s))
}

// versionTitle returns the label for a version branch, e.g. "[Version 2 created 5 minutes ago]".
func (p *Printer) versionTitle(sv *vault.SecretVersion) string {
	now := p.now
	if now.IsZero() {
		now = time.Now()
	}

	status, ref := "created", sv.CreatedTime

	switch {
	case sv.Destroyed:
		status = "destroyed"

		if sv.DeletionTime != nil {
			ref = *sv.DeletionTime
		}
	case sv.DeletionTime != nil:
		status = "deleted"
		ref = *sv.DeletionTime
	}

	return versionStyle(fmt.Sprintf("[Version %d %s %s]", sv.Version, status, humanizeTimeAgo(ref, now)))
}

// humanizeTimeAgo renders a timestamp as a relative duration, e.g. "58 minutes ago".
func humanizeTimeAgo(t, now time.Time) string {
	if t.IsZero() {
		return "unknown time"
	}

	d := now.Sub(t)
	if d < 0 {
		d = 0
	}

	switch {
	case d < time.Minute:
		return pluralizeAgo(int(d.Seconds()), "second")
	case d < time.Hour:
		return pluralizeAgo(int(d.Minutes()), "minute")
	case d < 24*time.Hour:
		return pluralizeAgo(int(d.Hours()), "hour")
	case d < 30*24*time.Hour:
		return pluralizeAgo(int(d.Hours()/24), "day")
	case d < 365*24*time.Hour:
		return pluralizeAgo(int(d.Hours()/(24*30)), "month")
	default:
		return pluralizeAgo(int(d.Hours()/(24*365)), "year")
	}
}

func pluralizeAgo(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s ago", unit)
	}

	return fmt.Sprintf("%d %ss ago", n, unit)
}

// sortedKeys returns the secret paths of a VersionedSecrets map in lexical order.
func sortedKeys(vs vault.VersionedSecrets) []string {
	keys := make([]string, 0, len(vs))
	for k := range vs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
