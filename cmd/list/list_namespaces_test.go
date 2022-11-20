package list

import (
	"bytes"
	"io"
	"strings"
	"testing"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/namespace"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestListNamespacesCommand() {
	testCases := []struct {
		name  string
		args  []string
		ns    vault.Namespaces
		expNS string
	}{
		{
			name: "list all ns",
			args: []string{"--all"},
			ns: vault.Namespaces{
				"":  []string{},
				"a": []string{},
				"b": []string{},
			},
			expNS: `a
b
`,
		},
		{
			name: "list ns from a",
			args: []string{"-n=a"},
			ns: vault.Namespaces{
				"":  []string{},
				"a": []string{},
				"b": []string{},
			},
			expNS: `a
`,
		},
		{
			name: "list all ns with regex",
			args: []string{"--regex=a"},
			ns: vault.Namespaces{
				"":  []string{},
				"a": []string{},
				"b": []string{},
			},
			expNS: `a
`,
		},
		{
			name: "list all ns in json",
			args: []string{"--format=json"},
			ns: vault.Namespaces{
				"":  []string{},
				"a": []string{},
				"b": []string{},
			},
			expNS: `{
  "namespaces": [
    "a",
    "b"
  ]
}
`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// create engines
			for _, ns := range utils.SortMapKeys(utils.ToMapStringInterface(tc.ns)) {
				nsParts := strings.Split(ns, "/")
				nsParent := strings.Join(nsParts[:len(nsParts)-1], "/")
				nsName := nsParts[len(nsParts)-1]

				if ns != "" {
					if nsName != "" {
						require.NoError(s.Suite.T(), s.client.CreateNamespaceErrorIfNotForced(nsParent, nsName, false), tc.name)
					} else {
						require.NoError(s.Suite.T(), s.client.CreateNamespaceErrorIfNotForced("", nsParent, false), tc.name)
					}
				}
			}

			// run cmd
			b := bytes.NewBufferString("")

			listCmd := newListNamespacesCmd(b, s.client)
			listCmd.SetOut(b)
			listCmd.SetArgs(tc.args)

			require.NoError(s.Suite.T(), listCmd.Execute(), tc.name)

			out, _ := io.ReadAll(b)

			b.Reset()

			// assert
			assert.Equal(s.Suite.T(), tc.expNS, string(out), tc.name)
		})
	}
}

func TestNamespaceOutputFormat(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		expected printer.OutputFormat
		err      bool
	}{
		{
			name:     "json",
			err:      false,
			format:   "json",
			expected: printer.JSON,
		},
		{
			name:     "yaml",
			err:      false,
			format:   "YamL",
			expected: printer.YAML,
		},
		{
			name:     "yml",
			err:      false,
			format:   "yml",
			expected: printer.YAML,
		},
		{
			name:     "invalid",
			err:      true,
			format:   "invalid",
			expected: printer.YAML,
		},
		{
			name:     "base",
			err:      false,
			format:   "base",
			expected: printer.Base,
		},
	}

	for _, tc := range testCases {
		o := &listNamespaceOptions{
			FormatString: tc.format,
		}

		err := o.validateFlags()

		if tc.err {
			assert.ErrorIs(t, err, printer.ErrInvalidFormat, tc.name)
		} else {
			assert.NoError(t, err, tc.name, tc.name)
			assert.Equal(t, tc.expected, o.outputFormat, tc.name)
		}
	}
}
