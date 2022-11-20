package snapshot

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/testutils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/suite"
)

type VaultSuite struct {
	suite.Suite
	c      *testutils.TestContainer
	client *vault.Vault
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

	s.client = v
}

func TestVaultSuite(t *testing.T) {
	// github actions doenst offer the docker sock, which we need
	// to run this test suite
	if runtime.GOOS == "linux" {
		suite.Run(t, new(VaultSuite))
	}
}

func TestOutputFormat(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		expected printer.OutputFormat
		err      bool
	}{
		{
			name:     "json",
			err:      false,
			format:   "json",
			expected: printer.JSON,
		},
		{
			name:     "yaml",
			err:      false,
			format:   "YamL",
			expected: printer.YAML,
		},
		{
			name:     "yml",
			err:      false,
			format:   "yml",
			expected: printer.YAML,
		},
		{
			name:     "invalid",
			err:      true,
			format:   "invalid",
			expected: printer.YAML,
		},
		{
			name:     "export",
			err:      false,
			format:   "Export",
			expected: printer.Export,
		},
		{
			name:     "markdown",
			err:      false,
			format:   "Markdown",
			expected: printer.Markdown,
		},
		{
			name:     "base",
			err:      false,
			format:   "base",
			expected: printer.Base,
		},
		{
			name:     "template",
			err:      false,
			format:   "template",
			expected: printer.Template,
		},
		{
			name:     "tmpl",
			err:      false,
			format:   "tmpl",
			expected: printer.Template,
		},
	}

	for _, tc := range testCases {
		fmt.Println(tc)
	}
}
