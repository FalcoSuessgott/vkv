package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		content []byte
		err     bool
	}{
		{
			name:    "valid",
			path:    "testdata/file_1.txt",
			content: []byte("Hello World"),
			err:     false,
		},
		{
			name:    "invalid",
			path:    "testdata/invalid",
			content: nil,
			err:     true,
		},
	}

	for _, tc := range testCases {
		out, err := ReadFile(tc.path)

		if tc.err {
			assert.Error(t, err, tc.name)

			continue
		}

		assert.Equal(t, tc.content, out, tc.name)
	}
}

func TestCreateDirectory(t *testing.T) {
	assert.NoError(t, CreateDirectory("a/b/c"))

	_, err := os.Stat("a/b/c")
	assert.NoError(t, err)
}
