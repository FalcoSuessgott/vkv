package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	v   *vault.Vault
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

	vaultClient, err := vault.NewClient()
	if err != nil {
		log.Fatalf("error creating vault client: %v", err)
	}

	s.v = vaultClient
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

func (s *VaultSuite) TestE2E() {
	testCases := []struct {
		name       string
		enginePath string
		rootPath   string
		expected   string
		err        bool
		importArgs []string
		exportArgs []string
	}{
		{
			name:       "yaml",
			rootPath:   "yaml",
			importArgs: []string{"-f=testdata/1.yaml", "-p=yaml"},
			exportArgs: []string{"-p=yaml", "-f=yaml", "--show-values"},
			expected:   "testdata/1.yaml",
		},
		{
			name:       "json",
			rootPath:   "json",
			importArgs: []string{"-f=testdata/2.json", "-p=json"},
			exportArgs: []string{"-p=json", "-f=json", "--show-values"},
			expected:   "testdata/2.json",
		},
	}

	for _, tc := range testCases {
		// 1. import secrets
		importCmd := newImportCmd()

		importCmd.SetOut(io.Discard)
		importCmd.SetArgs(tc.importArgs)

		require.NoError(s.Suite.T(), importCmd.Execute())

		// 2. read secrets and compare
		b := bytes.NewBufferString("")

		rootCmd := newRootCmd("", b)
		rootCmd.SetOut(b)
		rootCmd.SetArgs(tc.exportArgs)

		err := rootCmd.Execute()
		if tc.err {
			require.NoError(s.Suite.T(), err, tc.name)

			continue
		}

		// assert
		out, _ := io.ReadAll(b)

		exp, _ := utils.ReadFile(tc.expected)
		assert.Equal(s.Suite.T(), string(exp), string(out), tc.name)
	}
}
