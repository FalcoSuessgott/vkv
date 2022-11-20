package vault

import (
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VaultSuite struct {
	suite.Suite
	c      *testutils.TestContainer
	client *Vault
}

func (s *VaultSuite) TearDownSubTest() {
	if err := s.c.Terminate(); err != nil {
		log.Fatal(err)
	}
}

func (s *VaultSuite) SetupSubTest() {
	vc, err := testutils.StartTestContainer()
	if err != nil {
		log.Fatal(err)
	}

	s.c = vc

	v, err := NewClient(vc.URI, vc.Token)
	if err != nil {
		log.Fatal(err)
	}

	s.client = v
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

		_, err := NewDefaultClient()

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
	// github actions doenst offer the docker sock, which we need
	// to run this test suite
	if runtime.GOOS == "linux" {
		suite.Run(t, new(VaultSuite))
	}
}
