package cmd

import (
	"bytes"
	"io"
)

func (s *VaultSuite) TestVersion() {
	testCases := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "v1.0.0",
			version:  "v1.0.0",
			expected: "vkv v1.0.0\n",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// inject writer and version
			Version = tc.version
			b := bytes.NewBufferString("")
			writer = b

			// run command
			s.Require().NoError(NewVersionCmd().Execute(), tc.name)

			// assert output
			out, _ := io.ReadAll(b)
			s.Require().Equal(tc.expected, string(out), tc.name)
		})
	}
}
