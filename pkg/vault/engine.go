package vault

import (
	"fmt"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

const (
	mountEnginePath   = "sys/mounts/%s"
	listSecretEngines = "sys/mounts"
)

// Engines struct that hols all engines key is the namespace.
type Engines map[string][]string

// EnableKV1Engine enables the kvv1 engine at a specified path.
func (v *Vault) EnableKV1Engine(rootPath string) error {
	options := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"path":    rootPath,
			"version": 1,
		},
	}

	_, err := v.Client.Logical().Write(fmt.Sprintf(mountEnginePath, rootPath), options)
	if err != nil {
		return err
	}

	return nil
}

// EnableKV2Engine enables the kvv2 engine at a specified path.
func (v *Vault) EnableKV2Engine(rootPath string) error {
	options := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"path":    rootPath,
			"version": 2,
		},
	}

	_, err := v.Client.Logical().Write(fmt.Sprintf(mountEnginePath, rootPath), options)
	if err != nil {
		return err
	}

	return nil
}

// EnableKV2EngineErrorIfNotForced enables a KVv2 Engine and errors if
// already enabled, unless force is set to true.
func (v *Vault) EnableKV2EngineErrorIfNotForced(force bool, path string) error {
	rootPath, _ := utils.SplitPath(path)

	if len(strings.Split(path, utils.Delimiter)) > 1 {
		//nolint: nilerr
		if err := v.EnableKV2Engine(rootPath); err != nil {
			return nil
		}
	}

	if v.EnableKV2Engine(rootPath) != nil && !force {
		return fmt.Errorf("a secret engine under \"%s\" is already enabled. Use --force for overwriting", rootPath)
	}

	if err := v.DisableKV2Engine(rootPath); err != nil {
		return fmt.Errorf("error disabling secret engine \"%s\": %w", rootPath, err)
	}

	if err := v.EnableKV2Engine(rootPath); err != nil {
		return fmt.Errorf("error enabling secret engine \"%s\": %w", rootPath, err)
	}

	return nil
}

// IsKV2 returns true if the specified path is a KVv2 engine.
func (v *Vault) IsKV2(rootPath string) (bool, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf("sys/mounts/%s", rootPath))
	if err != nil {
		return false, fmt.Errorf("error detecting KV version: %v", err)
	}

	if data.Data["type"].(string) != "kv" {
		return false, fmt.Errorf("%s is not a KV mount", rootPath)
	}

	options := data.Data["options"].(map[string]interface{})

	if options["version"] == "1" {
		return false, nil
	}

	return true, nil
}

// ListKVSecretEngines returns a list of all visible KV secret engines.
func (v *Vault) ListKVSecretEngines(ns string) ([]string, error) {
	v.Client.SetNamespace(ns)

	data, err := v.Client.Logical().Read((listSecretEngines))
	if err != nil {
		return nil, err
	}

	v.Client.ClearNamespace()

	engineList := []string{}

	if data != nil {
		for k, v := range data.Data {
			t, ok := v.(map[string]interface{})["type"]
			if !ok {
				return nil, fmt.Errorf("cannot get type of engine: %s", k)
			}

			if fmt.Sprintf("%v", t) == "kv" {
				engineList = append(engineList, k)
			}
		}

		return engineList, nil
	}

	return nil, fmt.Errorf("could not list secret engines for namespace: \"%s\". Perhaps invalid namespace", ns)
}

// ListAllKVSecretEngines returns a list of all visible KV secret engines.
func (v *Vault) ListAllKVSecretEngines(ns string) (Engines, error) {
	res := make(Engines)

	nsList, err := v.ListAllNamespaces(ns)
	if err != nil {
		return nil, err
	}

	for k, subNS := range nsList {
		engines, err := v.ListKVSecretEngines(k)
		if err != nil {
			return nil, err
		}

		res[k] = engines

		for _, n := range subNS {
			engines, err := v.ListKVSecretEngines(path.Join(k, n))
			if err != nil {
				return nil, err
			}

			res[path.Join(k, n)] = engines
		}
	}

	return res, nil
}
