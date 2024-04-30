package testutils

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/vault"
)

var token = "root"

// TestContainer vault dev container wrapper.
type TestContainer struct {
	Container testcontainers.Container
	URI       string
	Token     string
}

// StartTestContainer Starts a fresh vault in development mode.
func StartTestContainer(commands ...string) (*TestContainer, error) {
	ctx := context.Background()

	vaultContainer, err := vault.RunContainer(ctx,
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
