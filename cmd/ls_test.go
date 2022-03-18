package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/stretchr/testify/assert"
)

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			name: "test: yaml and json",
			err:  true,
			args: []string{"--json", "--yaml"},
		},
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
			name: "test: export 2",
			err:  true,
			args: []string{"--json", "--markdown"},
		},
		{
			name: "test: export 2",
			err:  true,
			args: []string{"--json", "--yaml"},
		},
	}

	for _, tc := range testCases {
		c := lsCmd()
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

func TestMaxValueLength(t *testing.T) {
	testCases := []struct {
		name           string
		opts           *lsOptions
		maxValueLength string
		expected       int
		err            bool
	}{
		{
			name:           "no env set no flag defined",
			opts:           defaultLSOptions(),
			maxValueLength: "",
			expected:       printer.MaxValueLength,
		},
		{
			name:           "env set no flag defined",
			opts:           defaultLSOptions(),
			maxValueLength: "2",
			expected:       2,
		},
		{
			name: "env set and flag defined",
			opts: &lsOptions{
				maxValueLength: 23,
			},
			maxValueLength: "2",
			expected:       23,
		},
		{
			name: "no env set and flag defined",
			opts: &lsOptions{
				maxValueLength: 23,
			},
			maxValueLength: "",
			expected:       23,
		},
		{
			name: "no env set and flag defined",
			opts: &lsOptions{
				maxValueLength: 23,
			},
			maxValueLength: "",
			expected:       23,
		},
		{
			name:           "invalid env",
			opts:           defaultLSOptions(),
			maxValueLength: "invalid",
			err:            true,
		},
		{
			name:           "length off",
			opts:           defaultLSOptions(),
			maxValueLength: "-1",
			expected:       -1,
		},
	}

	for _, tc := range testCases {
		if tc.maxValueLength != "" {
			os.Setenv(maxValueLengthEnvVar, tc.maxValueLength)
		}

		err := tc.opts.validateFlags()
		if !tc.err {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tc.opts.maxValueLength, tc.name)
		} else {
			assert.Error(t, err)
		}
	}
}
