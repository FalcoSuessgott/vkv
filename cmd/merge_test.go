package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMergeFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			name: "test: src and dest empty",
			err:  true,
			args: []string{"--src-path", "", "--dest-path", ""},
		},
		{
			name: "test: src and dest equal",
			err:  true,
			args: []string{"--src-path", "a", "--dest-path", "b"},
		},
		{
			name: "test: src missing",
			err:  true,
			args: []string{"--dest-path", "b"},
		},
		{
			name: "test: dest missing",
			err:  true,
			args: []string{"--src-path", "b"},
		},
	}

	for _, tc := range testCases {
		c := mergeCmd()
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
