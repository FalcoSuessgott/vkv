package vault

import (
	"context"
	"fmt"
	"path"
	"sort"
)

const (
	listNamespaces  = "sys/namespaces"
	createNamespace = "sys/namespaces/%s"
)

// Namespaces represents vault hierarchical namespaces.
type Namespaces map[string][]string

// ListAllNamespaces lists all namespaces of a specified namespace recursively.
func (v *Vault) ListAllNamespaces(ctx context.Context, ns string) (Namespaces, error) {
	m := make(Namespaces)

	//nolint: errcheck
	v.namespaceIterator(ctx, ns, &m)

	if len(m) == 0 {
		m = Namespaces{
			"": []string{},
		}
	}

	return m, nil
}

// nolint: godox
func (v *Vault) namespaceIterator(ctx context.Context, ns string, res *Namespaces) error {
	nsList, err := v.ListNamespaces(ctx, ns)
	if err != nil {
		return err
	}

	if len(nsList) > 0 {
		sort.Strings(nsList)

		(*res)[ns] = nsList

		for _, n := range nsList {
			if err := v.namespaceIterator(ctx, path.Join(ns, n), res); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListNamespaces list the namespaces of the specified namespace.
func (v *Vault) ListNamespaces(ctx context.Context, ns string) ([]string, error) {
	v.Client.SetNamespace(ns)

	data, err := v.Client.Logical().ListWithContext(ctx, listNamespaces)
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
func (v *Vault) DeleteNamespace(ctx context.Context, parentns, ns string) error {
	v.Client.SetNamespace(parentns)

	_, err := v.Client.Logical().DeleteWithContext(ctx, fmt.Sprintf(createNamespace, ns))
	if err != nil {
		return err
	}

	v.Client.SetNamespace("")

	return nil
}

// CreateNamespaceErrorIfNotForced creates a namespace returns no error if force is true.
func (v *Vault) CreateNamespaceErrorIfNotForced(ctx context.Context, parentNS, nsName string, force bool) error {
	v.Client.SetNamespace(parentNS)

	if _, err := v.Client.Logical().WriteWithContext(ctx, fmt.Sprintf(createNamespace, nsName), nil); err != nil {
		if force {
			return nil
		}

		return err
	}

	v.Client.SetNamespace("")

	return nil
}
