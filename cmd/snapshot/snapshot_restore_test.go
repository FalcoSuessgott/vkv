package snapshot

import (
	"io"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestSnapShotRestoreCommand() {
	testCases := []struct {
		name       string
		args       []string
		expEngines vault.Engines
	}{
		{
			name: "restore",
			args: []string{"--source=testdata/vkv-snapshot-export"},
			expEngines: vault.Engines{
				"":                 []string{"secret/", "secret_2/"},
				"sub":              []string{"sub_secret/", "sub_secret_2/"},
				"sub/sub2":         []string{"sub_sub2_secret/", "sub_sub2_secret_2/"},
				"test":             []string{},
				"test/test2":       []string{},
				"test/test2/test3": []string{"test_test2_test3_secret/", "test_test2_test3_secret_2/"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := newSnapshotRestoreCmd(io.Discard, s.client)
			cmd.SetOut(io.Discard)
			cmd.SetArgs(tc.args)

			require.NoError(s.Suite.T(), cmd.Execute())

			engines, err := s.client.ListAllKVSecretEngines("")
			require.NoError(s.Suite.T(), err)

			for expNS, expEngines := range tc.expEngines {
				resEngines := engines[expNS]

				assert.ElementsMatch(s.Suite.T(), expEngines, resEngines, tc.name)

				for _, engine := range expEngines {
					if engine == "secret/" {
						continue
					}

					secret, err := s.client.ListRecursive(path.Join(expNS, engine), "")
					require.NoError(s.Suite.T(), err)

					out, err := fs.ReadFile(path.Join("testdata/vkv-snapshot-export", expNS, strings.TrimSuffix(engine, "/")+".yaml"))
					require.NoError(s.Suite.T(), err)

					res, _ := utils.FromJSON(out)

					assert.Equal(s.Suite.T(), res, utils.ToMapStringInterface(secret), tc.name)
				}
			}
		})
	}
}
