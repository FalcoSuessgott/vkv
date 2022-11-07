package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintExport(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: export format",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Export),
				ShowValues(true),
			},
			output: `export key='value'
export user='password'
`,
		},
		{
			name: "test: empty export",
			s:    map[string]interface{}{},
			opts: []Option{
				ToFormat(Export),
			},
			output: "",
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath+"/"] = tc.s
		assert.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}
