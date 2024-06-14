package render

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// String renders byte array input with the given data.
func Apply(tmpl []byte, input interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	tpl, err := template.New("template").
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
		Parse(string(tmpl))
	if err != nil {
		return buf, err
	}

	if err := tpl.Execute(&buf, input); err != nil {
		return buf, err
	}

	return buf, nil
}
