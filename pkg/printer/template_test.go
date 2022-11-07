package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: template",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Template),
				WithTemplate(`{{ range $path, $data := . }}{{ range $entry := $data }}{{ printf "%s:\t%s=%v\n" $path $entry.Key $entry.Value }}{{ end }}{{ end }}`, ""),
			},
			output: `root/secret:	key=*****
root/secret:	user=********

`,
		},
		{
			name:     "test: template show values",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Template),
				ShowValues(true),
				WithTemplate(`{{ range $path, $data := . }}{{ range $entry := $data }}{{ printf "%s:\t%s=%v\n" $path $entry.Key $entry.Value }}{{ end }}{{ end }}`, ""),
			},
			output: `root/secret:	key=value
root/secret:	user=password

`,
		},
		{
			name:     "test: template file show values",
			rootPath: "root",
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
			output: `
path "root/secret/*" {
    capabilities = [ "create", "read" ]
}


`,
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
