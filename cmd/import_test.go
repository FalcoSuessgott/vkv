package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportEnvVars(t *testing.T) {
	testCases := []struct {
		name     string
		envs     map[string]interface{}
		err      bool
		expected *ImportOptions
	}{
		{
			name: "force",
			err:  false,
			envs: map[string]interface{}{
				"VKV_IMPORT_FORCE":       true,
				"VKV_IMPORT_DRY_RUN":     true,
				"VKV_IMPORT_SILENT":      true,
				"VKV_IMPORT_PATH":        "path",
				"VKV_IMPORT_ENGINE_PATH": "engine",
				"VKV_IMPORT_FILE":        "file",
			},
			expected: &ImportOptions{
				Force:          true,
				DryRun:         true,
				Silent:         true,
				Path:           "path",
				MaxValueLength: 12,
				File:           "file",
			},
		},
		{
			name: "error",
			err:  true,
			envs: map[string]interface{}{
				"VKV_IMPORT_FORCE": "invalid",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		o := &ImportOptions{}

		os.Clearenv()

		for k, v := range tc.envs {
			os.Setenv(k, fmt.Sprintf("%v", v))
		}

		err := o.parseEnvs()

		for k := range tc.envs {
			os.Unsetenv(k)
		}

		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err, tc.name)
		assert.Equal(t, tc.expected, o, tc.name)
	}
}

func TestValidateImportFlags(t *testing.T) {
	testCases := []struct {
		name string
		opts *ImportOptions
		err  bool
	}{
		{
			name: "force and dryrun",
			err:  true,
			opts: &ImportOptions{
				Force:  true,
				DryRun: true,
				Path:   "tmp",
			},
		},
		{
			name: "no paths",
			err:  true,
			opts: &ImportOptions{
				Path: "",
			},
		},
		{
			name: "silent and dryrun",
			err:  true,
			opts: &ImportOptions{
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
			name:  "file",
			stdin: false,
			file:  "testdata/1.yaml",
			expected: `yaml/:
  secret:
    user: password

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
		cmd := newImportCmd()
		o := &ImportOptions{
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

		assert.Equal(t, tc.expected, string(out), tc.name)
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
		cmd := newImportCmd()

		o := &ImportOptions{
			writer: cmd.OutOrStdout(),
		}
		m, err := o.parseInput(tc.input)

		if tc.err {
			require.Error(t, err, tc.name)
		}

		assert.Equal(t, tc.expected, m, tc.name)
	}
}
