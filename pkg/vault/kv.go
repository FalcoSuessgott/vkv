package vault

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

const (
	//nolint: gosec
	readWriteSecretsPath = "%s/data/%s"
	//nolint: gosec
	listSecretsPath = "%s/metadata/%s"
)

// Secrets holds all recursive secrets of a certain path.
type Secrets map[string]interface{}

// ListRecursive returns secrets to a path recursive.
func (v *Vault) ListRecursive(rootPath, subPath string) (*Secrets, error) {
	s := make(Secrets)

	keys, err := v.ListKeys(rootPath, subPath)
	if err != nil {
		// no sub directories in here, but lets check for normal kv pairs then..
		secrets, err := v.ReadSecrets(rootPath, subPath)
		if err != nil {
			return nil, fmt.Errorf("could not read secrets from %s/%s: %w", rootPath, subPath, err)
		}

		return (*Secrets)(&secrets), nil
	}

	for _, k := range keys {
		if strings.HasSuffix(k, utils.Delimiter) {
			secrets, err := v.ListRecursive(rootPath, path.Join(subPath, k))
			if err != nil {
				return &s, err
			}

			(s)[k] = secrets
		} else {
			secrets, err := v.ReadSecrets(rootPath, path.Join(subPath, k))
			if err != nil {
				return nil, err
			}

			(s)[k] = secrets
		}
	}

	return &s, nil
}

// ListKeys returns all keys from vault kv secret path.
func (v *Vault) ListKeys(rootPath, subPath string) ([]string, error) {
	data, err := v.Client.Logical().List(fmt.Sprintf(listSecretsPath, rootPath, subPath))
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

// ReadSecrets returns a map with all secrets from a kv engine path.
func (v *Vault) ReadSecrets(rootPath, subPath string) (map[string]interface{}, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(readWriteSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
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
	options := map[string]interface{}{
		"data": secrets,
	}

	_, err := v.Client.Logical().Write(fmt.Sprintf(readWriteSecretsPath, rootPath, subPath), options)
	if err != nil {
		return err
	}

	return nil
}

// ReadSecretMetadata read the metadata of the secret.
func (v *Vault) ReadSecretMetadata(rootPath, subPath string) (interface{}, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(listSecretsPath, rootPath, subPath))
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
	data, err := v.Client.Logical().Read(fmt.Sprintf(listSecretsPath, rootPath, subPath))
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
