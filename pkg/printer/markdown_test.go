package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintMarkdown(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: markdown",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
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
				"secret": map[string]interface{}{
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
				"secret": map[string]interface{}{
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
