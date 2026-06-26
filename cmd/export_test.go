package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer/secret"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
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
		{
			name: "all-versions rejects markdown format",
			args: []string{"-p=1", "--all-versions", "-f=markdown"},
			err:  true,
		},
		{
			name: "all-versions rejects policy format",
			args: []string{"-p=1", "--all-versions", "-f=policy"},
			err:  true,
		},
		{
			name: "all-versions and merge-paths mutually exclusive",
			args: []string{"-p=1", "--all-versions", "--merge-paths"},
			err:  true,
		},
		{
			name: "all-versions and only-paths mutually exclusive",
			args: []string{"-p=1", "--all-versions", "--only-paths"},
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
			expected: `secret:
  user: password
`,
		},
		{
			name:          "import secrets, overwrite path, read from path",
			importCmdArgs: []string{"-f=testdata/1.yaml", "-p=yaml2"},
			exportCmdArgs: []string{"-p=yaml2", "-f=yaml", "--show-values"},
			expected: `secret:
  user: password
`,
		},
		{
			name:          "import secrets, overwrite path and subpath, read from path",
			importCmdArgs: []string{"-f=testdata/1.yaml", "-p=yaml2/sub"},
			exportCmdArgs: []string{"-p=yaml2", "-f=yaml", "--show-values"},
			expected: `sub/secret:
  user: password
`,
		},
		{
			name:          "import secrets, overwrite path with engine path, export from engine path",
			importCmdArgs: []string{"-f=testdata/2.json", "-e=engine/path"},
			exportCmdArgs: []string{"-e=engine/path", "-f=json", "--show-values"},
			expected: `{
  "admin": {
    "sub": "password"
  }
}
`,
		},
		{
			name:          "import secrets, overwrite path with engine path and subpath, export from engine path",
			importCmdArgs: []string{"-f=testdata/2.json", "-e=engine/path", "-p=sub"},
			exportCmdArgs: []string{"-e=engine/path", "-f=json", "--show-values"},
			expected: `{
  "sub/admin": {
    "sub": "password"
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

func (s *VaultSuite) TestExportAllVersions() {
	s.Run("export all versions", func() {
		ctx := context.Background()

		s.Require().NoError(vaultClient.EnableKV2Engine(ctx, "versions"))

		// write two versions of the same secret + a nested secret
		s.Require().NoError(vaultClient.WriteSecrets(ctx, "versions", "admin", map[string]interface{}{"user": "v1"}))
		s.Require().NoError(vaultClient.WriteSecrets(ctx, "versions", "admin", map[string]interface{}{"user": "v2"}))
		s.Require().NoError(vaultClient.WriteSecrets(ctx, "versions", "sub/demo", map[string]interface{}{"foo": "bar"}))

		b := bytes.NewBufferString("")
		writer = b

		exportCmd := NewExportCmd()
		exportCmd.SetArgs([]string{"-p=versions", "--all-versions", "--show-values", "--with-hyperlink=false"})

		s.Require().NoError(exportCmd.Execute())

		out, _ := io.ReadAll(b)
		output := string(out)

		// both versions of admin appear, each showing its own value
		s.Require().Contains(output, "[Version 2 created")
		s.Require().Contains(output, "[Version 1 created")
		s.Require().Contains(output, "user=v2")
		s.Require().Contains(output, "user=v1")
		// nested secret appears with its value
		s.Require().Contains(output, "sub")
		s.Require().Contains(output, "foo=bar")

		// same data is also exportable as JSON (real values, versioned schema)
		jb := bytes.NewBufferString("")
		writer = jb

		jsonCmd := NewExportCmd()
		jsonCmd.SetArgs([]string{"-p=versions", "--all-versions", "-f=json"})

		s.Require().NoError(jsonCmd.Execute())

		var parsed map[string]vault.VersionedSecret
		s.Require().NoError(json.Unmarshal(jb.Bytes(), &parsed), "all-versions JSON must be valid")
		s.Require().Len(parsed["admin"].Versions, 2)
		s.Require().Equal("v2", parsed["admin"].Versions[0].Data["user"])
		s.Require().Equal("v1", parsed["admin"].Versions[1].Data["user"])
	})
}

func (s *VaultSuite) TestExportAllVersionsKVv1() {
	s.Run("export all versions on a KVv1 engine errors", func() {
		ctx := context.Background()

		s.Require().NoError(vaultClient.EnableKV1Engine(ctx, "kvv1"))
		s.Require().NoError(vaultClient.WriteSecrets(ctx, "kvv1", "admin", map[string]interface{}{"user": "v1"}))

		writer = io.Discard

		exportCmd := NewExportCmd()
		exportCmd.SetArgs([]string{"-p=kvv1", "--all-versions"})

		s.Require().Error(exportCmd.Execute(), "--all-versions should fail on a KVv1 engine")
	})
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
