package testutils

import (
	"context"
	"fmt"
	"os"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/vault"
)

var token = "root"

var (
	VaultVersionEnv = "VAULT_VERSION"
	VaultVersion    = "1.20.0"
)

// TestContainer vault dev container wrapper.
type TestContainer struct {
	Container testcontainers.Container
	URI       string
	Token     string
}

// StartTestContainer Starts a fresh vault in development mode.
func StartTestContainer(commands ...string) (*TestContainer, error) {
	ctx := context.Background()

	if v, ok := os.LookupEnv(VaultVersionEnv); ok {
		VaultVersion = v
	}

	vaultContainer, err := vault.Run(ctx, "hashicorp/vault:"+VaultVersion,
		vault.WithToken(token),
		vault.WithInitCommand(commands...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	uri, err := vaultContainer.HttpHostAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("error returning container mapped port: %w", err)
	}

	return &TestContainer{
		Container: vaultContainer,
		URI:       uri,
		Token:     token,
	}, nil
}

// Terminate terminates the testcontainer.
func (v *TestContainer) Terminate() error {
	return v.Container.Terminate(context.Background())
}
