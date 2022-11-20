package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRegex(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		str     string
		match   bool
		err     bool
	}{
		{
			name:    "1",
			pattern: "test*",
			str:     "test/sub",
			match:   true,
		},
		{
			name:    "2",
			pattern: "te",
			str:     "test/sub",
			match:   true,
		},
		{
			name:    "2",
			pattern: "*",
			str:     "test/sub",
			match:   true,
			err:     true,
		},
	}

	for _, tc := range testCases {
		res, err := MatchRegex(tc.pattern, tc.str)

		if tc.err {
			assert.Error(t, err)

			continue
		}

		assert.Equal(t, tc.match, res, tc.name)
	}
}
