package snapshot

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"

// 	"github.com/FalcoSuessgott/vkv/pkg/fs"
// 	"github.com/FalcoSuessgott/vkv/pkg/vault"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func (s *VaultSuite) TestSnapshotSaveCommand() {
// 	testCases := []struct {
// 		name        string
// 		restoreArgs []string
// 		saveArgs    []string
// 		expEngines  vault.Engines
// 	}{
// 		{
// 			name:        "restore",
// 			restoreArgs: []string{"--source=testdata/vkv-snapshot-export"},
// 			saveArgs:    []string{"--destination=../vkv-snapshot-export"},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			// restore a snapshot
// 			restoreCmd := newSnapshotRestoreCmd(io.Discard, s.client)
// 			restoreCmd.SetOut(io.Discard)
// 			restoreCmd.SetArgs(tc.restoreArgs)

// 			require.NoError(s.Suite.T(), restoreCmd.Execute())

// 			// save snapshot to tmp dir
// 			saveCmd := newSnapshotSaveCmd(io.Discard, s.client)
// 			saveCmd.SetOut(io.Discard)
// 			saveCmd.SetArgs(tc.saveArgs)

// 			require.NoError(s.Suite.T(), saveCmd.Execute())

// 			err := filepath.Walk("testdata/vkv-snapshot-export", func(p string, info os.FileInfo, err error) error {
// 				if err != nil {
// 					s.Suite.T().Fail()
// 				}

// 				expOut, err := fs.ReadFile(info.Name())
// 				require.NoError(s.Suite.T(), err)

// 				resOut, err := fs.ReadFile(fmt.Sprintf("../vkv-snapshot-export/%s", info.Name()))
// 				require.NoError(s.Suite.T(), err)

// 				assert.Equal(s.Suite.T(), expOut, resOut, info.Name())

// 				return nil
// 			})
// 			if err != nil {
// 				s.Suite.T().Fail()
// 			}
// 		})
// 	}
// }
