package cmd

func (s *VaultSuite) TestManPages() {
	testCases := []struct {
		name string
	}{
		{
			name: "mapage gen",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().NoError(NewManCmd().Execute(), tc.name)
		})
	}
}
