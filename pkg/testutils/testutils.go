package testutils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	vaultVersion = "latest"
	image        = fmt.Sprintf("hashicorp/vault:%s", vaultVersion)
	envs         = map[string]string{}
	token        = "root"
)

// TestContainer vault dev container wrapper.
type TestContainer struct {
	Container testcontainers.Container
	ctx       context.Context
	URI       string
	Token     string
}

// StartTestContainer Starts a fresh vault in development mode.
func StartTestContainer() (*TestContainer, error) {
	ctx := context.Background()

	if v, ok := os.LookupEnv("VAULT_VERSION"); ok {
		vaultVersion = v
	}

	// use OSS image per default, if license is available use enterprise
	if license, ok := os.LookupEnv("VAULT_LICENSE"); ok {
		envs["VAULT_LICENSE"] = license
		image = fmt.Sprintf("hashicorp/vault-enterprise:%s", vaultVersion)
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"8200/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Root Token:").WithPollInterval(1*time.Second).WithStartupTimeout(3*time.Minute),
			wait.ForListeningPort("8200/tcp"),
		),
		Cmd:        []string{"server", "-dev", "-dev-root-token-id", token},
		AutoRemove: true,
		Env:        envs,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.CapAdd = []string{"IPC_LOCK"}
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := c.MappedPort(ctx, "8200")
	if err != nil {
		return nil, err
	}

	return &TestContainer{
		Container: c, ctx: ctx,
		URI:   fmt.Sprintf("http://127.0.0.1:%s", mappedPort.Port()),
		Token: token,
	}, nil
}

// Terminate terminates the testcontainer.
func (v *TestContainer) Terminate() error {
	if err := v.Container.Terminate(v.ctx); err != nil {
		return err
	}

	return nil
}
