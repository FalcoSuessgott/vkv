package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name string
		cmd  []string
		exp  string
		err  bool
	}{
		{
			name: "simple command",
			cmd:  []string{"echo", "hallo"},
			exp:  "hallo\n",
			err:  false,
		},
		{
			name: "pipe command",
			cmd:  []string{"echo", "hallo world", "|", "cut", "-d", "\" \"", "-f2"},
			exp:  "world\n",
			err:  false,
		},
		{
			name: "error command",
			cmd:  []string{"cat invalid_file.txt"},
			err:  true,
		},
	}

	for _, tc := range testCases {
		out, err := Run(tc.cmd)

		if tc.err {
			require.Error(t, err, tc.name)
		}

		assert.Equal(t, tc.exp, string(out), tc.name)
	}
}
