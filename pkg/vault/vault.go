package vault

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client
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

	return &Vault{c}, nil
}

// ListPath returns all keys from vault kv secret path.
func (v *Vault) ListPath(rootPath, subPath string) ([]string, error) {
	path := fmt.Sprintf("%s/metadata/%s", rootPath, subPath)

	data, err := v.Client.Logical().List(path)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets found")
	}

	keys := []string{}

	if data.Data != nil {
		for _, k := range data.Data["keys"].([]interface{}) {
			keys = append(keys, k.(string))
		}
	}

	return keys, nil
}

// ReadSecrets returns a map with all secrets from a kv engine path.
func (v *Vault) ReadSecrets(rootPath, subPath string) (map[string]interface{}, error) {
	path := fmt.Sprintf("%s/data/%s", rootPath, subPath)

	data, err := v.Client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if data.Data == nil {
		return nil, fmt.Errorf("no secrets found")
	}

	return data.Data["data"].(map[string]interface{}), nil
}
