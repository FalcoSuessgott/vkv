package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileFuncs(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		exists   bool
		isJSON   bool
		err      bool
		secret   map[string]interface{}
	}{
		{
			name:     "test: normal json",
			filename: "testdata/normal_json.json",
			exists:   true,
			isJSON:   true,
			secret: map[string]interface{}{
				"kv/secret": map[string]interface{}{
					"username": "username",
					"password": "passw0rd",
				},
			},
		},
		{
			name:     "test: normal json",
			filename: "testdata/normal_yaml.yml",
			exists:   true,
			isJSON:   false,
			secret: map[string]interface{}{
				"kv/secret": map[string]interface{}{
					"username": "username",
					"password": "passw0rd",
				},
			},
		},
		{
			name:     "test: file does not exist",
			filename: "testdata/invalid_file",
			exists:   false,
		},
		{
			name:     "test: neither json or yaml",
			filename: "testdata/invalid_format.md",
			exists:   true,
			isJSON:   false,
			err:      true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.exists, FileExists(tc.filename), "file exists "+tc.name)

		if tc.exists {
			content, err := ReadFile(tc.filename)
			assert.NoError(t, err)

			data, res, err := JSONorYAML(content)

			if tc.err {
				assert.Error(t, err, "file is not json or yaml "+tc.name)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.isJSON, res, "json or yaml "+tc.name)
				assert.Equal(t, tc.secret, data, "data "+tc.name)
			}
		}
	}
}

func TestRemoveEmptyElements(t *testing.T) {
	testCases := []struct {
		name     string
		parts    []string
		expected []string
	}{
		{
			name:     "test: root path",
			parts:    []string{"", "", "1", "2", "", "3"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "test: root path",
			parts:    []string{"1", "2", "3"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "test: root path",
			parts:    []string{"1", "", "3"},
			expected: []string{"1", "3"},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, removeEmptyElements(tc.parts), tc.name)
	}
}

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
			json: []byte("{\n\t\"key_1\": \"value\",\n\t\"key_2\": 12\n}"),
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
			json: []byte("{\n\t\"key_1\": \"value\",\n\t\"key_2\": 12,\n\t\"key_3\": {\n\t\t\"foo\": \"bar\",\n\t\t\"user\": \"password\"\n\t}\n}"),
		},
	}

	for _, tc := range testCases {
		out, err := ToJSON(tc.s)

		if tc.err {
			assert.Error(t, err)
		} else {
			assert.Equal(t, string(tc.json), string(out), tc.name)
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

func TestMergeMaps(t *testing.T) {
	testCases := []struct {
		name                string
		src, dest, expected map[string]interface{}
	}{
		{
			name: "merge with override",
			src: map[string]interface{}{
				"a/dir":  map[string]interface{}{"key": 12},
				"a/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
			dest: map[string]interface{}{
				"c/dir": map[string]interface{}{"key": 9},
			},
			expected: map[string]interface{}{
				"c/dir":  map[string]interface{}{"key": 12},
				"c/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
		},
		{
			name: "merge with override 2",
			src: map[string]interface{}{
				"a/dir":  map[string]interface{}{"key": 12},
				"a/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
			dest: map[string]interface{}{
				"c/dir": map[string]interface{}{"key": 9, "key3": "test"},
			},
			expected: map[string]interface{}{
				"c/dir":  map[string]interface{}{"key": 12, "key3": "test"},
				"c/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
		},
		{
			name: "merge with additional keys",
			src: map[string]interface{}{
				"a/dir":  map[string]interface{}{"key": 12},
				"a/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
			dest: map[string]interface{}{
				"c/dir": map[string]interface{}{"key2": 9},
			},
			expected: map[string]interface{}{
				"c/dir":  map[string]interface{}{"key": 12, "key2": 9},
				"c/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
		},
		{
			name: "empty dest",
			src: map[string]interface{}{
				"a/dir":  map[string]interface{}{"key": 12},
				"a/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
			dest: map[string]interface{}{},
			expected: map[string]interface{}{
				"c/dir":  map[string]interface{}{"key": 12},
				"c/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
		},
		{
			name: "empty src",
			src:  map[string]interface{}{},
			dest: map[string]interface{}{
				"a/dir":  map[string]interface{}{"key": 12},
				"a/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
			expected: map[string]interface{}{
				"c/dir":  map[string]interface{}{"key": 12},
				"c/dir2": map[string]interface{}{"key": "value", "user": "password"},
			},
		},
		{
			name:     "empty both",
			src:      map[string]interface{}{},
			dest:     map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, MergeMaps(tc.src, tc.dest, "c"), tc.name)
	}
}
