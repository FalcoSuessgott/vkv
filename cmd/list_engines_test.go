package cmd

import (
	"bytes"
	"io"

	prt "github.com/FalcoSuessgott/vkv/pkg/printer/engine"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
)

func (s *VaultSuite) TestListEnginesCommand() {
	testCases := []struct {
		name       string
		args       []string
		engines    vault.Engines
		expEngines string
	}{
		{
			name: "list all engines",
			args: []string{"--all=true"},
			engines: vault.Engines{
				"": []string{"secret_1", "secret_2"},
			},
			expEngines: `secret/
secret_1/
secret_2/
`,
		},
		{
			name: "list all engines in json",
			args: []string{"--format=json"},
			engines: vault.Engines{
				"": []string{"secret_1", "secret_2"},
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
			b := bytes.NewBufferString("")
			writer = b

			for _, e := range tc.engines[""] {
				s.Require().NoError(vaultClient.EnableKV2Engine(e), tc.name)
			}

			// run cmd
			listCmd := newListEngineCmd()
			listCmd.SetArgs(tc.args)

			s.Require().NoError(listCmd.Execute(), tc.name)

			out, _ := io.ReadAll(b)

			b.Reset()

			s.Require().Equal(tc.expEngines, string(out), tc.name)
		})
	}
}

func (s *VaultSuite) TestEnginesOutputFormat() {
	testCases := []struct {
		name     string
		expected prt.OutputFormat
		err      bool
	}{
		{
			name:     "json",
			expected: prt.JSON,
		},
		{
			name:     "yaml",
			expected: prt.YAML,
		},
		{
			name:     "yml",
			expected: prt.YAML,
		},
		{
			name:     "invalid",
			err:      true,
			expected: prt.YAML,
		},
		{
			name:     "base",
			expected: prt.Base,
		},
	}

	for _, tc := range testCases {
		o := &listEnginesOptions{
			FormatString: tc.name,
		}

		err := o.Validate(nil, nil)

		s.Require().Equal(tc.err, err != nil, "error "+tc.name)

		// if no error -> assert output format
		if !tc.err {
			s.Require().Equal(tc.expected, o.outputFormat, "format "+tc.name)
		}
	}
}
