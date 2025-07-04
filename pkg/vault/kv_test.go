package vault

import (
	"context"
	"path"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestListRecursive() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		err      bool
		v1       bool
		secrets  Secrets
		expected Secrets
	}{
		{
			name:     "simple secret",
			rootPath: "kvv2",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user": "password",
				},
				"sub2": map[string]interface{}{
					"user": false,
				},
			},
			expected: map[string]interface{}{
				"kvv2": Secrets{
					"sub": map[string]interface{}{
						"user": "password",
					},
					"sub2": map[string]interface{}{
						"user": false,
					},
				},
			},
		},
		{
			name:     "simple secret",
			rootPath: "kvv1",
			v1:       true,
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user": "password",
				},
				"sub2": map[string]interface{}{
					"user": false,
				},
			},
			expected: map[string]interface{}{
				"kvv1": Secrets{
					"sub": map[string]interface{}{
						"user": "password",
					},
					"sub2": map[string]interface{}{
						"user": false,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// write secrets
			if tc.v1 {
				require.NoError(s.T(), s.client.EnableKV1Engine(context.Background(), tc.rootPath))
			} else {
				require.NoError(s.T(), s.client.EnableKV2Engine(context.Background(), tc.rootPath))
			}

			for k, secrets := range tc.secrets {
				if m, ok := secrets.(map[string]interface{}); ok {
					require.NoError(s.T(), s.client.WriteSecrets(context.Background(), tc.rootPath, path.Join(tc.subPath, k), m))
				}
			}

			// read secrets
			res := make(Secrets)
			secrets, err := s.client.ListRecursive(context.Background(), tc.rootPath, tc.subPath, false)
			require.NoError(s.T(), err)

			res[tc.rootPath] = *secrets

			// assert
			if tc.err {
				require.Error(s.T(), err)
			} else {
				require.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expected, res, tc.name)
			}

			require.NoError(s.T(), s.client.DisableKV2Engine(context.Background(), tc.rootPath))
		})
	}
}

func (s *VaultSuite) TestReadSecretMetadataVersion() {
	testCases := []struct {
		name     string
		rootPath string
		subPath  string
		secrets  Secrets
		version  string
		metadata interface{}
	}{
		{
			name:     "simple secret",
			rootPath: "kv",
			subPath:  "subpath",
			secrets: map[string]interface{}{
				"sub": map[string]interface{}{
					"user": "password",
				},
				"sub2": map[string]interface{}{
					"user": false,
				},
			},
			metadata: nil,
			version:  "1",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// write secrets
			require.NoError(s.T(), s.client.EnableKV2Engine(context.Background(), tc.rootPath))

			for k, v := range tc.secrets {
				if m, ok := v.(map[string]interface{}); ok {
					require.NoError(s.T(), s.client.WriteSecrets(context.Background(), tc.rootPath, path.Join(tc.subPath, k), m), tc.name)
				}
			}

			// read metadata
			for k := range tc.secrets {
				md, err := s.client.ReadSecretMetadata(context.Background(), tc.rootPath, path.Join(tc.subPath, k))
				require.NoError(s.T(), err, tc.name)
				require.EqualValues(s.T(), tc.metadata, md, "we currently cant write metadata")

				v, err := s.client.ReadSecretVersion(context.Background(), tc.rootPath, path.Join(tc.subPath, k))
				require.NoError(s.T(), err, tc.name)

				// assert
				require.EqualValues(s.T(), tc.version, v, "version")
			}

			require.NoError(s.T(), s.client.DisableKV2Engine(context.Background(), tc.rootPath))
		})
	}
}
