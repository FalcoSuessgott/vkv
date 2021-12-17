package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name    string
		options *lsOptions
		err     bool
	}{
		{
			name:    "test: valid options",
			options: defaultLsOptions(),
		},
		{
			name: "test: yaml and json",
			err:  true,
			options: &lsOptions{
				json: true,
				yaml: true,
			},
		},
		{
			name: "test: only keys and only paths",
			err:  true,
			options: &lsOptions{
				onlyKeys:  true,
				onlyPaths: true,
			},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			options: &lsOptions{
				onlyKeys:    true,
				showSecrets: true,
			},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			options: &lsOptions{
				onlyPaths:   true,
				showSecrets: true,
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
