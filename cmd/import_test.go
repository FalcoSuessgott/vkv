package cmd

import (
	"bytes"
	"io"

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
