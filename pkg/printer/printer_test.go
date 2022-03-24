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
			name:    "test: default opions",
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
		p := NewPrinter(tc.input, tc.options...)

		p.maskSecrets()

		assert.Equal(t, tc.output, p.secrets, tc.name)
	}
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		name   string
		s      map[string]interface{}
		opts   []Option
		output string
	}{
		{
			name: "test: default opions",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ShowSecrets(false),
			},
			output: `key_1
	key=*****
	user=********
key_2
	key=**
`,
		},
		{
			name: "test: show secrets",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ShowSecrets(true),
			},
			output: `key_1
	key=value
	user=password
key_2
	key=12
`,
		},
		{
			name: "test: only paths",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				OnlyPaths(true),
				ShowSecrets(true),
			},
			output: `key_1
key_2
`,
		},
		{
			name: "test: only keys",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				OnlyKeys(true),
				ShowSecrets(true),
			},
			output: `key_1
	key
	user
key_2
	key
`,
		},
		{
			name: "test: normal map to json",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(JSON),
				ShowSecrets(true),
			},
			output: "{\n\t\"key_1\": {\n\t\t\"key\": \"value\",\n\t\t\"user\": \"password\"\n\t},\n\t\"key_2\": {\n\t\t\"key\": 12\n\t}\n}",
		},
		{
			name: "test: normal map to json only keys",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(JSON),
				OnlyKeys(true),
			},
			output: "{\n\t\"key_1\": {\n\t\t\"key\": \"\",\n\t\t\"user\": \"\"\n\t},\n\t\"key_2\": {\n\t\t\"key\": \"\"\n\t}\n}",
		},
		{
			name: "test: normal map to yaml",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(YAML),
				ShowSecrets(true),
			},
			output: "key_1:\n  key: value\n  user: password\nkey_2:\n  key: 12\n",
		},
		{
			name: "test: normal map to yaml only keys",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(YAML),
				OnlyKeys(true),
			},
			output: "key_1:\n  key: \"\"\n  user: \"\"\nkey_2:\n  key: \"\"\n",
		},
		{
			name: "test: export format",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(Export),
				ShowSecrets(true),
			},
			output: "export key=value\nexport user=password\nexport key=12\n",
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
			name: "test: markdown",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(Markdown),
			},
			output: "| PATHS | KEYS |  VALUES  |\n|-------|------|----------|\n| key_1 | key  | *****    |\n|       | user | ******** |\n| key_2 | key  | **       |\n",
		},
		{
			name: "test: markdown",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				ToFormat(Markdown),
			},
			output: "| PATHS | KEYS |  VALUES  |\n|-------|------|----------|\n| key_1 | key  | *****    |\n|       | user | ******** |\n| key_2 | key  | **       |\n",
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.s, tc.opts...)
		assert.NoError(t, p.Out())

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
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts:     []Option{},
			expected: []string{"paths", "keys", "values"},
		},
		{
			name: "only paths",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				OnlyPaths(true),
			},
			expected: []string{"paths"},
		},
		{
			name: "only keys",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			opts: []Option{
				OnlyKeys(true),
			},
			expected: []string{"paths", "keys"},
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.s, tc.opts...)
		headers, _ := p.buildMarkdownTable()

		assert.Equal(t, tc.expected, headers, tc.name)
	}
}
