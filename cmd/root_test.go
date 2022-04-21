package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/stretchr/testify/assert"
)

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
	}

	for _, tc := range testCases {
		o := &Options{
			FormatString: tc.format,
		}

		err := o.validateFlags()

		if tc.err {
			assert.ErrorIs(t, err, printer.ErrInvalidFormat, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
			assert.Equal(t, tc.expected, o.outputFormat)
		}
	}
}

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			name: "test: only keys and only paths",
			err:  true,
			args: []string{"--only-keys", "--only-paths"},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			args: []string{"--only-keys", "--show-values"},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			args: []string{"--only-paths", "--show-values"},
		},
		{
			name: "test: no paths",
			err:  false,
			args: []string{"--path", ""},
		},
	}

	for _, tc := range testCases {
		c := newRootCmd("")
		b := bytes.NewBufferString("")

		c.SetArgs(tc.args)
		c.SetOut(b)

		os.Setenv("VAULT_ADDR", "")
		os.Setenv("VAULT_TOKEN", "")

		err := c.Execute()
		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err, tc.name)
	}
}

func TestEnvVars(t *testing.T) {
	testCases := []struct {
		name     string
		envs     map[string]interface{}
		err      bool
		expected *Options
	}{
		{
			name: "only keys",
			err:  false,
			envs: map[string]interface{}{
				"VKV_ONLY_KEYS": true,
			},
			expected: &Options{
				MaxValueLength: printer.MaxValueLength,
				Paths:          []string{"kv"},
				FormatString:   "base",
				OnlyKeys:       true,
			},
		},
		{
			name: "only paths",
			err:  false,
			envs: map[string]interface{}{
				"VKV_ONLY_PATHS": true,
			},
			expected: &Options{
				MaxValueLength: printer.MaxValueLength,
				Paths:          []string{"kv"},
				FormatString:   "base",
				OnlyPaths:      true,
			},
		},
		{
			name: "invalid value only paths",
			err:  true,
			envs: map[string]interface{}{
				"VKV_ONLY_PATHS": "invalid",
			},
		},
		{
			name: "show values and max value length",
			err:  false,
			envs: map[string]interface{}{
				"VKV_SHOW_VALUES":      true,
				"VKV_MAX_VALUE_LENGTH": 213,
			},
			expected: &Options{
				MaxValueLength: 213,
				Paths:          []string{"kv"},
				FormatString:   "base",
				ShowValues:     true,
			},
		},
		{
			name: "format",
			err:  false,
			envs: map[string]interface{}{
				"VKV_FORMAT": "yaml",
			},
			expected: &Options{
				MaxValueLength: 12,
				Paths:          []string{"kv"},
				FormatString:   "yaml",
			},
		},
		{
			name: "show values and max value length",
			err:  false,
			envs: map[string]interface{}{
				"VKV_PATHS":            "kv1,kv2,kv3",
				"VKV_MAX_VALUE_LENGTH": 213,
			},
			expected: &Options{
				MaxValueLength: 213,
				Paths:          []string{"kv1", "kv2", "kv3"},
				FormatString:   "base",
			},
		},
	}

	for _, tc := range testCases {
		o := &Options{}

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
		assert.Equal(t, tc.expected, o)
	}
}
