package cmd

import (
	"bytes"
	"io"
)

func (s *VaultSuite) TestValidateExportFlags() {
	testCases := []struct {
		name string
		args []string
		opts *importOptions
		err  bool
	}{
		{
			name: "only keys and only paths mutually exclusive",
			args: []string{"-p=1", "--only-keys", "--only-paths"},
			err:  true,
		},
		{
			name: "only keys and show values mutually exclusive",
			args: []string{"-p=1", "--only-keys", "--show-values"},
			err:  true,
		},
		{
			name: "only paths and show values mutually exclusive",
			args: []string{"-p=1", "--only-paths", "--show-values"},
			err:  true,
		},
	}

	for _, tc := range testCases {
		cmd := NewExportCmd()
		cmd.SetArgs(tc.args)

		err := cmd.Execute()

		s.Require().Equal(tc.err, err != nil, tc.name)
	}
}

func (s *VaultSuite) TestExportImportCommand() {
	testCases := []struct {
		name          string
		expected      string
		err           bool
		importCmdArgs []string
		exportCmdArgs []string
	}{
		{
			name:          "import secrets, export from path",
			importCmdArgs: []string{"-f=testdata/1.yaml"},
			exportCmdArgs: []string{"-p=yaml", "-f=yaml", "--show-values"},
			expected: `yaml/:
  secret:
    user: password
`,
		},
		{
			name:          "import secrets, overwrite path, read from path",
			importCmdArgs: []string{"-f=testdata/1.yaml", "-p=yaml2"},
			exportCmdArgs: []string{"-p=yaml2", "-f=yaml", "--show-values"},
			expected: `yaml2/:
  secret:
    user: password
`,
		},
		{
			name:          "import secrets, overwrite path and subpath, read from path",
			importCmdArgs: []string{"-f=testdata/1.yaml", "-p=yaml2/sub"},
			exportCmdArgs: []string{"-p=yaml2", "-f=yaml", "--show-values"},
			expected: `yaml2/:
  sub/:
    secret:
      user: password
`,
		},
		{
			name:          "import secrets, overwrite path with engine path, export from engine path",
			importCmdArgs: []string{"-f=testdata/2.json", "-e=engine/path"},
			exportCmdArgs: []string{"-e=engine/path", "-f=json", "--show-values"},
			expected: `{
  "engine/path/": {
    "admin": {
      "sub": "password"
    }
  }
}
`,
		},
		{
			name:          "import secrets, overwrite path with engine path and subpath, export from engine path",
			importCmdArgs: []string{"-f=testdata/2.json", "-e=engine/path", "-p=sub"},
			exportCmdArgs: []string{"-e=engine/path", "-f=json", "--show-values"},
			expected: `{
  "engine/path/": {
    "sub/": {
      "admin": {
        "sub": "password"
      }
    }
  }
}
`,
		},
		{
			name:          "no output, error, dryrun",
			err:           true,
			importCmdArgs: []string{"-f=testdata/2.json", "-d", "-p=json"},
			exportCmdArgs: []string{"-p=json", "-f=json", "--show-values"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// 0. inject writer & vault client
			writer = io.Discard

			// 1. import secrets
			importCmd := NewImportCmd()
			importCmd.SetArgs(tc.importCmdArgs)

			s.Require().NoError(importCmd.Execute(), "import "+tc.name)

			// 2. read secrets, capture output and assert
			b := bytes.NewBufferString("")
			writer = b

			exportCmd := NewExportCmd()
			exportCmd.SetArgs(tc.exportCmdArgs)

			err := exportCmd.Execute()

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

// func (s *VaultSuite) TestExportOutputFormat() {
// 	testCases := []struct {
// 		name     string
// 		expected prt.OutputFormat
// 		err      bool
// 	}{
// 		{
// 			name:     "json",
// 			expected: prt.JSON,
// 		},
// 		{
// 			name:     "yaml",
// 			expected: prt.YAML,
// 		},
// 		{
// 			name:     "yml",
// 			expected: prt.YAML,
// 		},
// 		{
// 			name:     "invalid",
// 			err:      true,
// 			expected: prt.YAML,
// 		},
// 		{
// 			name:     "export",
// 			expected: prt.Export,
// 		},
// 		{
// 			name:     "markdown",
// 			expected: prt.Markdown,
// 		},
// 		{
// 			name:     "base",
// 			expected: prt.Base,
// 		},
// 		{
// 			name:     "template",
// 			expected: prt.Template,
// 		},
// 		{
// 			name:     "tmpl",
// 			expected: prt.Template,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		o := &exportOptions{
// 			FormatString:   tc.name,
// 			Path:           "kv",
// 			TemplateString: "tmpl",
// 		}

// 		err := o.validateFlags(nil, nil)

// 		s.Require().Equal(tc.err, err != nil, "error "+tc.name)

// 		// if no error -> assert output format
// 		if !tc.err {
// 			s.Require().Equal(tc.expected, o.outputFormat, "format "+tc.name)
// 		}
// 	}
// }
