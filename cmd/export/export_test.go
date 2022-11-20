package export

import (
	"bytes"
	"io"
	"log"
	"runtime"
	"testing"

	imp "github.com/FalcoSuessgott/vkv/cmd/imp"
	"github.com/FalcoSuessgott/vkv/pkg/fs"
	printer "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/testutils"
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

func (s *VaultSuite) TestExportCommand() {
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
			importCmd := imp.NewImportCmd(io.Discard, s.client)

			importCmd.SetOut(io.Discard)
			importCmd.SetArgs(tc.importArgs)

			err := importCmd.Execute()
			if tc.importErr {
				require.Error(s.Suite.T(), err, tc.name)
			}

			require.NoError(s.Suite.T(), err, tc.name)

			// 2. read secrets and compare
			b := bytes.NewBufferString("")

			exportCmd := NewExportCmd(b, s.client)
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
		o := &exportOptions{
			Path:         "kv",
			FormatString: tc.format,
			TemplateFile: "o", // needed for testing template format
		}

		err := o.validateFlags()

		if tc.err {
			assert.ErrorIs(t, err, printer.ErrInvalidFormat, tc.name)
		} else {
			assert.NoError(t, err, tc.name, tc.name)
			assert.Equal(t, tc.expected, o.outputFormat, tc.name)
		}
	}
}

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name string
		opts *exportOptions
		err  bool
	}{
		{
			name: "test: only keys and only paths",
			err:  true,
			opts: &exportOptions{Path: "kv", OnlyKeys: true, OnlyPaths: true},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			opts: &exportOptions{Path: "kv", OnlyKeys: true, ShowValues: true},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			opts: &exportOptions{Path: "kv", OnlyPaths: true, ShowValues: true},
		},
		{
			name: "test: no paths",
			err:  false,
			opts: &exportOptions{FormatString: "base", Path: "kv"},
		},
		{
			name: "test: template with file",
			err:  false,
			opts: &exportOptions{Path: "kv", FormatString: "tmpl", TemplateFile: "OK"},
		},
		{
			name: "test: template with string",
			err:  false,
			opts: &exportOptions{Path: "kv", FormatString: "tmpl", TemplateString: "OK"},
		},
		{
			name: "test: template no file or string",
			err:  true,
			opts: &exportOptions{Path: "kv", FormatString: "tmpl"},
		},
		{
			name: "test: template file and string",
			err:  true,
			opts: &exportOptions{Path: "kv", FormatString: "tmpl", TemplateString: "OK", TemplateFile: "OK"},
		},
	}

	for _, tc := range testCases {
		err := tc.opts.validateFlags()
		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err, tc.name)
	}
}
