package vault

import (
	"bytes"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/printer"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/andreyvit/diff"
)

func (s *VaultSuite) TestPrinterDefault() {
	testCases := []struct {
		name    string
		printer string
		fOpts   *FormatOptions
		kv      *KVSecrets
		err     bool
	}{
		{
			name:    "default",
			printer: "default",
			kv:      exampleKVSecrets(),
			fOpts: &FormatOptions{
				ShowDiff:    false,
				OnlyKeys:    false,
				MaskSecrets: false,
			},
		},
		{
			name:    "default_diff",
			printer: "default",
			kv:      exampleKVSecrets(),
			fOpts: &FormatOptions{
				ShowDiff: true,
			},
		},
		{
			name:    "default_masked",
			printer: "default",
			kv:      exampleKVSecrets(),
			fOpts: &FormatOptions{
				MaskSecrets:    true,
				MaxValueLength: 12,
			},
		},
		{
			name:    "default_only-keys",
			printer: "default",
			kv:      exampleKVSecrets(),
			fOpts: &FormatOptions{
				OnlyKeys: true,
			},
		},
		{
			name:    "default_only-keys_diff",
			printer: "default",
			kv:      exampleKVSecrets(),
			fOpts: &FormatOptions{
				ShowDiff: true,
				OnlyKeys: true,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// setup buffer for output
			var b bytes.Buffer
			opts := printer.PrinterOptions{}
			opts.Writer = &b
			opts.Format = tc.printer

			// dependency injection
			tc.kv.Vault = s.client

			// disable colored output for test purposes
			s.Suite.T().Setenv("NO_COLOR", "true")
			s.Suite.T().Setenv("NO_HYPERLINKS", "true")

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

func exampleKVSecrets() *KVSecrets {
	return &KVSecrets{
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
				{
					Version: 2,
					Data: map[string]interface{}{
						"foo": "bar",
						"new": "element",
					},
				},
				{
					Version: 3,
					Data: map[string]interface{}{
						"foo": "change",
						"new": "element",
					},
				},
				{
					Version: 4,
					Data:    map[string]interface{}{},
				},
			},
		},
	}
}
