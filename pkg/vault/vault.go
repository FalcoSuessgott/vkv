package vault

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/hashicorp/vault/api"
)

//nolint: gosec
const (
	mountEnginePath      = "sys/mounts/%s"
	readWriteSecretsPath = "%s/data/%s"
	listSecretsPath      = "%s/metadata/%s"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client
}

// Secrets holds all recursive secrets of a certain path.
type Secrets map[string]interface{}

// NewClient returns a new vault client wrapper.
func NewClient() (*Vault, error) {
	_, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		return nil, fmt.Errorf("VAULT_ADDR required but not set")
	}

	vaultToken, ok := os.LookupEnv("VAULT_TOKEN")
	if !ok {
		return nil, fmt.Errorf("VAULT_TOKEN required but not set")
	}

	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}

	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c.SetToken(vaultToken)

	vaultNamespace, ok := os.LookupEnv("VAULT_NAMESPACE")
	if ok {
		c.SetNamespace(vaultNamespace)
	}

	return &Vault{Client: c}, nil
}

// ListRecursive returns secrets to a path recursive.
func (s *Secrets) ListRecursive(v *Vault, rootPath, subPath string) error {
	keys, err := v.ListSecrets(rootPath, subPath)
	if err != nil {
		// no sub directories in here, but lets check for normal kv pairs then..
		secrets, e := v.ReadSecrets(rootPath, subPath)
		if e == nil {
			(*s)[path.Join(rootPath, subPath)] = secrets

			return nil
		}

		return err
	}

	for _, k := range keys {
		if strings.HasSuffix(k, utils.Delimiter) {
			if err := s.ListRecursive(v, rootPath, path.Join(subPath, k)); err != nil {
				return err
			}
		} else {
			secrets, err := v.ReadSecrets(rootPath, path.Join(subPath, k))
			if err != nil {
				(*s)[path.Join(rootPath, subPath, k)] = map[string]interface{}{}

				continue
			}

			(*s)[path.Join(rootPath, subPath, k)] = secrets
		}
	}

	return nil
}

// ListSecrets returns all keys from vault kv secret path.
func (v *Vault) ListSecrets(rootPath, subPath string) ([]string, error) {
	data, err := v.Client.Logical().List(fmt.Sprintf(listSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets under path \"%s\" found", path.Join(rootPath, subPath))
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

	return nil, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
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

// EnableKV2Engine enables the kv2 engine at a specified path.
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

// DisableKV2Engine disables the kv2 engine at a specified path.
func (v *Vault) DisableKV2Engine(rootPath string) error {
	_, err := v.Client.Logical().Delete(fmt.Sprintf(mountEnginePath, rootPath))
	if err != nil {
		return err
	}

	return nil
}
