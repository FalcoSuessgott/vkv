package vault

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *VaultSuite) TestReadAllVersions() {
	s.Run("read all versions", func() {
		ctx := context.Background()
		rootPath := "kvv2"
		subPath := "admin"

		require.NoError(s.T(), s.client.EnableKV2Engine(ctx, rootPath))

		// write two versions of the same secret
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, subPath, map[string]interface{}{"user": "v1"}))
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, subPath, map[string]interface{}{"user": "v2", "extra": "x"}))

		secret, err := s.client.ReadAllVersions(ctx, rootPath, subPath)
		require.NoError(s.T(), err)
		require.Len(s.T(), secret.Versions, 2)

		// newest version first
		assert.Equal(s.T(), 2, secret.Versions[0].Version)
		assert.Equal(s.T(), 1, secret.Versions[1].Version)

		assert.Equal(s.T(), map[string]interface{}{"user": "v2", "extra": "x"}, secret.Versions[0].Data)
		assert.Equal(s.T(), map[string]interface{}{"user": "v1"}, secret.Versions[1].Data)

		// timestamps are populated
		assert.False(s.T(), secret.Versions[0].CreatedTime.IsZero())
		assert.False(s.T(), secret.Versions[1].CreatedTime.IsZero())
	})
}

func (s *VaultSuite) TestReadAllVersionsDeleted() {
	s.Run("read all versions with a deleted version", func() {
		ctx := context.Background()
		rootPath := "kvv2"
		subPath := "admin"

		require.NoError(s.T(), s.client.EnableKV2Engine(ctx, rootPath))

		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, subPath, map[string]interface{}{"user": "v1"}))
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, subPath, map[string]interface{}{"user": "v2"}))

		// soft-delete version 1 (metadata remains, data is gone)
		require.NoError(s.T(), s.deleteVersion(ctx, rootPath, subPath, 1))

		secret, err := s.client.ReadAllVersions(ctx, rootPath, subPath)
		require.NoError(s.T(), err)
		require.Len(s.T(), secret.Versions, 2)

		// version 1 is deleted: deletion time set, no data
		deleted := secret.Versions[1]
		assert.Equal(s.T(), 1, deleted.Version)
		assert.NotNil(s.T(), deleted.DeletionTime)
		assert.Nil(s.T(), deleted.Data)
	})
}

func (s *VaultSuite) TestListRecursiveAllVersions() {
	s.Run("list recursive all versions", func() {
		ctx := context.Background()
		rootPath := "kvv2"

		require.NoError(s.T(), s.client.EnableKV2Engine(ctx, rootPath))

		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, "admin", map[string]interface{}{"user": "v1"}))
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, "admin", map[string]interface{}{"user": "v2"}))
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, "sub/demo", map[string]interface{}{"foo": "bar"}))

		vs, err := s.client.ListRecursiveAllVersions(ctx, rootPath, "", false)
		require.NoError(s.T(), err)

		require.Contains(s.T(), vs, "admin")
		require.Contains(s.T(), vs, "sub/demo")
		assert.Len(s.T(), vs["admin"].Versions, 2)
		assert.Len(s.T(), vs["sub/demo"].Versions, 1)
	})
}

func (s *VaultSuite) TestListRecursiveAllVersionsKVv1() {
	s.Run("list recursive all versions on a KVv1 engine errors", func() {
		ctx := context.Background()
		rootPath := "kvv1"

		require.NoError(s.T(), s.client.EnableKV1Engine(ctx, rootPath))
		require.NoError(s.T(), s.client.WriteSecrets(ctx, rootPath, "admin", map[string]interface{}{"user": "v1"}))

		_, err := s.client.ListRecursiveAllVersions(ctx, rootPath, "", false)
		require.Error(s.T(), err, "--all-versions should fail on a KVv1 engine")
	})
}

// deleteVersion soft-deletes a single KVv2 secret version.
func (s *VaultSuite) deleteVersion(ctx context.Context, rootPath, subPath string, version int) error {
	_, err := s.client.Client.Logical().WriteWithContext(ctx,
		rootPath+"/delete/"+subPath,
		map[string]interface{}{"versions": []int{version}},
	)

	return err
}
