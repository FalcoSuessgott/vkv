package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

func (s *VaultSuite) TestValidateImportFlags() {
	testCases := []struct {
		name string
		args []string
		opts *importOptions
		err  bool
	}{
		{
			name: "force and dry-run mutually exclusive",
			args: []string{"--force", "--dry-run"},
			err:  true,
		},
		{
			name: "silent and dry-run should fail",
			err:  true,
			args: []string{"-p=o", "-s", "-d"},
		},
		{
			name: "file and STDIN should fail",
			err:  true,
			args: []string{"-f=sss", "-"},
		},
	}

	for _, tc := range testCases {
		cmd := NewImportCmd()

		writer = io.Discard

		cmd.SetArgs(tc.args)

		err := cmd.Execute()
		s.Require().Equal(tc.err, err != nil, tc.name)
	}
}

func (s *VaultSuite) TestGetInput() {
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
			file:  "testdata/1.yaml",
			expected: `yaml/:
  secret:
    user: password
`,
		},
		{
			name:  "json file",
			stdin: false,
			file:  "testdata/2.json",
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
		writer = io.Discard

		cmd := NewImportCmd()
		o := &importOptions{
			input: bytes.NewReader([]byte(tc.stdinInput)),
		}

		args := []string{}

		if tc.stdin {
			args = append(args, "-")
		} else {
			o.File = tc.file
		}

		cmd.SetArgs(args)

		out, err := o.getInput()

		s.Require().Equal(tc.err, err != nil, tc.name)

		if !tc.err {
			s.Require().Equal(tc.expected, utils.RemoveCarriageReturns(string(out)), tc.name)
		}
	}
}

func (s *VaultSuite) TestParseInput() {
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
		writer = io.Discard

		o := &importOptions{}

		m, err := o.parseInput(tc.input)

		s.Require().Equal(tc.err, err != nil, "error "+tc.name)

		if !tc.err {
			s.Require().Equal(tc.expected, m, "expected "+tc.name)
		}
	}
}

func (s *VaultSuite) TestImportCommand() {
	testCases := []struct {
		name            string
		importPath      string
		exportPath      string
		err             bool
		originalSecrets string
		newSecrets      string
		expected        string
	}{
		{
			name:       "simple merge of 2 secrets",
			importPath: "yaml",
			exportPath: "yaml",
			originalSecrets: `yaml/:
  secret:
    user: password
`,
			newSecrets: `yaml/:
  secret2:
    new: key
`,
			expected: `yaml/:
  secret:
    user: password
  secret2:
    new: key
`,
		},
		{
			name:       "merge of different nested secrets",
			importPath: "yaml",
			exportPath: "yaml",
			originalSecrets: `yaml/:
  secret:
    user: password
`,
			newSecrets: `yaml/:
  secret2/:
    sub/:
      foo:
        new: key
`,
			expected: `yaml/:
  secret:
    user: password
  secret2/:
    sub/:
      foo:
        new: key
`,
		},
		{
			name:       "write secrets to a subpath",
			importPath: "yaml/secret2",
			exportPath: "yaml",
			originalSecrets: `yaml/:
  secret:
    user: password
`,
			newSecrets: `yaml/:
  secret:
    user: password
`,
			expected: `yaml/:
  secret:
    user: password
  secret2/:
    secret:
      user: password
`,
		},
		{
			name:       "no original secrets",
			importPath: "yaml",
			exportPath: "yaml",
			newSecrets: `yaml/:
  secret:
    user: password
`,
			expected: `yaml/:
  secret:
    user: password
`,
		},
		{
			name:       "overwriting secrets",
			importPath: "yaml",
			exportPath: "yaml",
			originalSecrets: `yaml/:
  secret:
    user: password
`,
			newSecrets: `yaml/:
  secret:
    new: key
`,
			expected: `yaml/:
  secret:
    new: key
`,
		},
	}

	for _, tc := range testCases {
		//nolint: perfsprint, gosec
		s.Run(tc.name, func() {
			// 0. inject writer & vault client
			writer = io.Discard

			// 1. import original secrets
			if tc.originalSecrets != "" {
				o, err := os.CreateTemp(s.Suite.T().TempDir(), "original-secrets")
				s.Require().NoError(err, "temp file")

				s.Require().NoError(os.WriteFile(o.Name(), []byte(tc.originalSecrets), 0o644), "write original secrets")

				importCmd := NewImportCmd()
				importCmd.SetArgs([]string{fmt.Sprintf("-f=%s", o.Name())})
				s.Require().NoError(importCmd.Execute(), "import original "+tc.name)
			}

			// 2. import new secrets
			n, err := os.CreateTemp(s.Suite.T().TempDir(), "original-secrets")
			s.Require().NoError(err, "temp file")

			s.Require().NoError(os.WriteFile(n.Name(), []byte(tc.newSecrets), 0o644), "write new secrets")

			importCmd := NewImportCmd()
			importCmd.SetArgs([]string{fmt.Sprintf("-p=%s", tc.importPath), "--force", fmt.Sprintf("-f=%s", n.Name())})
			s.Require().NoError(importCmd.Execute(), "import new "+tc.name)

			// 3. export & assert
			b := bytes.NewBufferString("")
			writer = b

			exportCmd := NewExportCmd()
			exportCmd.SetArgs([]string{fmt.Sprintf("-p=%s", tc.exportPath), "-f=yaml", "--show-values"})

			err = exportCmd.Execute()

			if tc.err {
				s.Require().Error(err, "export "+tc.name)
			} else {
				// if no error - compare exported secrets with expected value
				out, _ := io.ReadAll(b)

				s.Require().Equal(tc.expected, string(out), "secrets "+tc.name)
			}
		})
	}
}
