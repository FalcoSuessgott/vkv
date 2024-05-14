package namespace

import (
	"bytes"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintNamespaces(t *testing.T) {
	testCases := []struct {
		name     string
		ns       vault.Namespaces
		opts     []Option
		expected string
		err      bool
	}{
		{
			name: "root & 2 sub",
			ns: map[string][]string{
				"": {"a", "b"},
			},
			opts: []Option{
				ToFormat(Base),
			},
			expected: `a
b
`,
		},
		{
			name: "only root",
			ns: map[string][]string{
				"": {},
			},
			opts: []Option{
				ToFormat(Base),
			},
			expected: ``,
		},
		{
			name: "multi leveled",
			ns: map[string][]string{
				"":     {"a", "b"},
				"a":    {"a1", "a2"},
				"a/a1": {"1", "2"},
			},
			opts: []Option{
				ToFormat(Base),
			},
			expected: `a
a/a1
a/a1/1
a/a1/2
a/a2
b
`,
		},
		{
			name: "empty",
			ns:   map[string][]string{},
			opts: []Option{
				ToFormat(Base),
			},
			err: true,
		},
		{
			name: "regex",
			ns: map[string][]string{
				"":  {"a", "b"},
				"a": {"a1", "a2"},
			},
			opts: []Option{
				ToFormat(Base),
				WithRegex("a"),
			},
			expected: `a
a/a1
a/a2
`,
		},
		{
			name: "json",
			ns: map[string][]string{
				"": {"a", "b"},
			},
			opts: []Option{
				ToFormat(JSON),
			},
			expected: `{
  "namespaces": [
    "a",
    "b"
  ]
}
`,
		},
		{
			name: "yaml",
			ns: map[string][]string{
				"": {"a", "b"},
			},
			opts: []Option{
				ToFormat(YAML),
			},
			expected: `namespaces:
- a
- b
`,
		},
		{
			name: "invalid regex",
			ns: map[string][]string{
				"": {"a", "b"},
			},
			opts: []Option{
				ToFormat(YAML),
				WithRegex("*"),
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer

		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewNamespacePrinter(tc.opts...)

		err := p.Out(tc.ns)

		if tc.err {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
			assert.Equal(t, tc.expected, b.String(), tc.name)
		}
	}
}
