package cmd

import (
	"bytes"
	"io"
	"testing"

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
		c := newVersionCmd(tc.version)
		b := bytes.NewBufferString("")

		c.SetOut(b)

		err := c.Execute()
		assert.NoError(t, err)

		out, _ := io.ReadAll(b)
		assert.Equal(t, tc.expected, string(out), tc.name)
	}
}
