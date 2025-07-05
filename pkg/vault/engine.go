package vault

import (
	"context"
	"fmt"
	"path"
)

const (
	mountEnginePath   = "sys/mounts/%s"
	listSecretEngines = "sys/mounts"
)

// Engines struct that hols all engines key is the namespace.
type Engines map[string][]string

// GetEngineDescription returns the description of the engine.
func (v *Vault) GetEngineDescription(ctx context.Context, rootPath string) (string, error) {
	data, err := v.Client.Logical().ReadWithContext(ctx, fmt.Sprintf(mountEnginePath, rootPath))
	if err != nil {
		return "", err
	}

	if data != nil {
		desc, ok := data.Data["description"]
		if !ok {
			return "", nil
		}

		//nolint: forcetypeassert
		return desc.(string), nil
	}

	return "", fmt.Errorf("could not get engine description for path: \"%s\"", rootPath)
}

func (v *Vault) GetEngineTypeVersion(ctx context.Context, rootPath string) (string, string, error) {
	data, err := v.Client.Logical().ReadWithContext(ctx, fmt.Sprintf(mountEnginePath, rootPath))
	if err != nil {
		return "", "", err
	}

	if data != nil {
		var eType, eVersion string

		t, ok := data.Data["type"]

		if ok {
			//nolint: forcetypeassert
			eType = t.(string)
		}

		// "options" is nil in "generic" type
		//nolint: goconst
		if eType == "generic" {
			return "kv", "1", nil
		}

		v, ok := data.Data["options"]
		if ok {
			//nolint: forcetypeassert
			eVersion = v.(map[string]interface{})["version"].(string)
		}

		return eType, eVersion, nil
	}

	return "", "", fmt.Errorf("could not get engine type for path: \"%s\"", rootPath)
}

// EnableKV2Engine enables the kv2 engine at a specified path.
func (v *Vault) EnableKV2Engine(ctx context.Context, rootPath string) error {
	options := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"path":    rootPath,
			"version": 2,
		},
	}

	_, err := v.Client.Logical().WriteWithContext(ctx, fmt.Sprintf(mountEnginePath, rootPath), options)
	if err != nil {
		return err
	}

	return nil
}

// EnableKV1Engine enables the kv1 engine at a specified path.
func (v *Vault) EnableKV1Engine(ctx context.Context, rootPath string) error {
	options := map[string]interface{}{
		"type": "kv",
		"options": map[string]interface{}{
			"path":    rootPath,
			"version": 1,
		},
	}

	_, err := v.Client.Logical().WriteWithContext(ctx, fmt.Sprintf(mountEnginePath, rootPath), options)
	if err != nil {
		return err
	}

	return nil
}

// EnableKV2EngineErrorIfNotForced enables a KVv2 Engine and errors if
// already enabled, unless force is set to true.
func (v *Vault) EnableKV2EngineErrorIfNotForced(ctx context.Context, force bool, path string) error {
	// check if engine exists
	engineType, kvVersion, err := v.GetEngineTypeVersion(ctx, path)
	// engine does not exists, so we enable it and exit
	if err != nil {
		if err := v.EnableKV2Engine(ctx, path); err != nil {
			return fmt.Errorf("error enabling secret engine \"%s\": %w", path, err)
		}

		return nil
	}

	// engine exists, but is not of type kvv2
	if err == nil && (engineType != "kv" || kvVersion != "2") {
		return fmt.Errorf("engine \"%s\" is not of type kv2", path)
	}

	// engine exists but no force flag used for using that engine
	if err == nil && !force {
		return fmt.Errorf("a secret engine under \"%s\" is already enabled. Use --force for overwriting", path)
	}

	return nil
}

// ListKVSecretEngines returns a list of all visible KV secret engines.
func (v *Vault) ListKVSecretEngines(ctx context.Context, ns string) ([]string, error) {
	v.Client.SetNamespace(ns)

	data, err := v.Client.Logical().ReadWithContext(ctx, listSecretEngines)
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

			if fmt.Sprintf("%v", t) == "kv" || fmt.Sprintf("%v", t) == "generic" {
				engineList = append(engineList, k)
			}
		}

		return engineList, nil
	}

	return nil, fmt.Errorf("could not list secret engines for namespace: \"%s\". Perhaps invalid namespace", ns)
}

// ListAllKVSecretEngines returns a list of all visible KV secret engines.
func (v *Vault) ListAllKVSecretEngines(ctx context.Context, ns string) (Engines, error) {
	res := make(Engines)

	nsList, err := v.ListAllNamespaces(ctx, ns)
	if err != nil {
		return nil, err
	}

	for k, subNS := range nsList {
		engines, err := v.ListKVSecretEngines(ctx, k)
		if err != nil {
			return nil, err
		}

		res[k] = engines

		for _, n := range subNS {
			engines, err := v.ListKVSecretEngines(ctx, path.Join(k, n))
			if err != nil {
				return nil, err
			}

			res[path.Join(k, n)] = engines
		}
	}

	return res, nil
}
