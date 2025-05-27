package cmd

import (
	"io"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
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
				"": []string{"secret/", "secret_2/"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			writer = io.Discard

			cmd := NewSnapshotRestoreCmd()
			cmd.SetArgs(tc.args)

			s.Require().NoError(cmd.Execute())

			engines, err := vaultClient.ListAllKVSecretEngines(rootContext, "")
			s.Require().NoError(err)

			for expNS, expEngines := range tc.expEngines {
				resEngines := engines[expNS]

				s.Require().ElementsMatch(expEngines, resEngines, tc.name)

				for _, engine := range expEngines {
					if engine == "secret/" {
						continue
					}

					secret, err := vaultClient.ListRecursive(rootContext, path.Join(expNS, engine), "", false)
					s.Require().NoError(err)

					out, err := fs.ReadFile(path.Join("testdata/vkv-snapshot-export", expNS, strings.TrimSuffix(engine, "/")+".yaml"))
					s.Require().NoError(err)

					res, _ := utils.FromJSON(out)

					s.Require().Equal(res, utils.ToMapStringInterface(secret), tc.name)
				}
			}
		})
	}
}
