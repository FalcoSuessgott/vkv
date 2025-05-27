package vault

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestGetCapabilities() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		s        Secrets
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
		s.Run(tc.name, func() {
			// enable kv engine
			require.NoError(s.Suite.T(), s.client.EnableKV2Engine(context.Background(), tc.rootPath))

			// enable kv engine again, so it erros
			require.Error(s.Suite.T(), s.client.EnableKV2Engine(context.Background(), tc.rootPath))

			// read secrets- find none, so it errors
			_, err := s.client.ReadSecrets(context.Background(), tc.rootPath, tc.subPath)
			require.Error(s.Suite.T(), err)

			// actual write the secrets
			if err = s.client.WriteSecrets(context.Background(), tc.rootPath, tc.subPath, tc.s); err != nil {
				s.Suite.T().Fail()
			}

			caps, err := s.client.GetCapabilities(context.Background(), tc.rootPath)
			require.NoError(s.Suite.T(), err)

			assert.Equal(s.Suite.T(), tc.expected, caps, tc.name)
		})
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
