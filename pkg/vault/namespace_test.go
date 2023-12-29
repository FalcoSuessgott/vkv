package vault

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint: errcheck
func (s *VaultSuite) TestNamespaces() {
	testCases := []struct {
		name string
		ns   Namespaces
	}{
		{
			name: "ns",
			ns: Namespaces{
				"":     {"a/", "b/"},
				"a":    {"a1/", "a2/"},
				"a/a1": {"b1/", "b2/"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			for k, v := range tc.ns {
				for _, n := range v {
					require.NoError(s.Suite.T(), s.client.CreateNamespaceErrorIfNotForced(k, n, true), tc.name)
				}
			}

			res, err := s.client.ListAllNamespaces("")
			require.NoError(s.Suite.T(), err)

			for parent, subs := range res {
				assert.ElementsMatch(s.Suite.T(), tc.ns[parent], subs, tc.name)
			}

			for _, ns := range res["a/a1"] {
				require.NoError(s.Suite.T(), s.client.DeleteNamespace("a/a1", ns))
			}
		})
	}
}
