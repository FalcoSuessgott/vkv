package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateImportFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			name: "test: path is required",
			err:  true,
			args: []string{"-"},
		},
		{
			name: "test: input is required",
			err:  true,
			args: []string{},
		},
		{
			name: "test: file not found",
			err:  true,
			args: []string{"-f", "nonexistingfile.txt"},
		},
		{
			name: "test: cannot specify both",
			err:  true,
			args: []string{"-f", "nonexistingfile.txt", "-"},
		},
	}

	for _, tc := range testCases {
		c := importCmd()
		b := bytes.NewBufferString("")

		c.SetArgs(tc.args)
		c.SetOut(b)

		err := c.Execute()
		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.NoError(t, err, tc.name)
	}
}
