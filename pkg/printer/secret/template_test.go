package secret

import (
	"bytes"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		s        vault.Secrets
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name: "test: template",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Template),
				ShowValues(true),
				WithTemplate(`{{ range $path, $secret:= . }}
{{- range $key, $value := $secret -}}
{{ $key}}={{ $value }}
{{ end -}}
{{ end -}}`, ""),
			},
			output: `key=value
user=password
`,
		},
		{
			name: "test: template file show values",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Template),
				ShowValues(true),
				WithTemplate("", "testdata/policies.tmpl"),
			},
			output: `path "root/secret/*" {
    capabilities = [ "create", "read" ]
}
`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewSecretPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath] = tc.s
		require.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, utils.RemoveCarriageReturns(b.String()), tc.name)
	}
}
