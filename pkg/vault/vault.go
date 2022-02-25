package vault

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	Client  *api.Client
	Secrets map[string]interface{}
}

// NewClient returns a new vault client wrapper.
// VAULT_ADDR and VAULT_TOKEN are required
// VAULT_SKIP_VERIFY is considered, if defined
// reads the proxy configuration via HTTP_PROXY and HTTPS_PROXY.
func NewClient() (*Vault, error) {
	client := &http.Client{}

	vaultAddr, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		return nil, fmt.Errorf("VAULT_ADDR required but not set")
	}

	vaultToken, ok := os.LookupEnv("VAULT_TOKEN")
	if !ok {
		return nil, fmt.Errorf("VAULT_TOKEN required but not set")
	}

	_, skipVerify := os.LookupEnv("VAULT_SKIP_VERIFY")
	if skipVerify {
		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			//nolint: gosec
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	config := &api.Config{
		Address:    vaultAddr,
		HttpClient: client,
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

	return &Vault{Client: c, Secrets: make(map[string]interface{})}, nil
}

// ListRecursive returns secrets to a path recursive.
func (v *Vault) ListRecursive(rootPath, subPath string) error {
	keys, err := v.ListSecrets(rootPath, subPath)
	if err != nil {
		return err
	}

	for _, k := range keys {
		if strings.HasSuffix(k, "/") {
			if err := v.ListRecursive(rootPath, filepath.Join(subPath, k)); err != nil {
				return err
			}
		} else {
			secrets, err := v.ReadSecrets(rootPath, filepath.Join(subPath, k))
			if err != nil {
				return err
			}

			v.Secrets[filepath.Join(rootPath, subPath, k)] = secrets
		}
	}

	return nil
}

// ListSecrets returns all keys from vault kv secret path.
func (v *Vault) ListSecrets(rootPath, subPath string) ([]string, error) {
	data, err := v.Client.Logical().List(fmt.Sprintf(listSecretsPath, rootPath, strings.ReplaceAll(subPath, "\\", "/")))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets found")
	}

	if data.Data != nil {
		keys := []string{}

		for _, k := range data.Data["keys"].([]interface{}) {
			keys = append(keys, k.(string))
		}

		return keys, nil
	}

	return nil, fmt.Errorf("no secrets found")
}

// ReadSecrets returns a map with all secrets from a kv engine path.
func (v *Vault) ReadSecrets(rootPath, subPath string) (map[string]interface{}, error) {
	data, err := v.Client.Logical().Read(fmt.Sprintf(readWriteSecretsPath, rootPath, strings.ReplaceAll(subPath, "\\", "/")))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets found")
	}

	if d, ok := data.Data["data"]; ok {
		if m, ok := d.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return nil, fmt.Errorf("no secrets found")
}

// WriteSecrets writes kv secrets to a specified path.
func (v *Vault) WriteSecrets(rootPath, subPath string, secrets map[string]interface{}) error {
	options := map[string]interface{}{
		"data": secrets,
	}

	_, err := v.Client.Logical().Write(fmt.Sprintf(readWriteSecretsPath, rootPath, strings.ReplaceAll(subPath, "\\", "/")), options)
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
