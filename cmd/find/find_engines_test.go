package find

import (
	"bytes"
	"io"
	"strings"
	"testing"

	printer "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestFindEnginesCommand() {
	testCases := []struct {
		name       string
		args       []string
		engines    vault.Engines
		expEngines string
	}{
		{
			name: "find all engines",
			args: []string{"--all=true"},
			engines: vault.Engines{
				"":  []string{},
				"a": []string{"secret_1", "secret_2"},
			},
			expEngines: `secret/
secret_1/
secret_2/
`,
		},
		{
			name: "find ns from a",
			args: []string{"-n=a"},
			engines: vault.Engines{
				"":  []string{},
				"a": []string{"secret_1", "secret_2"},
			},
			expEngines: `secret_1/
secret_2/
`,
		},
		{
			name: "find all ns with regex",
			args: []string{"--all", "--regex=secret"},
			engines: vault.Engines{
				"":  []string{},
				"a": []string{"b", "c"},
			},
			expEngines: `secret/
`,
		},
		{
			name: "find all ns in json",
			args: []string{"--all", "--format=json"},
			engines: vault.Engines{
				"":  []string{},
				"a": []string{"secret_1", "secret_2"},
			},
			expEngines: `{
  "engines": [
    "secret/",
    "secret_1/",
    "secret_2/"
  ]
}
`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// create engines
			for _, ns := range utils.SortMapKeys(utils.ToMapStringInterface(tc.engines)) {
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

				for _, e := range tc.engines[ns] {
					s.client.Client.SetNamespace(ns)
					require.NoError(s.Suite.T(), s.client.EnableKV2Engine(e), tc.name)
					s.client.Client.ClearNamespace()
				}
			}

			// run cmd
			b := bytes.NewBufferString("")

			findCmd := newFindEngineCmd(b, s.client)
			findCmd.SetOut(b)
			findCmd.SetArgs(tc.args)

			require.NoError(s.Suite.T(), findCmd.Execute(), tc.name)

			out, _ := io.ReadAll(b)

			b.Reset()

			// assert
			assert.Equal(s.Suite.T(), tc.expEngines, string(out), tc.name)
		})
	}
}

func TestEnginesOutputFormat(t *testing.T) {
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
		o := &findEnginesOptions{
			FormatString: tc.format,
		}

		err := o.validateFlags()

		if tc.err {
			require.ErrorIs(t, err, printer.ErrInvalidFormat, tc.name)
		} else {
			require.NoError(t, err, tc.name, tc.name)
			assert.Equal(t, tc.expected, o.outputFormat, tc.name)
		}
	}
}
