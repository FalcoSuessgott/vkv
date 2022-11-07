package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestGetCapabilities() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		s        map[string]interface{}
		expected *Capability
	}{
		{
			name:     "root",
			rootPath: "cap",
			subPath:  "secret",
			s: map[string]interface{}{
				"key":  "value",
				"user": "password",
			},
			expected: &Capability{Root: true},
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
		if err = s.v.WriteSecrets(tc.rootPath, tc.subPath, tc.s); err != nil {
			s.Suite.T().Fail()
		}

		caps, err := s.v.GetCapabilities(tc.rootPath)
		assert.NoError(s.Suite.T(), err)

		assert.Equal(s.Suite.T(), tc.expected, caps, tc.name)
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		c        *Capability
		expected string
	}{
		{
			name: "simple",
			c: &Capability{
				Create: true,
				Update: true,
			},
			expected: "✔\t✖\t✔\t✖\t✖\t✖\n",
		},
		{
			name: "simple",
			c: &Capability{
				Create: true,
				Update: true,
				Root:   true,
			},
			expected: "✔\t✖\t✔\t✖\t✖\t✔\n",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.c.String(), tc.name)
	}
}
