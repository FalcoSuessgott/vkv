package vault

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var vaultVersion = "latest"

type VaultContainer struct {
	testcontainers.Container
	URI string
}

type VaultSuite struct {
	suite.Suite
	vc  *VaultContainer
	v   *Vault
	ctx context.Context
}

func spinUpVault(ctx context.Context) (*VaultContainer, error) {
	if v, ok := os.LookupEnv("VAULT_VERSION"); ok {
		vaultVersion = v
	}

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("hashicorp/vault:%s", vaultVersion),
		ExposedPorts: []string{"8200/tcp"},
		WaitingFor:   wait.ForListeningPort("8200/tcp"),
		Cmd:          []string{"server", "-dev", "-dev-root-token-id", "root"},
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
		log.Fatalf("error getting ip: %v", err)
	}

	mappedPort, err := c.MappedPort(ctx, "8200")
	if err != nil {
		log.Fatalf("error getting socket: %v", err)
	}

	uri := fmt.Sprintf("http://%s", net.JoinHostPort(ip, mappedPort.Port()))

	return &VaultContainer{Container: c, URI: uri}, nil
}

func (s *VaultSuite) SetupTest() {
	vc, err := spinUpVault(s.ctx)
	if err != nil {
		log.Fatalf("error spinning up vault: %v", err)
	}

	s.vc = vc

	os.Setenv("VAULT_ADDR", s.vc.URI)
	// os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	os.Setenv("VAULT_TOKEN", "root")

	vaultClient, err := NewClient()
	if err != nil {
		log.Fatalf("error creating vault client: %v", err)
	}

	s.v = vaultClient
}

func (s *VaultSuite) TestNewClient() {
	testCases := []struct {
		name    string
		envVars map[string]string
		err     bool
	}{
		{
			name: "test: valid options",
			envVars: map[string]string{
				"VAULT_ADDR":  "localhost",
				"VAULT_TOKEN": "root",
			},
			err: false,
		},
		{
			name: "test: vault addr missing",
			envVars: map[string]string{
				"VAULT_TOKEN": "root",
			},
			err: true,
		},
		{
			name: "test: vault token missing",
			envVars: map[string]string{
				"VAULT_ADDR": "root",
			},
			err: true,
		},
		{
			name:    "test: vault token and address missing",
			envVars: map[string]string{},
			err:     true,
		},
	}

	for _, tc := range testCases {
		os.Unsetenv("VAULT_ADDR")
		os.Unsetenv("VAULT_TOKEN")

		for k, v := range tc.envVars {
			os.Setenv(k, v)
		}

		_, err := NewClient()

		if tc.err {
			assert.Error(s.Suite.T(), err, tc.name)

			continue
		}

		assert.NoError(s.Suite.T(), err, tc.name)

		for k := range tc.envVars {
			os.Unsetenv(k)
		}
	}
}

func TestVaultSuite(t *testing.T) {
	vs := new(VaultSuite)
	vs.ctx = context.Background()

	// github actions doenst offer the docker sock, which we need
	// to run this test suite
	if runtime.GOOS == "linux" {
		suite.Run(t, vs)

		//nolint: errcheck
		defer vs.vc.Terminate(vs.ctx)
	}
}
