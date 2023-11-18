package vault

import (
	"fmt"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/exec"
	"github.com/hashicorp/vault/api"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client
}

// NewDefaultClient returns a new vault client wrapper.
func NewDefaultClient() (*Vault, error) {
	_, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		return nil, fmt.Errorf("VAULT_ADDR required but not set")
	}

	vaultToken, tokenExported := os.LookupEnv("VAULT_TOKEN")

	cmd, ok := os.LookupEnv("VKV_LOGIN_COMMAND")
	if !tokenExported && ok {
		cmdParts := strings.Split(cmd, " ")

		token, err := exec.Run(cmdParts)
		if err != nil {
			return nil, fmt.Errorf("error running VKV_LOGIN_CMD (%s): %w", cmd, err)
		}

		vaultToken = strings.TrimSpace(string(token))
	}

	if vaultToken == "" {
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

	_, err = c.Auth().Token().Lookup(vaultToken)
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Perhaps not a valid token")
	}

	return &Vault{Client: c}, nil
}

// NewClient returns a new vault client wrapper.
func NewClient(addr, token string) (*Vault, error) {
	cfg := &api.Config{
		Address: addr,
	}

	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)

	return &Vault{Client: c}, nil
}
