package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPath(t *testing.T) {
	testCases := []struct {
		name     string
		elements []string
		expected string
	}{
		{
			name:     "simple paths",
			elements: []string{"root", "sub"},
			expected: "root/sub",
		},
		{
			name:     "simple sub paths",
			elements: []string{"root", "sub", "sub2"},
			expected: "root/sub/sub2",
		},
		{
			name:     "empty sub path",
			elements: []string{"root", "", "sub"},
			expected: "root/sub",
		},
		{
			name:     "trailing slash",
			elements: []string{"root", "sub/", "sub", "demo/"},
			expected: "root/sub/sub/demo",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, buildPath(tc.elements...), tc.name)
	}
}

func TestValidateFlags(t *testing.T) {
	testCases := []struct {
		name    string
		options *Options
		err     bool
	}{
		{
			name:    "test: valid options",
			options: defaultOptions(),
		},
		{
			name: "test: yaml and json",
			err:  true,
			options: &Options{
				json: true,
				yaml: true,
			},
		},
		{
			name: "test: only keys and only paths",
			err:  true,
			options: &Options{
				onlyKeys:  true,
				onlyPaths: true,
			},
		},
		{
			name: "test: only keys and show secrets ",
			err:  true,
			options: &Options{
				onlyKeys:    true,
				showSecrets: true,
			},
		},
		{
			name: "test: only paths and show secrets ",
			err:  true,
			options: &Options{
				onlyPaths:   true,
				showSecrets: true,
			},
		},
	}

	for _, tc := range testCases {
		err := tc.options.validateFlags()

		if tc.err {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
	}
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		name    string
		s       secrets
		options *Options
		output  string
	}{
		{
			name: "test: default opions",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			options: defaultOptions(),
			output: `kv2/
key_1	key=***** user=********
key_2	key=**
`,
		},
		{
			name: "test: show secrets",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			options: &Options{
				rootPath:    "root",
				showSecrets: true,
			},
			output: `root/
key_1	key=value user=password
key_2	key=12
`,
		},
		{
			name: "test: only paths",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			options: &Options{
				rootPath:  "root",
				onlyPaths: true,
			},
			output: `root/
key_1
key_2
`,
		},
		{
			name: "test: only keys",
			s: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": 12},
			},
			options: &Options{
				rootPath: "root",
				onlyKeys: true,
			},
			output: `root/
key_1	key user
key_2	key
`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer

		tc.options.writer = &b

		tc.options.evalModifyFlags(tc.s)

		tc.options.print(tc.s)
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		name       string
		s          secrets
		sortedKeys []string
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1":     map[string]interface{}{"key": "value", "user": "password"},
				"key_2":     map[string]interface{}{"key": 12},
				"key_1/a":   map[string]interface{}{"key": "value", "user": "password"},
				"key_2/b":   map[string]interface{}{"key": 12},
				"key_1/a/c": map[string]interface{}{"key": "value", "user": "password"},
				"key_2/b/d": map[string]interface{}{"key": 12},
			},
			sortedKeys: []string{"key_1", "key_1/a", "key_1/a/c", "key_2", "key_2/b", "key_2/b/d"},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.sortedKeys, sortMapKeys(tc.s), tc.name)
	}
}

func TestToJson(t *testing.T) {
	testCases := []struct {
		name string
		s    secrets
		json string
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			json: "{\"key_1\":\"value\",\"key_2\":12}",
		},
		{
			name: "test: empty map",
			s:    map[string]interface{}{},
			json: "{}",
		},
		{
			name: "test: multiple values",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			json: "{\"key_1\":\"value\",\"key_2\":12,\"key_3\":{\"foo\":\"bar\",\"user\":\"password\"}}",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.json, tc.s.toJSON(), tc.name)
	}
}

func TestToYAML(t *testing.T) {
	testCases := []struct {
		name string
		s    secrets
		yaml string
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			yaml: `key_1: value
key_2: 12
`,
		},
		{
			name: "test: empty map",
			s:    map[string]interface{}{},
			yaml: `{}
`, // this is ok since we stop earlier when there are no secrets read
		},
		{
			name: "test: multiple values ",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			yaml: `key_1: value
key_2: 12
key_3:
  foo: bar
  user: password
`,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.yaml, tc.s.toYAML(), tc.name)
	}
}
