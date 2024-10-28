package vault

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/exec"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/tokenhelper"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client
}

// NewDefaultClient returns a new vault client wrapper.
func NewDefaultClient() (*Vault, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	// create vault client using defaults (recommended)
	c, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	c.SetToken(token)

	// self lookup current auth for verification
	if _, err := c.Auth().Token().LookupSelf(); err != nil {
		return nil, fmt.Errorf("not authenticated, perhaps not a valid token: %w", err)
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

// getToken finds the token configured by the user via env vars or token helpers
// Precedence: 1. VAULT_TOKEN, 2. VKV_LOGIN_COMMAND, 3. Vault Token Helper.
//
//nolint:cyclop
func getToken() (string, error) {
	// warn user if more than one is configured
	envToken, envTokenOk := os.LookupEnv("VAULT_TOKEN")
	tokenCommand, tokenCommandOk := os.LookupEnv("VKV_LOGIN_COMMAND")

	th, err := tokenhelper.NewInternalTokenHelper()
	if err != nil {
		return "", fmt.Errorf("error creating default token helper: %w", err)
	}

	thToken, err := th.Get()
	if err != nil {
		return "", fmt.Errorf("error getting token from default token helper: %w", err)
	}

	var (
		// number of tokens configured
		tokenSources int

		// if we issue a warning to the user, we also want to inform with what token option we went
		warn bool
	)

	if envTokenOk {
		tokenSources++
	}

	if tokenCommandOk {
		tokenSources++
	}

	if thToken != "" {
		tokenSources++
	}

	// check whether user disabled warnings
	_, disableWarn := os.LookupEnv("VKV_DISABLE_WARNING")

	if tokenSources > 1 {
		warn = true

		if !disableWarn {
			fmt.Println("[WARN] More than one token source configured (either VAULT_TOKEN, VKV_LOGIN_COMMAND or ~/.vault-token).")
			fmt.Println("[WARN] See https://falcosuessgott.github.io/vkv/authentication/ for vkv's token precedence logic. Disable these warnings with VKV_DISABLE_WARNING.")
		}
	}

	// if VAULT_TOKEN is set - return it
	if envToken != "" {
		if warn && !disableWarn {
			fmt.Println("[INFO] Using VAULT_TOKEN.")
			fmt.Println()
		}

		return envToken, nil
	}

	// if VKV_LOGIN_COMMAND
	if tokenCommand != "" {
		if warn && !disableWarn {
			fmt.Println("[INFO] Using VKV_LOGIN_COMMAND.")
			fmt.Println()
		}

		return runVaultTokenCommand(tokenCommand)
	}

	if thToken != "" {
		if warn && !disableWarn {
			fmt.Println("[INFO] Using ~/.vault-token.")
			fmt.Println()
		}

		return thToken, nil
	}

	return "", errors.New("no token provided")
}

func runVaultTokenCommand(cmd string) (string, error) {
	cmdParts := strings.Split(cmd, " ")

	token, err := exec.Run(cmdParts)
	if err != nil {
		return "", fmt.Errorf("error running VKV_LOGIN_CMD (%s): %w", cmd, err)
	}

	vaultToken := strings.TrimSpace(string(token))
	if vaultToken == "" {
		return "", errors.New("VKV_LOGIN_COMMAND required but not set")
	}

	return vaultToken, nil
}
