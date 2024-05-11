package secret

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintJSON(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: normal map to json",
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
}`,
		},
		{
			name:     "test: normal map to json only keys",
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
}`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath+"/"] = tc.s
		require.NoError(t, p.Out(tc.rootPath, m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}
