package vault

import (
	"fmt"
	"path"
	"sort"

	"github.com/FalcoSuessgott/vkv/pkg/markdown"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (ns *Namespaces) PrintJSON() ([]byte, error) {
	return utils.ToJSON(ns)
}

func (ns *Namespaces) PrintYAML() ([]byte, error) {
	return utils.ToYAML(ns)
}

func (ns *Namespaces) PrintMarkdown() ([]byte, error) {
	out, err := markdown.Table([]string{"test"}, [][]string{
		{"ok"},
		{"test"},
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (ns *Namespaces) PrintBase() ([]byte, error) {
	return nil, nil
}

// ListAllNamespaces lists all namespaces of a specified namespace recursively.
func (v *Vault) ListAllNamespaces(ns string) (Namespaces, error) {
	m := make(Namespaces)

	//nolint: errcheck
	v.namespaceIterator(ns, &m)

	if len(m) == 0 {
		m = Namespaces{
			"": []string{},
		}
	}

	return m, nil
}

// nolint: godox
func (v *Vault) namespaceIterator(ns string, res *Namespaces) error {
	nsList, err := v.ListNamespaces(ns)
	if err != nil {
		return err
	}

	if len(nsList) > 0 {
		sort.Strings(nsList)

		(*res)[ns] = nsList

		for _, n := range nsList {
			if err := v.namespaceIterator(path.Join(ns, n), res); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListNamespaces list the namespaces of the specified namespace.
func (v *Vault) ListNamespaces(ns string) ([]string, error) {
	v.Client.SetNamespace(ns)

	data, err := v.Client.Logical().List(listNamespaces)
	if err != nil {
		return nil, err
	}

	res := []string{}

	if data != nil {
		nsList, ok := data.Data["key_info"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot list namespaces for ns %s", ns)
		}

		for k := range nsList {
			res = append(res, k)
		}
	}

	v.Client.SetNamespace("")

	return res, nil
}

// DeleteNamespace deletes a namespace.
func (v *Vault) DeleteNamespace(parentns, ns string) error {
	v.Client.SetNamespace(parentns)

	_, err := v.Client.Logical().Delete(fmt.Sprintf(createNamespace, ns))
	if err != nil {
		return err
	}

	v.Client.SetNamespace("")

	return nil
}

// CreateNamespaceErrorIfNotForced creates a namespace returns no error if force is true.
func (v *Vault) CreateNamespaceErrorIfNotForced(parentNS, nsName string, force bool) error {
	v.Client.SetNamespace(parentNS)

	if _, err := v.Client.Logical().Write(fmt.Sprintf(createNamespace, nsName), nil); err != nil {
		if force {
			return nil
		}

		return err
	}

	v.Client.SetNamespace("")

	return nil
}
