package imp

import (
	"bytes"
	"io"
	"log"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/cmd/export"
	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/testutils"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (s *VaultSuite) TestImportCommand() {
	testCases := []struct {
		name       string
		enginePath string
		rootPath   string
		expected   string
		importErr  bool
		exportErr  bool
		importArgs []string
		exportArgs []string
	}{
		{
			name:       "yaml",
			rootPath:   "yaml",
			importArgs: []string{"-f=../testdata/1.yaml", "-p=yaml"},
			exportArgs: []string{"-p=yaml", "-f=yaml", "--show-values"},
			expected:   "../testdata/1.yaml",
		},
		{
			name:       "json",
			rootPath:   "json",
			importArgs: []string{"-f=../testdata/2.json", "-p=json"},
			exportArgs: []string{"-p=json", "-f=json", "--show-values"},
			expected:   "../testdata/2.json",
		},
		{
			name:       "dryrun",
			rootPath:   "json",
			exportErr:  true,
			importArgs: []string{"-f=../testdata/2.json", "-d", "-p=json"},
			exportArgs: []string{"-p=json", "-f=json", "--show-values"},
			expected:   "../testdata/2.json",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// 1. import secrets
			importCmd := NewImportCmd(io.Discard, s.client)

			importCmd.SetOut(io.Discard)
			importCmd.SetArgs(tc.importArgs)

			err := importCmd.Execute()
			if tc.importErr {
				require.Error(s.Suite.T(), err, tc.name)
			}

			require.NoError(s.Suite.T(), err, tc.name)

			// 2. read secrets and compare
			b := bytes.NewBufferString("")

			exportCmd := export.NewExportCmd(b, s.client)
			exportCmd.SetOut(b)
			exportCmd.SetArgs(tc.exportArgs)

			err = exportCmd.Execute()
			if tc.exportErr {
				require.Error(s.Suite.T(), err, tc.name)
			} else {
				// assert
				out, _ := io.ReadAll(b)

				exp, _ := fs.ReadFile(tc.expected)
				assert.Equal(s.Suite.T(), string(exp), string(out), tc.name)
			}
		})
	}
}

func TestValidateImportFlags(t *testing.T) {
	testCases := []struct {
		name string
		opts *importOptions
		err  bool
	}{
		{
			name: "force and dryrun",
			err:  true,
			opts: &importOptions{
				Force:  true,
				DryRun: true,
				Path:   "tmp",
			},
		},
		{
			name: "no paths",
			err:  true,
			opts: &importOptions{
				Path: "",
			},
		},
		{
			name: "silent and dryrun",
			err:  true,
			opts: &importOptions{
				Silent: true,
				DryRun: true,
				Path:   "tmp",
			},
		},
	}

	for _, tc := range testCases {
		err := tc.opts.validateFlags([]string{})
		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err, tc.name)
	}
}

func TestGetInput(t *testing.T) {
	testCases := []struct {
		name       string
		stdin      bool
		stdinInput string
		file       string
		expected   string
		err        bool
	}{
		{
			name:       "stdin",
			stdin:      true,
			stdinInput: "text",
			expected:   "text",
		},
		{
			name:  "yaml file",
			stdin: false,
			file:  "../testdata/1.yaml",
			expected: `yaml/:
  secret:
    user: password

`,
		},
		{
			name:  "json file",
			stdin: false,
			file:  "../testdata/2.json",
			expected: `{
  "json/": {
    "admin": {
      "sub": "password"
    }
  }
}
`,
		},
		{
			name:       "empty input",
			stdin:      true,
			stdinInput: "",
			err:        true,
		},
	}

	for _, tc := range testCases {
		cmd := NewImportCmd(nil, nil)
		o := &importOptions{
			writer: cmd.OutOrStdout(),
		}

		args := []string{}

		if tc.stdin {
			args = append(args, "-")
		} else {
			o.File = tc.file
		}

		cmd.SetArgs(args)
		cmd.SetIn(bytes.NewReader([]byte(tc.stdinInput)))

		out, err := o.getInput(cmd)

		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err)

		assert.Equal(t, tc.expected, utils.RemoveCarriageReturns(string(out)), tc.name)
	}
}

func TestParseInput(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[string]interface{}
		err      bool
	}{
		{
			name:  "invalid input",
			err:   true,
			input: []byte("invalid input"),
		},
		{
			name: "json input",
			err:  false,
			input: []byte(`{
"secret/": {
  "admin": {
    "sub": "value",
    }
  }
}`),
			expected: map[string]interface{}{
				"secret/": map[string]interface{}{
					"admin": map[string]interface{}{
						"sub": "value",
					},
				},
			},
		},
		{
			name: "yaml input",
			err:  false,
			input: []byte(`secret/:
  admin:
    sub: 'value'
`),
			expected: map[string]interface{}{
				"secret/": map[string]interface{}{
					"admin": map[string]interface{}{
						"sub": "value",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		cmd := NewImportCmd(io.Discard, nil)

		o := &importOptions{
			writer: cmd.OutOrStdout(),
		}
		m, err := o.parseInput(tc.input)

		if tc.err {
			require.Error(t, err, tc.name)
		}

		assert.Equal(t, tc.expected, m, tc.name)
	}
}
