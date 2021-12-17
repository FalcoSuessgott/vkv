package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	testCases := []struct {
		name            string
		path            string
		expectedRoot    string
		expectedSubPath string
	}{
		{
			name:            "test: root path",
			path:            "kv",
			expectedRoot:    "kv",
			expectedSubPath: "",
		},
		{
			name:            "test: sub path",
			path:            "kv/sub",
			expectedRoot:    "kv",
			expectedSubPath: "sub",
		},
		{
			name:            "test: sub sub path",
			path:            "kv/sub/sub2",
			expectedRoot:    "kv",
			expectedSubPath: "sub/sub2",
		},
		{
			name:            "test: empty path",
			path:            "",
			expectedRoot:    "",
			expectedSubPath: "",
		},
	}

	for _, tc := range testCases {
		rootPath, subPath := SplitPath(tc.path)

		assert.Equal(t, tc.expectedRoot, rootPath, tc.name)
		assert.Equal(t, tc.expectedSubPath, subPath, tc.name)
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		name       string
		s          map[string]interface{}
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
		assert.Equal(t, tc.sortedKeys, SortMapKeys(tc.s), tc.name)
	}
}

func TestToJson(t *testing.T) {
	testCases := []struct {
		name string
		s    map[string]interface{}
		json []byte
		err  bool
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			json: []byte("{\"key_1\":\"value\",\"key_2\":12}"),
		},
		{
			name: "test: empty map",
			s:    map[string]interface{}{},
			json: []byte("{}"),
		},
		{
			name: "test: multiple values",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			json: []byte("{\"key_1\":\"value\",\"key_2\":12,\"key_3\":{\"foo\":\"bar\",\"user\":\"password\"}}"),
		},
	}

	for _, tc := range testCases {
		out, err := ToJSON(tc.s)

		if tc.err {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tc.json, out, tc.name)
		}
	}
}

func TestToYAML(t *testing.T) {
	testCases := []struct {
		name string
		s    map[string]interface{}
		yaml []byte
		err  bool
	}{
		{
			name: "test: normal map",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
			},
			yaml: []byte(`key_1: value
key_2: 12
`),
		},
		{
			name: "test: empty map",
			s:    map[string]interface{}{},
			yaml: []byte(`{}
`), // this is ok since we stop earlier when there are no secrets read
		},
		{
			name: "test: multiple values ",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": 12,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			yaml: []byte(`key_1: value
key_2: 12
key_3:
  foo: bar
  user: password
`),
		},
	}

	for _, tc := range testCases {
		out, err := ToYAML(tc.s)

		if tc.err {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tc.yaml, out, tc.name)
		}
	}
}
