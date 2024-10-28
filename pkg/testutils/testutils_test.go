package testutils

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/suite"
)

type VaultSuite struct {
	suite.Suite

	c      *TestContainer
	client *vault.Vault
}

func (s *VaultSuite) TearDownSubTest() {
	if err := s.c.Terminate(); err != nil {
		log.Fatal(err)
	}
}

func (s *VaultSuite) SetupSubTest() {
	vc, err := StartTestContainer()
	if err != nil {
		log.Fatal(err)
	}

	s.c = vc

	v, err := vault.NewClient(vc.URI, vc.Token)
	if err != nil {
		log.Fatal(err)
	}

	s.client = v
}

func (s *VaultSuite) TestVaultConnection() {
	s.Run("test", func() {
		health, err := s.client.Client.Sys().Health()
		s.Require().NoError(err)

		fmt.Println(health)
		s.Require().True(health.Initialized, "initialized")
		s.Require().False(health.Sealed, "unsealed")
	})
}

func TestVaultSuite(t *testing.T) {
	// github actions doesn't offer the docker socket, which we need to run this test suite
	if runtime.GOOS != "windows" {
		suite.Run(t, new(VaultSuite))
	}
}
