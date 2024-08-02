package vault

import (
	"fmt"
	"path"
)

func (s *VaultSuite) TestListRecursive() {
	testCases := []struct {
		name     string
		kv       *KVSecrets
		err      bool
		v1       bool
		expected map[string][]*Secret
	}{
		{
			name: "simple secret",
			kv:   exampleKVSecrets(true),
			expected: map[string][]*Secret{
				"secret/test/admin": {
					{
						Version:        1,
						CustomMetadata: map[string]interface{}{"key": "value"},
						Data: map[string]interface{}{
							"foo": "bar",
						},
					},
					{
						Version: 2,
						Data: map[string]interface{}{
							"foo": "bar",
							"new": "element",
						},
					},
					{
						Version: 3,
						Data: map[string]interface{}{
							"foo": "change",
							"new": "element",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// write secrets
			for path, secrets := range tc.kv.Secrets {
				for _, secret := range secrets {
					if !secret.Deleted && !secret.Destroyed && len(secret.Data) > 0 {
						s.Suite.Require().NoError(s.client.WriteSecrets(tc.kv.MountPath, path, secret.Data), "writing secrets "+tc.name)
					}
				}
			}

			// read secrets
			data, err := s.client.NewKVSecrets(tc.kv.MountPath, "", false, true)
			s.Require().NoError(err, "reading secrets "+tc.name)

			fmt.Println(data.Secrets)

			// assert
			for p, secrets := range tc.expected {
				_, ok := data.Secrets[path.Join(tc.kv.MountPath, p)]

				s.Require().True(ok, "matching paths "+tc.name)

				for i, version := range secrets {
					s.Require().Equal(version.Data, data.Secrets[path.Join(tc.kv.MountPath, p)][i].Data, "matching data "+tc.name)
				}
			}
		})
	}
}

func (s *VaultSuite) TestGetDescription() {
	s.Run("description", func() {
		desc, err := s.client.GetEngineDescription("secret")

		s.Require().NoError(err)

		s.Require().Equal("key/value secret storage", desc)
	})
}

func (s *VaultSuite) TestGetEngineVersionType() {
	s.Run("description", func() {
		engineType, version, err := s.client.GetEngineTypeVersion("secret")

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
				s.Require().NoError(s.client.EnableKV2Engine(tc.path))
			}

			err := s.client.EnableKV2EngineErrorIfNotForced(tc.force, tc.path)

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
				s.Require().NoError(s.client.EnableKV2Engine(engine), tc.name)
			}

			res, err := s.client.ListAllKVSecretEngines("")
			s.Require().NoError(err, tc.name)

			s.Require().ElementsMatch(tc.expected[""], res[""], tc.name)
		})
	}
}
