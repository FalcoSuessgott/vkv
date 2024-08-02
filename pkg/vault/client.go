package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/exec"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/tokenhelper"
)

// NewDefaultClient returns a new vault client wrapper.
func NewDefaultClient(ctx context.Context) (*Vault, error) {
	// create vault client using defaults (recommended)
	c, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	// use tokenhelper if available
	th, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return nil, fmt.Errorf("error creating default token helper: %w", err)
	}

	token, err := th.Get()
	if err != nil {
		return nil, fmt.Errorf("error getting token from default token helper: %w", err)
	}

	if token != "" {
		c.SetToken(token)
	}

	// custom: if VKV_LOGIN_COMMAND is set, execute it and set the output as token
	cmd, ok := os.LookupEnv("VKV_LOGIN_COMMAND")
	if ok && cmd != "" {
		cmdParts := strings.Split(cmd, " ")

		token, err := exec.Run(cmdParts)
		if err != nil {
			return nil, fmt.Errorf("error running VKV_LOGIN_CMD (%s): %w", cmd, err)
		}

		vaultToken := strings.TrimSpace(string(token))
		if vaultToken == "" {
			return nil, errors.New("VKV_LOGIN_COMMAND required but not set")
		}

		// set token
		c.SetToken(vaultToken)
	}

	// self lookup current auth for verification
	if _, err := c.Auth().Token().LookupSelf(); err != nil {
		return nil, fmt.Errorf("not authenticated, perhaps not a valid token: %w", err)
	}

	return &Vault{
		Client:  c,
		Context: ctx,
	}, nil
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
