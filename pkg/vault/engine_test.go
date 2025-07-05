package vault

import (
	"context"

	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestGetDescription() {
	s.Run("description", func() {
		desc, err := s.client.GetEngineDescription(s.T().Context(), "secret")

		s.Require().NoError(err)

		s.Require().Equal("key/value secret storage", desc)
	})
}

func (s *VaultSuite) TestGetEngineVersionType() {
	s.Run("description", func() {
		engineType, version, err := s.client.GetEngineTypeVersion(s.T().Context(), "secret")

		s.Require().NoError(err)

		s.Require().Equal("kv", engineType)
		s.Require().Equal("2", version)
	})
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
		s.Run(tc.name, func() {
			if tc.prepare {
				s.Require().NoError(s.client.EnableKV2Engine(context.Background(), tc.path))
			}

			err := s.client.EnableKV2EngineErrorIfNotForced(context.Background(), tc.force, tc.path)

			s.Require().Equal(tc.err, err != nil, tc.name)
		})
	}
}

func (s *VaultSuite) TestListAllKVSecretEngines() {
	testCases := []struct {
		name     string
		engines  []string
		expected Engines
	}{
		{
			name:    "test",
			engines: []string{"1", "2", "3"},
			expected: Engines{
				"": []string{"secret/", "1/", "2/", "3/"}, // enabled by default
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			for _, engine := range tc.engines {
				require.NoError(s.T(), s.client.EnableKV2Engine(context.Background(), engine), tc.name)
			}

			res, err := s.client.ListAllKVSecretEngines(context.Background(), "")
			s.Require().NoError(err, tc.name)

			s.Require().ElementsMatch(tc.expected[""], res[""], tc.name)
		})
	}
}
