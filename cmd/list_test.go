package cmd

func (s *VaultSuite) TestListCommand() {
	s.Require().NoError(NewListCmd().Execute(), "list cmd")
}
