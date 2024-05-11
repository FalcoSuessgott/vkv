package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveExtension(t *testing.T) {
	s := "path/to/file.txt"
	assert.Equal(t, "path/to/file", RemoveExtension(s))
}

func TestRemoveCarriageReturn(t *testing.T) {
	s := "new line\r\n"
	assert.Equal(t, "new line\n", RemoveCarriageReturns(s))
}

func TestTransformMap(t *testing.T) {
	testCases := []struct {
		name     string
		m        map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "2 level",
			m: map[string]interface{}{
				"secret": map[string]interface{}{
					"root": map[string]interface{}{
						"kv":   12,
						"bool": false,
					},
				},
			},
			expected: map[string]interface{}{
				"secret/root": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
		},
		{
			name: "1 level",
			m: map[string]interface{}{
				"secret": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
			expected: map[string]interface{}{
				"secret": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
		},
		{
			name: "3 level",
			m: map[string]interface{}{
				"secret": map[string]interface{}{
					"root": map[string]interface{}{
						"sub": map[string]interface{}{
							"kv":   12,
							"bool": false,
						},
					},
				},
			},
			expected: map[string]interface{}{
				"secret/root/sub": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
		},
		{
			name: "3 multi level",
			m: map[string]interface{}{
				"secret": map[string]interface{}{
					"root": map[string]interface{}{
						"sub": map[string]interface{}{
							"kv":   12,
							"bool": false,
						},
						"sub2": map[string]interface{}{
							"kv":   12,
							"bool": false,
						},
					},
				},
			},
			expected: map[string]interface{}{
				"secret/root/sub": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
				"secret/root/sub2": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
		},
		{
			name: "3 multi level 2",
			m: map[string]interface{}{
				"secret": map[string]interface{}{
					"root": map[string]interface{}{
						"sub": map[string]interface{}{
							"kv":   12,
							"bool": false,
						},
					},
					"root2": map[string]interface{}{
						"kv":   12,
						"bool": false,
					},
				},
			},
			expected: map[string]interface{}{
				"secret/root/sub": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
				"secret/root2": map[string]interface{}{
					"kv":   12,
					"bool": false,
				},
			},
		},
	}

	for _, tc := range testCases {
		res := make(map[string]interface{})

		TransformMap("", tc.m, &res)

		assert.Equal(t, tc.expected, res, tc.expected)
	}
}

func TestPathMap(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		m            map[string]interface{}
		isSecretPath bool
		expected     map[string]interface{}
	}{
		{
			name: "secret path",
			path: "root/sub",
			m: map[string]interface{}{
				"k":  "v",
				"k2": 12,
			},
			isSecretPath: true,
			expected: map[string]interface{}{
				"root/": map[string]interface{}{
					"sub": map[string]interface{}{
						"k":  "v",
						"k2": 12,
					},
				},
			},
		},
		{
			name:         "directory path",
			path:         "root/sub",
			m:            map[string]interface{}{},
			isSecretPath: false,
			expected: map[string]interface{}{
				"root/": map[string]interface{}{
					"sub/": map[string]interface{}{},
				},
			},
		},
		{
			name:         "only root",
			path:         "root",
			m:            map[string]interface{}{},
			isSecretPath: false,
			expected: map[string]interface{}{
				"root/": map[string]interface{}{},
			},
		},
		{
			name:         "only root",
			path:         "",
			m:            map[string]interface{}{},
			isSecretPath: false,
			expected:     map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, PathMap(tc.path, tc.m, tc.isSecretPath), tc.name)
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

func TestRemoveDuplicates(t *testing.T) {
	l := []string{"a", "a", "b", "c", "c"}

	assert.Equal(t, []string{"a", "b", "c"}, RemoveDuplicates(l))
}

// nolint: dupword
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

func TestToMapStringInterface(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
	}{
		{
			name: "test: normal map",
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": false},
			},
			expected: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
				"key_2": map[string]interface{}{"key": false},
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, ToMapStringInterface(tc.input), tc.name)
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
				"key_1":   "value",
				"key_2":   false,
				"special": "passw0rd<",
			},
			json: []byte(`{
  "key_1": "value",
  "key_2": false,
  "special": "passw0rd<"
}
`),
		},
		{
			name: "test: empty map",
			s:    map[string]interface{}{},
			json: []byte(`{}
`),
		},
		{
			name: "test: multiple values",
			s: map[string]interface{}{
				"key_1": "value",
				"key_2": false,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
			json: []byte(`{
  "key_1": "value",
  "key_2": false,
  "key_3": {
    "foo": "bar",
    "user": "password"
  }
}
`),
		},
	}

	for _, tc := range testCases {
		out, err := ToJSON(tc.s)

		if tc.err {
			require.Error(t, err)
		} else {
			assert.Equal(t, string(tc.json), string(out), tc.name)
		}
	}
}

func TestFromJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[string]interface{}
		err      bool
	}{
		{
			name: "test: normal map",
			input: []byte(`{
  "key_1": "value",
  "key_2": false
}`),
			expected: map[string]interface{}{
				"key_1": "value",
				"key_2": false,
			},
		},
		{
			name:     "test: empty map",
			input:    []byte("{}"),
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		out, err := FromJSON(tc.input)

		if tc.err {
			require.Error(t, err)
		} else {
			assert.Equal(t, tc.expected, out, tc.name)
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
			require.Error(t, err)
		} else {
			assert.Equal(t, tc.yaml, out, tc.name)
		}
	}
}

func TestFromYAML(t *testing.T) {
	testCases := []struct {
		name     string
		expected map[string]interface{}
		input    []byte
		err      bool
	}{
		{
			name: "test: normal map",
			input: []byte(`key_1: value
key_2: false
`),
			expected: map[string]interface{}{
				"key_1": "value",
				"key_2": false,
			},
		},
		{
			name: "test: multiple values ",
			input: []byte(`key_1: value
key_2: false
key_3:
  foo: bar
  user: password
`),
			expected: map[string]interface{}{
				"key_1": "value",
				"key_2": false,
				"key_3": map[string]interface{}{"foo": "bar", "user": "password"},
			},
		},
	}

	for _, tc := range testCases {
		out, err := FromYAML(tc.input)

		if tc.err {
			require.Error(t, err)
		} else {
			assert.Equal(t, tc.expected, out, tc.name)
		}
	}
}

func TestMergeMap(t *testing.T) {
	testCases := []struct {
		name             string
		m1, m2, expected map[string]interface{}
	}{
		{
			name: "simple maps",
			m1: map[string]interface{}{
				"a": "b",
				"c": 12,
			},
			m2: map[string]interface{}{
				"d": map[string]interface{}{
					"12": false,
				},
			},
			expected: map[string]interface{}{
				"a": "b",
				"c": 12,
				"d": map[string]interface{}{
					"12": false,
				},
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, DeepMergeMaps(tc.m1, tc.m2), tc.name)
	}
}

func TestHandleEnginePath(t *testing.T) {
	testCases := []struct {
		name             string
		enginePath       string
		path             string
		expectedRootPath string
		expectedSubPath  string
	}{
		{
			name:             "only path",
			path:             "1/2/3/4",
			expectedRootPath: "1",
			expectedSubPath:  "2/3/4",
		},
		{
			name:             "one element path",
			path:             "1",
			expectedRootPath: "1",
			expectedSubPath:  "",
		},
		{
			name:             "engine path and path",
			enginePath:       "1/2/3/4",
			path:             "5/6",
			expectedRootPath: "1/2/3/4",
			expectedSubPath:  "5/6",
		},
		{
			name:             " only engine path",
			enginePath:       "1/2/3/4",
			expectedRootPath: "1/2/3/4",
		},
	}

	for _, tc := range testCases {
		rootPath, subPath := HandleEnginePath(tc.enginePath, tc.path)

		assert.Equal(t, tc.expectedRootPath, rootPath, tc.name)
		assert.Equal(t, tc.expectedSubPath, subPath, tc.name)
	}
}

func TestParseEnvs(t *testing.T) {
	type test struct {
		Test string `env:"TEST"`
	}

	o := &test{}
	exp := "test"

	t.Setenv("test_TEST", exp)

	require.NoError(t, ParseEnvs("test_", o))
	require.Equal(t, exp, o.Test)
}
