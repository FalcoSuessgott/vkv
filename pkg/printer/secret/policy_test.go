package secret

import (
	"bytes"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintPolicy(t *testing.T) {
	testCases := []struct {
		name   string
		caps   map[string]*vault.Capability
		opts   []Option
		output string
	}{
		{
			name: "test: default options",
			caps: map[string]*vault.Capability{
				"root": {
					Read: true,
					Root: true,
				},
			},
			opts: []Option{
				ToFormat(Policy),
				ShowValues(false),
			},
			output: "root\t✖\t✔\t✖\t✖\t✖\t✔\n",
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewSecretPrinter(tc.opts...)

		require.NoError(t, p.printCapabilities(tc.caps))

		expected := header + tc.output
		assert.Equal(t, expected, b.String(), tc.name)
	}
}
