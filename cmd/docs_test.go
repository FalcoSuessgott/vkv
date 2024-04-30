package cmd

import "os"

func (s *VaultSuite) TestDocGen() {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "docgen",
			path: "./docs-test",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmdDocPath = tc.path

			s.Require().NoError(NewDocCmd().Execute(), tc.name)
		})

		s.T().Cleanup(func() {
			os.RemoveAll(tc.path)
		})
	}
}
