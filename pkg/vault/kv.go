package vault

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

// nolint: gosec
const (
	kvv1ReadWriteSecretsPath = "%s/%s"
	kvv1ListSecretsPath      = "%s/%s"

	kvv2ReadWriteSecretsPath = "%s/data/%s"
	kvv2ListSecretsPath      = "%s/metadata/%s"

	mountDetailsPath = "sys/internal/ui/mounts/%s"
)

// Secrets holds all recursive secrets of a certain path.
type Secrets map[string]interface{}

// ListRecursive returns secrets to a path recursive.
// nolint: cyclop
func (v *Vault) ListRecursive(rootPath, subPath string, skipErrors bool) (*Secrets, error) {
	s := make(Secrets)

	keys, err := v.ListKeys(rootPath, subPath)
	if err != nil {
		// no sub directories in here, but lets check for normal kv pairs then..
		secrets, err := v.ReadSecrets(rootPath, subPath)
		if !skipErrors && err != nil {
			return nil, fmt.Errorf("could not read secrets from %s/%s: %w.\n\nYou can skip this error using --skip-errors", rootPath, subPath, err)
		}

		return (*Secrets)(&secrets), nil
	}

	for _, k := range keys {
		if strings.HasSuffix(k, utils.Delimiter) {
			secrets, err := v.ListRecursive(rootPath, path.Join(subPath, k), skipErrors)
			if err != nil {
				return &s, err
			}

			(s)[k] = secrets
		} else {
			secrets, err := v.ReadSecrets(rootPath, path.Join(subPath, k))
			if !skipErrors && err != nil {
				return nil, err
			}

			// do not exit on errors, just an empty map, so json/yaml export still works
			if skipErrors && secrets == nil {
				secrets = make(Secrets)
			}

			(s)[k] = secrets
		}
	}

	return &s, nil
}

// ListKeys returns all keys from vault kv secret path.
func (v *Vault) ListKeys(rootPath, subPath string) ([]string, error) {
	apiPath := fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath)

	isV1, err := v.IsKVv1(rootPath)
	if err != nil {
		return nil, err
	}

	if isV1 {
		apiPath = fmt.Sprintf(kvv1ListSecretsPath, rootPath, subPath)
	}

	data, err := v.Client.Logical().List(apiPath)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no keys found in \"%s\"", path.Join(rootPath, subPath))
	}

	if data.Data != nil {
		keys := []string{}

		k, ok := data.Data["keys"].([]interface{})
		if !ok {
			log.Fatalf("did not found any keys in %s/%s", rootPath, subPath)
		}

		for _, e := range k {
			keys = append(keys, fmt.Sprintf("%v", e))
		}

		return keys, nil
	}

	return nil, fmt.Errorf("no keys found in \"%s\"", path.Join(rootPath, subPath))
}

// IsKVv1 returns true if the current path is a KVv1 Engine.
func (v *Vault) IsKVv1(rootPath string) (bool, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(mountDetailsPath, rootPath))
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, errors.New("cannot lookup mount type")
	}

	// early versions of Vaults KV engine are of type "generic"
	if data.Data["type"] == "generic" {
		return true, nil
	}

	if opts, ok := data.Data["options"].(map[string]interface{}); ok {
		//nolint: forcetypeassert
		if opts["version"].(string) == "1" {
			return true, nil
		}
	}

	return false, nil
}

// ReadSecrets returns a map with all secrets from a kv engine path.
func (v *Vault) ReadSecrets(rootPath, subPath string) (map[string]interface{}, error) {
	apiPath := fmt.Sprintf(kvv2ReadWriteSecretsPath, rootPath, subPath)

	isV1, err := v.IsKVv1(rootPath)
	if err != nil {
		return nil, err
	}

	if isV1 {
		apiPath = fmt.Sprintf(kvv1ReadWriteSecretsPath, rootPath, subPath)
	}

	data, err := v.Client.Logical().Read(apiPath)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
	}

	if isV1 {
		return data.Data, nil
	}

	if d, ok := data.Data["data"]; ok {
		if m, ok := d.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return nil, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
}

// WriteSecrets writes kv secrets to a specified path.
func (v *Vault) WriteSecrets(rootPath, subPath string, secrets map[string]interface{}) error {
	apiPath := fmt.Sprintf(kvv2ReadWriteSecretsPath, rootPath, subPath)
	options := map[string]interface{}{}

	isV1, err := v.IsKVv1(rootPath)
	if err != nil {
		return err
	}

	if isV1 {
		apiPath = fmt.Sprintf(kvv1ReadWriteSecretsPath, rootPath, subPath)
		options = secrets
	} else {
		options["data"] = secrets
	}

	_, err = v.Client.Logical().Write(apiPath, options)
	if err != nil {
		return err
	}

	return nil
}

// ReadSecretMetadata read the metadata of the secret.
func (v *Vault) ReadSecretMetadata(rootPath, subPath string) (interface{}, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("could not read secret %s metadata", path.Join(rootPath, subPath))
	}

	if d, ok := data.Data["custom_metadata"]; ok {
		return d, nil
	}

	return nil, fmt.Errorf("could not read secret %s metadata", path.Join(rootPath, subPath))
}

// ReadSecretVersion read the version of the secret.
func (v *Vault) ReadSecretVersion(rootPath, subPath string) (interface{}, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("could not read secret %s version", path.Join(rootPath, subPath))
	}

	if d, ok := data.Data["current_version"]; ok {
		return d, nil
	}

	return nil, fmt.Errorf("could not read secret %s version", path.Join(rootPath, subPath))
}

// DisableKV2Engine disables the kv2 engine at a specified path.
func (v *Vault) DisableKV2Engine(rootPath string) error {
	_, err := v.Client.Logical().Delete(fmt.Sprintf(mountEnginePath, rootPath))
	if err != nil {
		return err
	}

	return nil
}
