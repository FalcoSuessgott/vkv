package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskSecrets(t *testing.T) {
	testCases := []struct {
		name    string
		options []Option
		input   map[string]interface{}
		output  map[string]interface{}
	}{
		{
			name:    "test: normal secrets",
			options: nil,
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "*****", "user": "********"},
			},
		},
		{
			name:    "test: default options",
			options: nil,
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": 12, "user": false},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "**", "user": "*****"},
			},
		},
		{
			name:    "test: hit password length",
			options: []Option{CustomValueLength(3)},
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": 12, "user": "12345"},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "**", "user": "***"},
			},
		},
	}

	for _, tc := range testCases {
		p := NewPrinter(tc.options...)

		p.maskValues(tc.input)

		assert.Equal(t, tc.output, tc.input, tc.name)
	}
}

func TestPrint(t *testing.T) {
	//nolint: lll
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: default options",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(false),
			},
			output: `root/
└── secret
    ├── key=*****
    └── user=********
`,
		},
		{
			name:     "test: show secrets",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(true),
			},
			output: `root/
└── secret
    ├── key=value
    └── user=password
`,
		},
		{
			name:     "test: only paths",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyPaths(true),
			},
			output: `root/
└── secret
`,
		},
		{
			name:     "test: only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyKeys(true),
			},
			output: `root/
└── secret
    ├── key
    └── user
`,
		},
		{
			name:     "test: normal map to json",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(JSON),
				ShowValues(true),
			},
			output: `{
  "root": {
    "root/secret": {
      "key": "value",
      "user": "password"
    }
  }
}
`,
		},
		{
			name:     "test: normal map to json only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(JSON),
				OnlyKeys(true),
			},
			output: `{
  "root": {
    "root/secret": {
      "key": "",
      "user": ""
    }
  }
}
`,
		},
		{
			name:     "test: normal map to yaml",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(YAML),
				ShowValues(true),
			},
			output: `root:
  root/secret:
    key: value
    user: password

`,
		},
		{
			name:     "test: normal map to yaml only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(YAML),
				OnlyKeys(true),
			},
			output: `root:
  root/secret:
    key: ""
    user: ""

`,
		},
		{
			name:     "test: export format",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Export),
				ShowValues(true),
			},
			output: `export key="value"
export user="password"
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
		{
			name:     "test: markdown",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Markdown),
			},
			output: `|    PATH     | KEY  |  VALUE   |
|-------------|------|----------|
| root/secret | key  | *****    |
|             | user | ******** |
`,
		},
		{
			name:     "test: markdown only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Markdown),
				OnlyKeys(true),
			},
			output: `|    PATH     | KEY  |
|-------------|------|
| root/secret | key  |
|             | user |
`,
		},
		{
			name:     "test: markdown only paths",
			rootPath: "root",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Markdown),
				OnlyPaths(true),
			},
			output: `|    PATH     |
|-------------|
| root/secret |
`,
		},
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
				WithTemplate(`{{ range $entry := . }}{{ printf "%s:\t%s=%v\n" $entry.Path $entry.Key $entry.Value }}{{ end }}`, ""),
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
				WithTemplate(`{{ range $entry := . }}{{ printf "%s:\t%s=%v\n" $entry.Path $entry.Key $entry.Value }}{{ end }}`, ""),
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
    capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "root/secret/*" {
    capabilities = [ "create", "read", "update", "delete", "list" ]
}


`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath] = tc.s
		assert.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}

func TestMarkdownHeader(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		opts     []Option
		expected []string
	}{
		{
			name: "default",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts:     []Option{},
			expected: []string{"path", "key", "value"},
		},
		{
			name: "only paths",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				OnlyPaths(true),
			},
			expected: []string{"path"},
		},
		{
			name: "only keys",
			s: map[string]interface{}{
				"root/secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				OnlyKeys(true),
			},
			expected: []string{"path", "key"},
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		m := map[string]interface{}{}
		m["root"] = tc.s

		p := NewPrinter(tc.opts...)
		headers, _ := p.buildMarkdownTable(m)

		assert.Equal(t, tc.expected, headers, tc.name)
	}
}
