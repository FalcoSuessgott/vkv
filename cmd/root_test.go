package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/testutils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/suite"
)

type VaultSuite struct {
	suite.Suite

	c *testutils.TestContainer
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

	v, err := vault.NewClient(vc.URI, vc.Token)
	if err != nil {
		log.Fatal(err)
	}

	vaultClient = v
}

func (s *VaultSuite) TestMode() {
	testCases := []struct {
		name     string
		envs     map[string]string
		secrets  map[string]interface{}
		expected string
		err      bool
	}{
		{
			name: "export",
			envs: map[string]string{
				"VKV_MODE":                  "export",
				"VKV_EXPORT_PATH":           "e2e",
				"VKV_EXPORT_WITH_HYPERLINK": "false",
			},
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user": "password",
				},
				"sub2": map[string]interface{}{
					"key": false,
				},
			},
			expected: `e2e/ [type=kv2]
├── sub [v=1]
│   └── user=********
│   
└── sub2 [v=1]
    └── key=*****
`,
		},
		{
			name: "export",
			envs: map[string]string{
				"VKV_MODE": "invalid",
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// inject writer
			b := bytes.NewBufferString("")
			writer = b

			// enable kv engine
			s.Require().NoError(vaultClient.EnableKV2Engine("e2e"), "enabling KV engine")

			// write secrets
			for k, secrets := range tc.secrets {
				if m, ok := secrets.(map[string]interface{}); ok {
					s.Require().NoError(vaultClient.WriteSecrets("e2e", k, m))
				}
			}

			// set env vars
			for k, v := range tc.envs {
				s.Suite.T().Setenv(k, v)
			}

			err := NewRootCmd().Execute()
			fmt.Println(err)
			// run vkv
			s.Require().Equal(tc.err, err != nil, "error "+tc.name)

			// assert output
			if !tc.err {
				out, _ := io.ReadAll(b)
				s.Require().Equal(tc.expected, string(out), tc.name)
			}
		})
	}
}

func TestVaultSuite(t *testing.T) {
	// github actions doesn't offer the docker socket, which we need to run this test suite
	if runtime.GOOS != "windows" {
		suite.Run(t, new(VaultSuite))
	}
}
