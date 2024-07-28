package vault

import (
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestGetCapabilities() {
	testCases := []struct {
		name     string
		kv       *KVSecrets
		expected *Capability
	}{
		{
			name:     "root",
			kv:       exampleKVSecrets(false),
			expected: &Capability{true, true, true, true, true, true},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// write secrets
			for path, secrets := range tc.kv.Secrets {
				for _, secret := range secrets {
					s.Suite.Require().NoError(s.client.WriteSecrets(tc.kv.MountPath, path, secret.Data), "writing secrets "+tc.name)
				}
			}

			caps, err := s.client.GetCapabilities("secret/test/admin")
			s.Suite.Require().NoError(err, "reading caps "+tc.name)

			s.Suite.Require().Equal(tc.expected, caps, "comparing caps "+tc.name)
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
	}

	for _, tc := range testCases {
		t.Setenv(utils.NoColorEnv, "true")

		require.Equal(t, tc.expected, tc.c.String(), tc.name)
	}
}
