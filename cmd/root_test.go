package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	testCases := []struct {
		name   string
		s      Secrets
		output string
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			output: `key_1:	value
key_2:	12
`,
		},
		{
			name:   "test: empty map",
			s:      map[string]interface{}{},
			output: "",
		},
		// todo
		{
			name: "test: multiple values",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			output: `key_1:	value
key_2:	12
key_3:	map[foo:bar user:password]
`,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.output, tc.s.print(), tc.name)
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		name       string
		s          Secrets
		sortedKeys []string
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			sortedKeys: []string{"key_1", "key_2"},
		},
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_3": nil,
				"key_1": "value",
				"key_2": 12,
			},
			sortedKeys: []string{"key_1", "key_2", "key_3"},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.sortedKeys, tc.s.sortKeys(), tc.name)
	}
}

func TestToJson(t *testing.T) {
	testCases := []struct {
		name string
		s    Secrets
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
		s    Secrets
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
