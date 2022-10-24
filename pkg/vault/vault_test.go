package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"
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

func (s *VaultSuite) TestWriteReadSecrets() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "secret",
			secrets: map[string]interface{}{
				"user": "password",
			},
		},
		{
			name:     "test: multiple secrets",
			rootPath: "kv",
			subPath:  "secret",
			secrets: map[string]interface{}{
				"user":  "password",
				"value": "42",
			},
		},
		{
			name:     "test: sub path secrets",
			rootPath: "kv",
			subPath:  "secret/sub",
			secrets: map[string]interface{}{
				"user":  "password",
				"value": "42",
			},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		// enable kv engine again, so it erros
		assert.Error(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		// read secrets- find none, so it errors
		_, err := s.v.ReadSecrets(tc.rootPath, tc.subPath)
		assert.Error(s.Suite.T(), err)

		// actual write the secrets
		err = s.v.WriteSecrets(tc.rootPath, tc.subPath, tc.secrets)
		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)

			// read them, expect the exact same secrets as written before
			readSecrets, err := s.v.ReadSecrets(tc.rootPath, tc.subPath)
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.secrets, readSecrets, tc.name)
		}
		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
	}
}

func (s *VaultSuite) TestListSecrets() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
		expected []string
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub2": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: []string{"sub", "sub2"},
		},
		{
			name:     "test: multiple dirs",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub/sub2/sub3": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: []string{"sub", "sub/"},
		},
		{
			name:     "test: empty",
			rootPath: "kv",
			subPath:  "subpath",
			secrets:  nil,
			err:      true,
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		for k, v := range tc.secrets {
			m, ok := v.(map[string]interface{})
			if ok {
				assert.NoError(s.Suite.T(), s.v.WriteSecrets(tc.rootPath, path.Join(tc.subPath, k), m))
			} else {
				fmt.Println("no")
			}
		}

		// read them, expect the exact same secrets as written before
		elements, err := s.v.ListKeys(tc.rootPath, tc.subPath)

		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.expected, elements, tc.name)
		}

		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
	}
}

func (s *VaultSuite) TestListRecursive() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub2": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: map[string]interface{}{
				"kv": Secrets{
					"sub": map[string]interface{}{
						"user":  "password",
						"user1": "password",
						"user2": "password",
					},
					"sub2": map[string]interface{}{
						"user":  "password",
						"user1": "password",
						"user2": "password",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		for k, v := range tc.secrets {
			if m, ok := v.(map[string]interface{}); ok {
				assert.NoError(s.Suite.T(), s.v.WriteSecrets(tc.rootPath, path.Join(tc.subPath, k), m))
			}
		}

		// read them, expect the exact same secrets as written before
		secrets := map[string]interface{}{}
		tmp, err := s.v.ListRecursive(tc.rootPath, tc.subPath)
		secrets[tc.rootPath] = *tmp

		assert.NoError(s.Suite.T(), err)

		_, err = json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}

		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.expected, secrets, tc.name)
		}

		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
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
