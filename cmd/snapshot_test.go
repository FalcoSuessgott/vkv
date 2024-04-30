package cmd

func (s *VaultSuite) TestSnapshotCommand() {
	s.Require().NoError(NewSnapshotCmd().Execute(), "snapshot cmd")
}
