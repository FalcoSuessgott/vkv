package vault

import (
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/testutils"
	"github.com/stretchr/testify/require"
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
		name                           string
		envVars                        map[string]string
		addTestContainerURIAsVaultAddr bool
		err                            bool
	}{
		{
			name: "valid options",
			envVars: map[string]string{
				"VAULT_TOKEN": "root",
			},
			addTestContainerURIAsVaultAddr: true,
			err:                            false,
		},
		{
			name: "vault address missing",
			envVars: map[string]string{
				"VAULT_TOKEN": "root",
			},
			err: true,
		},
		{
			name: "vault token missing",
			envVars: map[string]string{
				"VAULT_ADDR": "root",
			},
			err: true,
		},
		{
			name:    "vault token and address missing",
			envVars: map[string]string{},
			err:     true,
		},
		{
			name: "vkv login command",
			envVars: map[string]string{
				"VKV_LOGIN_COMMAND": "echo root",
			},
			addTestContainerURIAsVaultAddr: true,
			err:                            false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// unset any VAULT env vars
			os.Unsetenv("VAULT_ADDR")
			os.Unsetenv("VAULT_TOKEN")

			// dirty hack since the uri is not available in the testcase declaration
			if tc.addTestContainerURIAsVaultAddr {
				tc.envVars["VAULT_ADDR"] = s.c.URI
			}

			// set the test case env vars
			for k, v := range tc.envVars {
				require.NoError(s.Suite.T(), os.Setenv(k, v), "error settings env var")
			}

			// auth
			_, err := NewDefaultClient()

			// assertions
			if tc.err {
				require.Error(s.Suite.T(), err, tc.name)
			} else {
				require.NoError(s.Suite.T(), err, tc.name)
			}

			// unsert test case env vars
			for k := range tc.envVars {
				os.Unsetenv(k)
			}
		})
	}
}

func TestVaultSuite(t *testing.T) {
	// github actions doenst offer the docker sock, which we need
	// to run this test suite
	if runtime.GOOS == "linux" {
		suite.Run(t, new(VaultSuite))
	}
}