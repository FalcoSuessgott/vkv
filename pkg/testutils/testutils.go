package testutils

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	vaultVersion = "latest"
	image        = fmt.Sprintf("hashicorp/vault:%s", vaultVersion)
	envs         = map[string]string{}
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
		fmt.Println(license)
		envs["VAULT_LICENSE"] = license
		image = fmt.Sprintf("hashicorp/vault-enterprise:%s", vaultVersion)
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"8200/tcp"},
		WaitingFor:   wait.ForListeningPort("8200/tcp"),
		Cmd:          []string{"server", "-dev", "-dev-root-token-id", "root"},
		AutoRemove:   true,
		Env:          envs,
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := c.MappedPort(ctx, "8200")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s", net.JoinHostPort(ip, mappedPort.Port()))

	fmt.Printf("started container: %s (%s)\n", c.GetContainerID(), uri)

	return &TestContainer{Container: c, ctx: ctx, URI: uri, Token: "root"}, nil
}

// Terminate terminates the testcontainer.
func (v *TestContainer) Terminate() error {
	time.Sleep(1 * time.Second)

	if err := v.Container.Terminate(v.ctx); err != nil {
		return err
	}

	fmt.Printf("terminated container: %s (%s)\n", v.Container.GetContainerID(), v.URI)

	return nil
}
