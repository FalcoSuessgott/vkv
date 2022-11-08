package vault

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestWriteReadSecrets() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "secret",
			secrets: map[string]interface{}{
				"user": "password",
			},
		},
		{
			name:     "test: multiple secrets",
			rootPath: "kv",
			subPath:  "secret",
			secrets: map[string]interface{}{
				"user":  "password",
				"value": "42",
			},
		},
		{
			name:     "test: sub path secrets",
			rootPath: "kv",
			subPath:  "secret/sub",
			secrets: map[string]interface{}{
				"user":  "password",
				"value": "42",
			},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		// enable kv engine again, so it erros
		assert.Error(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		// read secrets- find none, so it errors
		_, err := s.v.ReadSecrets(tc.rootPath, tc.subPath)
		assert.Error(s.Suite.T(), err)

		// actual write the secrets
		err = s.v.WriteSecrets(tc.rootPath, tc.subPath, tc.secrets)
		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)

			// read them, expect the exact same secrets as written before
			readSecrets, err := s.v.ReadSecrets(tc.rootPath, tc.subPath)
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.secrets, readSecrets, tc.name)
		}
		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
	}
}

func (s *VaultSuite) TestListSecrets() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
		expected []string
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub2": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: []string{"sub", "sub2"},
		},
		{
			name:     "test: multiple dirs",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub/sub2/sub3": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: []string{"sub", "sub/"},
		},
		{
			name:     "test: empty",
			rootPath: "kv",
			subPath:  "subpath",
			secrets:  nil,
			err:      true,
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		for k, v := range tc.secrets {
			m, ok := v.(map[string]interface{})
			if ok {
				assert.NoError(s.Suite.T(), s.v.WriteSecrets(tc.rootPath, path.Join(tc.subPath, k), m))
			} else {
				fmt.Println("no")
			}
		}

		// read them, expect the exact same secrets as written before
		elements, err := s.v.ListKeys(tc.rootPath, tc.subPath)

		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.expected, elements, tc.name)
		}

		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
	}
}

func (s *VaultSuite) TestEnableKV2EngineErrorIfNotForced() {
	testCases := []struct {
		name    string
		force   bool
		path    string
		prepare bool
		err     bool
	}{
		{
			name:  "engine does not exist, no force",
			force: false,
			path:  "case-1",
			err:   false,
		},
		{
			name:    "engine does exist, no force",
			force:   false,
			prepare: true,
			path:    "case-2",
			err:     true,
		},
		{
			name:    "engine does exist, force",
			force:   true,
			prepare: true,
			path:    "case-3",
			err:     false,
		},
		{
			name:    "engine does exist, no force",
			force:   false,
			prepare: true,
			path:    "case-4",
			err:     true,
		},
	}

	for _, tc := range testCases {
		if tc.prepare {
			//nolint: errcheck
			s.v.EnableKV2Engine(tc.path)
		}

		err := s.v.EnableKV2EngineErrorIfNotForced(tc.force, tc.path)
		if tc.err {
			require.Error(s.Suite.T(), err, tc.name)

			continue
		}

		require.NoError(s.Suite.T(), err, tc.name)
	}
}

func (s *VaultSuite) TestListRecursive() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		secrets  map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "test: simple secret",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
				"sub2": map[string]interface{}{
					"user":  "password",
					"user1": "password",
					"user2": "password",
				},
			},
			expected: map[string]interface{}{
				"kv": Secrets{
					"sub": map[string]interface{}{
						"user":  "password",
						"user1": "password",
						"user2": "password",
					},
					"sub2": map[string]interface{}{
						"user":  "password",
						"user1": "password",
						"user2": "password",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		// enable kv engine
		assert.NoError(s.Suite.T(), s.v.EnableKV2Engine(tc.rootPath))

		for k, v := range tc.secrets {
			if m, ok := v.(map[string]interface{}); ok {
				assert.NoError(s.Suite.T(), s.v.WriteSecrets(tc.rootPath, path.Join(tc.subPath, k), m))
			}
		}

		// read them, expect the exact same secrets as written before
		secrets := map[string]interface{}{}
		tmp, err := s.v.ListRecursive(tc.rootPath, tc.subPath)
		secrets[tc.rootPath] = *tmp

		assert.NoError(s.Suite.T(), err)

		_, err = json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}

		if tc.err {
			assert.Error(s.Suite.T(), err)
		} else {
			assert.NoError(s.Suite.T(), err)
			assert.Equal(s.Suite.T(), tc.expected, secrets, tc.name)
		}

		// disable kv engine, expect no error
		assert.NoError(s.Suite.T(), s.v.DisableKV2Engine(tc.rootPath))
	}
}
