package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "v1.0.0",
			version:  "v1.0.0",
			expected: "vkv v1.0.0\n",
		},
		{
			name:     "empty",
			version:  "",
			expected: "vkv \n",
		},
	}

	for _, tc := range testCases {
		c := newRootCmd(tc.version)
		b := bytes.NewBufferString("")

		c.SetArgs([]string{"-v"})
		c.SetOut(b)

		err := c.Execute()
		assert.NoError(t, err)

		out, _ := io.ReadAll(b)
		assert.Equal(t, tc.expected, string(out), tc.name)
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
		o := &Options{
			FormatString: tc.format,
			TemplateFile: "o", // needed for testing template format
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
		opts *Options
		err  bool
	}{
		{
			name: "test: only keys and only paths",
			err:  true,
			opts: &Options{OnlyKeys: true, OnlyPaths: true},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			opts: &Options{OnlyKeys: true, ShowValues: true},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			opts: &Options{OnlyPaths: true, ShowValues: true},
		},
		{
			name: "test: no paths",
			err:  false,
			opts: &Options{FormatString: "base", Path: ""},
		},
		{
			name: "test: template with file",
			err:  false,
			opts: &Options{FormatString: "tmpl", TemplateFile: "OK"},
		},
		{
			name: "test: template with string",
			err:  false,
			opts: &Options{FormatString: "tmpl", TemplateString: "OK"},
		},
		{
			name: "test: template no file or string",
			err:  true,
			opts: &Options{FormatString: "tmpl"},
		},
		{
			name: "test: template file and string",
			err:  true,
			opts: &Options{FormatString: "tmpl", TemplateString: "OK", TemplateFile: "OK"},
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
				Path:           "kv",
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
				Path:           "kv",
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
				Path:           "kv",
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
				Path:           "kv",
				FormatString:   "yaml",
			},
		},
		{
			name: "show path and max value length",
			err:  false,
			envs: map[string]interface{}{
				"VKV_PATH":             "kv1",
				"VKV_MAX_VALUE_LENGTH": 213,
			},
			expected: &Options{
				MaxValueLength: 213,
				Path:           "kv1",
				FormatString:   "base",
			},
		},
		{
			name: "show values and max value length",
			err:  false,
			envs: map[string]interface{}{
				"VKV_PATH":            "kv1",
				"VKV_FORMAT":          "template",
				"VKV_TEMPLATE_STRING": "string",
				"VKV_TEMPLATE_FILE":   "path",
			},
			expected: &Options{
				MaxValueLength: 12,
				Path:           "kv1",
				FormatString:   "template",
				TemplateFile:   "path",
				TemplateString: "string",
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
		assert.Equal(t, tc.expected, o, tc.name)
	}
}

func TestBuildEnginePath(t *testing.T) {
	testCases := []struct {
		name             string
		expectedRootPath string
		expectedSubPath  string
		opts             *Options
	}{
		{
			name: "only path",
			opts: &Options{
				Path: "1/2/3/4",
			},
			expectedRootPath: "1",
			expectedSubPath:  "2/3/4",
		},
		{
			name: "one element path",
			opts: &Options{
				Path: "1",
			},
			expectedRootPath: "1",
			expectedSubPath:  "",
		},
		{
			name: "engine path and path",
			opts: &Options{
				EnginePath: "1/2/3/4",
				Path:       "5/6",
			},
			expectedRootPath: "1/2/3/4",
			expectedSubPath:  "5/6",
		},
		{
			name: " only engine path",
			opts: &Options{
				EnginePath: "1/2/3/4",
			},
			expectedRootPath: "1/2/3/4",
		},
	}

	for _, tc := range testCases {
		rootPath, subPath := tc.opts.buildEnginePath()

		assert.Equal(t, tc.expectedRootPath, rootPath, tc.name)
		assert.Equal(t, tc.expectedSubPath, subPath, tc.name)
	}
}
