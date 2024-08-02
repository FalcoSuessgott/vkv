package vault

import (
	"bytes"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/andreyvit/diff"
)

func (s *VaultSuite) TestPrinterFormats() {
	testCases := []struct {
		name   string
		format string
		fOpts  *FormatOptions
		kv     *KVSecrets
		err    bool
	}{
		{
			name:   "default",
			format: "default",
			kv:     exampleKVSecrets(false),
			fOpts:  &FormatOptions{},
		},
		{
			name:   "default_masked",
			format: "default",
			kv:     exampleKVSecrets(false),
			fOpts:  &FormatOptions{maskSecrets: true},
		},
		{
			name:   "full",
			format: "full",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{},
		},
		{
			name:   "full_masked",
			format: "full",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{maskSecrets: true},
		},
		{
			name:   "full_masked-diff",
			format: "full",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{maskSecrets: true, showDiff: true},
		},
		{
			name:   "json",
			format: "json",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{},
		},
		// TODO: yaml is marshalled unpredictably
		// {
		// 	name:   "yaml",
		// 	format: "yaml",
		// 	kv:     exampleKVSecrets(),
		// 	fOpts:  &FormatOptions{},
		// },
		{
			name:   "export",
			format: "export",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{},
		},
		{
			name:   "policy",
			format: "policy",
			kv:     exampleKVSecrets(true),
			fOpts:  &FormatOptions{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// setup buffer for output
			var b bytes.Buffer
			opts := printer.PrinterOptions{}
			opts.Writer = &b
			opts.Format = tc.format

			// dependency injection
			tc.kv.Vault = s.client

			// disable colored output for test purposes
			s.Suite.T().Setenv(utils.NoColorEnv, "true")
			s.Suite.T().Setenv(utils.NoHyperlinksEnv, "true")
			s.Suite.T().Setenv(utils.MaxValueLengthEnv, "-1")

			// run printer
			err := printer.Print(tc.kv.PrinterFuncs(tc.fOpts), opts)

			// assertions
			if tc.err {
				s.Require().Error(err, tc.name)
			} else {
				s.Require().NoError(err, tc.name)

				exp, err := os.ReadFile("./testdata/" + tc.name + ".txt")
				s.Require().NoError(err, "golden file "+tc.name)

				if string(exp) != b.String() {
					s.Suite.T().Errorf("diff:\n%s\nwant:\n%s\ngot:\n%stest: %s",
						diff.LineDiff(string(exp), b.String()),
						string(exp),
						b.String(),
						tc.name,
					)
				}
			}
		})
	}
}

func exampleKVSecrets(allVersion bool) *KVSecrets {
	s := &KVSecrets{
		MountPath:   utils.NormalizePath("secret"),
		Type:        "kvv2",
		Description: "test",
		Secrets: map[string][]*Secret{
			"secret/test/admin": {
				{
					Version:        1,
					CustomMetadata: map[string]interface{}{"key": "value"},
					Data: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},
	}

	if allVersion {
		s.Secrets["secret/test/admin"] = append(
			s.Secrets["secret/test/admin"],
			&Secret{
				Version: 2,
				Data: map[string]interface{}{
					"foo": "bar",
					"new": "element",
				},
			},
			&Secret{
				Version: 3,
				Data: map[string]interface{}{
					"foo": "change",
					"new": "element",
				},
			},
			&Secret{
				Version: 4,
				Data:    map[string]interface{}{},
				Deleted: true,
			},
			&Secret{
				Version:   5,
				Data:      map[string]interface{}{},
				Destroyed: true,
			})
	}

	return s
}
