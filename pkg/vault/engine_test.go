package vault

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		s.Run(tc.name, func() {
			if tc.prepare {
				//nolint: errcheck
				s.client.EnableKV2Engine(tc.path)
			}

			err := s.client.EnableKV2EngineErrorIfNotForced(tc.force, tc.path)
			if tc.err {
				require.Error(s.Suite.T(), err, tc.name)
			} else {
				require.NoError(s.Suite.T(), err, tc.name)
			}
		})
	}
}

func (s *VaultSuite) TestListAllKVSecretEngines() {
	testCases := []struct {
		name     string
		ns       []string
		engines  []string
		expected Engines
	}{
		{
			name:    "test",
			ns:      []string{"a", "b"},
			engines: []string{"1", "2", "3"},
			expected: Engines{
				"":  []string{"secret/"}, // enabled by default
				"a": []string{"1/", "2/", "3/"},
				"b": []string{"1/", "2/", "3/"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// create ns
			for _, ns := range tc.ns {
				require.NoError(s.Suite.T(), s.client.CreateNamespaceErrorIfNotForced("", ns, true), tc.name)

				// create engines
				for _, engine := range tc.engines {
					s.client.Client.SetNamespace(ns)
					require.NoError(s.Suite.T(), s.client.EnableKV2Engine(engine), tc.name)
					s.client.Client.ClearNamespace()
				}
			}

			res, err := s.client.ListAllKVSecretEngines("")
			require.NoError(s.Suite.T(), err, tc.name)
			assert.ElementsMatch(s.Suite.T(), []string{"1/", "2/", "3/"}, res["a"], "all namespaces")
			assert.ElementsMatch(s.Suite.T(), []string{"1/", "2/", "3/"}, res["b"], "all namespaces")
		})
	}
}
