package cmd

import (
	"bytes"
	"io"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
)

func (s *VaultSuite) TestValidateExportFlags() {
	testCases := []struct {
		name string
		args []string
		opts *importOptions
		err  bool
	}{
		{
			name: "path and engine path mutually exclusive",
			args: []string{"--path=p", "--engine-path=e"},
			err:  true,
		},
		{
			name: "path or engine path required",
			args: []string{"--show-values"},
			err:  true,
		},
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
		path          string
		expected      string
		err           bool
		importCmdArgs []string
		exportCmdArgs []string
	}{
		{
			name:          "yaml",
			path:          "yaml",
			importCmdArgs: []string{"-f=testdata/1.yaml", "-p=yaml"},
			exportCmdArgs: []string{"-p=yaml", "-f=yaml", "--show-values"},
			expected:      "testdata/1.yaml",
		},
		{
			name:          "json",
			path:          "json",
			importCmdArgs: []string{"-f=testdata/2.json", "-p=json"},
			exportCmdArgs: []string{"-p=json", "-f=json", "--show-values"},
			expected:      "testdata/2.json",
		},
		{
			name:          "dryrun",
			path:          "json",
			err:           true,
			importCmdArgs: []string{"-f=testdata/2.json", "-d", "-p=json"},
			exportCmdArgs: []string{"-p=json", "-f=json", "--show-values"},
			expected:      "testdata/2.json",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// 0. inject writer & vault client
			writer = io.Discard

			// 1. import secrets
			importCmd := NewImportCmd()
			importCmd.SetArgs(tc.importCmdArgs)

			s.Require().NoError(importCmd.Execute(), ("import " + tc.name))

			// 2. read secrets and assert
			b := bytes.NewBufferString("")
			writer = b

			exportCmd := NewExportCmd()
			exportCmd.SetArgs(tc.exportCmdArgs)

			err := exportCmd.Execute()

			s.Require().Equal(tc.err, err != nil, "export "+tc.name)

			// if no error - compare exported secretd with expected value
			if !tc.err {
				out, _ := io.ReadAll(b)
				exp, _ := fs.ReadFile(tc.expected)

				s.Require().Equal(string(exp), string(out), "secrets "+tc.name)
			}
		})
	}
}

func (s *VaultSuite) TestExportOutputFormat() {
	testCases := []struct {
		name     string
		expected prt.OutputFormat
		err      bool
	}{
		{
			name:     "json",
			expected: prt.JSON,
		},
		{
			name:     "yaml",
			expected: prt.YAML,
		},
		{
			name:     "yml",
			expected: prt.YAML,
		},
		{
			name:     "invalid",
			err:      true,
			expected: prt.YAML,
		},
		{
			name:     "export",
			expected: prt.Export,
		},
		{
			name:     "markdown",
			expected: prt.Markdown,
		},
		{
			name:     "base",
			expected: prt.Base,
		},
		{
			name:     "template",
			expected: prt.Template,
		},
		{
			name:     "tmpl",
			expected: prt.Template,
		},
	}

	for _, tc := range testCases {
		o := &exportOptions{
			FormatString:   tc.name,
			Path:           "kv",
			TemplateString: "tmpl",
		}

		err := o.validateFlags(nil, nil)

		s.Require().Equal(tc.err, err != nil, "error "+tc.name)

		// if no error -> assert output format
		if !tc.err {
			s.Require().Equal(tc.expected, o.outputFormat, "format "+tc.name)
		}
	}
}
