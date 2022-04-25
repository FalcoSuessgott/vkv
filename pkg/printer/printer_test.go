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
		p := NewPrinter(tc.options...)

		p.maskValues(tc.input)

		assert.Equal(t, tc.output, tc.input, tc.name)
	}
}

//nolint: lll
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
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(false),
			},
			output: "secret/\n├── \n│   ├── key=*****\n│   └── user=********\n└── \n    └── key=**\n",
		},
		{
			name: "test: default opions multiple paths",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
				"secret_2": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(false),
			},
			output: "secret/\n├── \n│   ├── key=*****\n│   └── user=********\n└── \n    └── key=**\nsecret_2/\n├── \n│   ├── key=*****\n│   └── user=********\n└── \n    └── key=**\n",
		},
		{
			name: "test: show secrets",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(true),
			},
			output: "secret/\n├── \n│   ├── key=value\n│   └── user=password\n└── \n    └── key=12\n",
		},
		{
			name: "test: only paths",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyPaths(true),
				ShowValues(true),
			},
			output: "secret/\n├── \n└── \n",
		},
		{
			name: "test: only keys",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyKeys(true),
				ShowValues(true),
			},
			output: "secret/\n├── \n│   ├── key\n│   └── user\n└── \n    └── key\n",
		},
		{
			name: "test: normal map to json",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(JSON),
				ShowValues(true),
			},
			output: "{\n  \"secrets\": {\n    \"key_1\": {\n      \"key\": \"value\",\n      \"user\": \"password\"\n    },\n    \"key_2\": {\n      \"key\": 12\n    }\n  }\n}",
		},
		{
			name: "test: normal map to json only keys",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(JSON),
				OnlyKeys(true),
			},
			output: "{\n  \"secrets\": {\n    \"key_1\": {\n      \"key\": \"\",\n      \"user\": \"\"\n    },\n    \"key_2\": {\n      \"key\": \"\"\n    }\n  }\n}",
		},
		{
			name: "test: normal map to yaml",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(YAML),
				ShowValues(true),
			},
			output: "secrets:\n  key_1:\n    key: value\n    user: password\n  key_2:\n    key: 12\n",
		},
		{
			name: "test: normal map to yaml only keys",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(YAML),
				OnlyKeys(true),
			},
			output: "secrets:\n  key_1:\n    key: \"\"\n    user: \"\"\n  key_2:\n    key: \"\"\n",
		},
		{
			name: "test: export format",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Export),
				ShowValues(true),
			},
			output: "export key=\"value\"\nexport user=\"password\"\nexport key=\"12\"\n",
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
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Markdown),
			},
			output: "|  MOUNT  | PATHS | KEYS |  VALUES  |\n|---------|-------|------|----------|\n| secrets | key_1 | key  | *****    |\n|         |       | user | ******** |\n|         | key_2 | key  | **       |\n",
		},
		{
			name: "test: markdown only keys",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Markdown),
				OnlyKeys(true),
			},
			output: "|  MOUNT  | PATHS | KEYS |\n|---------|-------|------|\n| secrets | key_1 | key  |\n|         |       | user |\n|         | key_2 | key  |\n",
		},
		{
			name: "test: markdown only paths",
			s: map[string]interface{}{
				"secrets": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				ToFormat(Markdown),
				OnlyPaths(true),
			},
			output: "|  MOUNT  | PATHS |\n|---------|-------|\n| secrets | key_1 |\n|         | key_2 |\n",
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)
		assert.NoError(t, p.Out(tc.s))

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
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts:     []Option{},
			expected: []string{"mount", "paths", "keys", "values"},
		},
		{
			name: "only paths",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				OnlyPaths(true),
			},
			expected: []string{"mount", "paths"},
		},
		{
			name: "only keys",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key_1": map[string]interface{}{"key": "value", "user": "password"},
					"key_2": map[string]interface{}{"key": 12},
				},
			},
			opts: []Option{
				OnlyKeys(true),
			},
			expected: []string{"mount", "paths", "keys"},
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)
		headers, _ := p.buildMarkdownTable(tc.s)

		assert.Equal(t, tc.expected, headers, tc.name)
	}
}
