package testutils

import (
	"context"
	"fmt"
	//"net"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var vaultVersion = "latest"

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

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("hashicorp/vault-enterprise:%s", vaultVersion),
		ExposedPorts: []string{"8200/tcp"},
		WaitingFor:   wait.ForListeningPort("8200/tcp"),
		Cmd:          []string{"server", "-dev", "-dev-root-token-id", "root"},
		AutoRemove:   true,
		Env: map[string]string{
			"VAULT_LICENSE": os.Getenv("VAULT_LICENSE"),
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// ip, err := c.Host(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// mappedPort, err := c.MappedPort(ctx, "8200")
	// if err != nil {
	// 	return nil, err
	// }

	// uri := fmt.Sprintf("http://%s", net.JoinHostPort(ip, mappedPort.Port()))

	uri := ("http://127.0.0.1:8200")
	fmt.Printf("started container: %s (%s)\n", c.GetContainerID(), uri)

	time.Sleep(1 * time.Second)

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
