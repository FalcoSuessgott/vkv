package cmd

import (
	"os"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/stretchr/testify/assert"
)

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name    string
		options *Options
		err     bool
	}{
		{
			name:    "test: valid options",
			options: defaultOptions(),
		},
		{
			name: "test: yaml and json",
			err:  true,
			options: &Options{
				json: true,
				yaml: true,
			},
		},
		{
			name: "test: only keys and only paths",
			err:  true,
			options: &Options{
				onlyKeys:  true,
				onlyPaths: true,
			},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			options: &Options{
				onlyKeys:   true,
				showValues: true,
			},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			options: &Options{
				onlyPaths:  true,
				showValues: true,
			},
		},
		{
			name: "test: export 1",
			err:  true,
			options: &Options{
				export:     true,
				showValues: true,
			},
		},
		{
			name: "test: export 2",
			err:  true,
			options: &Options{
				export: true,
				json:   true,
			},
		},
		{
			name: "test: export 2",
			err:  true,
			options: &Options{
				export: true,
				yaml:   true,
			},
		},
	}

	for _, tc := range testCases {
		err := tc.options.validateFlags()

		if tc.err {
			assert.Error(t, err)

			continue
		}

		assert.NoError(t, err)
	}
}

func TestMaxValueLength(t *testing.T) {
	testCases := []struct {
		name           string
		opts           *Options
		maxValueLength string
		expected       int
		err            bool
	}{
		{
			name:           "no env set no flag defined",
			opts:           defaultOptions(),
			maxValueLength: "",
			expected:       printer.MaxValueLength,
		},
		{
			name:           "env set no flag defined",
			opts:           defaultOptions(),
			maxValueLength: "2",
			expected:       2,
		},
		{
			name: "env set and flag defined",
			opts: &Options{
				maxValueLength: 23,
			},
			maxValueLength: "2",
			expected:       23,
		},
		{
			name: "no env set and flag defined",
			opts: &Options{
				maxValueLength: 23,
			},
			maxValueLength: "",
			expected:       23,
		},
		{
			name: "no env set and flag defined",
			opts: &Options{
				maxValueLength: 23,
			},
			maxValueLength: "",
			expected:       23,
		},
		{
			name:           "invalid env",
			opts:           defaultOptions(),
			maxValueLength: "invalid",
			err:            true,
		},
		{
			name:           "length off",
			opts:           defaultOptions(),
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
