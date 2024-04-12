package vault

import (
	"errors"
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
	// error if no VAULT_ADDR exported
	_, ok := os.LookupEnv("VAULT_ADDR")
	if !ok {
		return nil, errors.New("VAULT_ADDR required but not set")
	}

	// get vault token
	vaultToken, tokenExported := os.LookupEnv("VAULT_TOKEN")

	// if none exported, check for VKV_LOGIN_COMMAND, execute it, and set the output as token
	cmd, ok := os.LookupEnv("VKV_LOGIN_COMMAND")
	if !tokenExported && ok {
		cmdParts := strings.Split(cmd, " ")

		token, err := exec.Run(cmdParts)
		if err != nil {
			return nil, fmt.Errorf("error running VKV_LOGIN_CMD (%s): %w", cmd, err)
		}

		vaultToken = strings.TrimSpace(string(token))
	}

	// if toke is still empty, error
	if vaultToken == "" {
		return nil, errors.New("VKV_LOGIN_COMMAND or VAULT_TOKEN required but not set")
	}

	// read all other vault env vars
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// set token
	c.SetToken(vaultToken)

	// and namespace
	if vaultNamespace, ok := os.LookupEnv("VAULT_NAMESPACE"); ok {
		c.SetNamespace(vaultNamespace)
	}

	// self lookup current auth for verification
	if _, err := c.Auth().Token().LookupSelf(); err != nil {
		return nil, fmt.Errorf("not authenticated. Perhaps not a valid token: %w", err)
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
