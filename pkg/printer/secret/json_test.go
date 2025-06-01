package secret

import (
	"bytes"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintJSON(t *testing.T) {
	testCases := []struct {
		name     string
		s        vault.Secrets
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "normal map to json",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(JSON),
				ShowValues(true),
			},
			output: `{
  "root/": {
    "secret": {
      "key": "value",
      "user": "password"
    }
  }
}
`,
		},
		{
			name:     "normal map to json only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(JSON),
				OnlyKeys(true),
			},
			output: `{
  "root/": {
    "secret": {
      "key": "",
      "user": ""
    }
  }
}
`,
		},
		{
			name:     "merge paths",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(JSON),
				MergePaths(true),
			},
			output: `{
  "root/secret": {
    "key": "value",
    "user": "password"
  }
}
`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewSecretPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath+"/"] = tc.s
		require.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}
