package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
)

func (s *VaultSuite) TestSnapshotSaveCommand() {
	testCases := []struct {
		name        string
		restoreArgs []string
		saveArgs    []string
		expEngines  vault.Engines
	}{
		{
			name:        "restore",
			restoreArgs: []string{"--source=testdata/vkv-snapshot-export"},
			saveArgs:    []string{"--destination=vkv-snapshot-export-test"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			writer = io.Discard

			// restore a snapshot
			restoreCmd := NewSnapshotRestoreCmd()
			restoreCmd.SetArgs(tc.restoreArgs)

			s.Require().NoError(restoreCmd.Execute(), "restore snapshot "+tc.name)

			// save snapshot to tmp dir
			saveCmd := NewSnapshotSaveCmd()
			saveCmd.SetArgs(tc.saveArgs)

			s.Require().NoError(saveCmd.Execute(), "save snapshot "+tc.name)

			err := filepath.Walk("testdata/vkv-snapshot-export", func(p string, info os.FileInfo, err error) error {
				s.Require().NoError(err, "filepath walk failed "+tc.name)

				if info.Name() == "vkv-snapshot-export" {
					return nil
				}

				// snapshot result file
				expOut, err := fs.ReadFile("testdata/vkv-snapshot-export/" + info.Name())
				s.Require().NoError(err, "error reading snapshot file "+info.Name())

				resOut, err := fs.ReadFile("vkv-snapshot-export-test/" + info.Name())
				s.Require().NoError(err, "error reading created snapshot file "+info.Name())

				s.Require().Equal(expOut, resOut, info.Name())

				return nil
			})

			s.Require().NoError(err, "test failed "+tc.name)
		})

		s.T().Cleanup(func() {
			os.RemoveAll("vkv-snapshot-export-test")
		})
	}
}
