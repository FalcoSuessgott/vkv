package secret

import (
	"bytes"
	"testing"
	"time"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintBaseAllVersions(t *testing.T) {
	now := time.Date(2024, 5, 28, 6, 0, 0, 0, time.UTC)
	min5 := now.Add(-5 * time.Minute)
	min6 := now.Add(-6 * time.Minute)

	testCases := []struct {
		name     string
		vs       vault.VersionedSecrets
		rootPath string
		opts     []Option
		output   string
	}{
		{
			name:     "two versions show their own contents, show values",
			rootPath: "secret",
			vs: vault.VersionedSecrets{
				"admin": {
					Versions: []*vault.SecretVersion{
						{Version: 2, CreatedTime: min5, Data: map[string]interface{}{"sub": "password", "w": "w"}},
						{Version: 1, CreatedTime: min6, Data: map[string]interface{}{"sub": "password"}},
					},
				},
			},
			opts: []Option{ShowValues(true)},
			output: `secret/
└── admin
    ├── [Version 2 created 5 minutes ago]
    │   ├── sub=password
    │   └── w=w
    └── [Version 1 created 6 minutes ago]
        └── sub=password
`,
		},
		{
			name:     "masked values with custom metadata",
			rootPath: "secret",
			vs: vault.VersionedSecrets{
				"admin": {
					CustomMetadata: map[string]interface{}{"key": "value"},
					Versions: []*vault.SecretVersion{
						{Version: 1, CreatedTime: min6, Data: map[string]interface{}{"user": "admin"}},
					},
				},
			},
			opts: []Option{ShowValues(false)},
			output: `secret/
└── admin {key=value}
    └── [Version 1 created 6 minutes ago]
        └── user=*****
`,
		},
		{
			name:     "only keys",
			rootPath: "secret",
			vs: vault.VersionedSecrets{
				"demo": {
					Versions: []*vault.SecretVersion{
						{Version: 1, CreatedTime: min6, Data: map[string]interface{}{"foo": "bar", "user": "admin"}},
					},
				},
			},
			opts: []Option{OnlyKeys(true)},
			output: `secret/
└── demo
    └── [Version 1 created 6 minutes ago]
        ├── foo
        └── user
`,
		},
		{
			name:     "deleted and destroyed versions render no data",
			rootPath: "secret",
			vs: vault.VersionedSecrets{
				"admin": {
					Versions: []*vault.SecretVersion{
						{Version: 3, CreatedTime: min5, Destroyed: true},
						{Version: 2, CreatedTime: min5, DeletionTime: &min5},
						{Version: 1, CreatedTime: min6, Data: map[string]interface{}{"user": "admin"}},
					},
				},
			},
			opts: []Option{ShowValues(true)},
			output: `secret/
└── admin
    ├── [Version 3 destroyed 5 minutes ago]
    ├── [Version 2 deleted 5 minutes ago]
    └── [Version 1 created 6 minutes ago]
        └── user=admin
`,
		},
		{
			name:     "nested paths",
			rootPath: "secret",
			vs: vault.VersionedSecrets{
				"sub/demo": {
					Versions: []*vault.SecretVersion{
						{Version: 1, CreatedTime: min6, Data: map[string]interface{}{"user": "admin"}},
					},
				},
			},
			opts: []Option{ShowValues(true)},
			output: `secret/
└── sub
    └── demo
        └── [Version 1 created 6 minutes ago]
            └── user=admin
`,
		},
	}

	for _, tc := range testCases {
		var b bytes.Buffer

		tc.opts = append(tc.opts,
			WithWriter(&b),
			WithEnginePath(utils.NormalizePath(tc.rootPath)),
		)

		p := NewSecretPrinter(tc.opts...)
		p.now = now

		require.NoError(t, p.Out(tc.vs), tc.name)
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}

func TestPrintAllVersionsJSONYAML(t *testing.T) {
	t1 := time.Date(2024, 5, 28, 4, 47, 54, 0, time.UTC)
	deleted := time.Date(2024, 5, 28, 5, 0, 0, 0, time.UTC)

	vs := vault.VersionedSecrets{
		"admin": {
			CustomMetadata: map[string]interface{}{"key": "value"},
			Versions: []*vault.SecretVersion{
				{Version: 2, CreatedTime: t1, DeletionTime: &deleted},
				{Version: 1, CreatedTime: t1, Data: map[string]interface{}{"user": "admin"}},
			},
		},
	}

	t.Run("json", func(t *testing.T) {
		var b bytes.Buffer

		p := NewSecretPrinter(ToFormat(JSON), WithWriter(&b), WithEnginePath("secret/"))
		require.NoError(t, p.Out(vs))

		assert.Equal(t, `{
  "admin": {
    "custom_metadata": {
      "key": "value"
    },
    "versions": [
      {
        "version": 2,
        "created_time": "2024-05-28T04:47:54Z",
        "deletion_time": "2024-05-28T05:00:00Z",
        "destroyed": false
      },
      {
        "version": 1,
        "created_time": "2024-05-28T04:47:54Z",
        "destroyed": false,
        "data": {
          "user": "admin"
        }
      }
    ]
  }
}
`, b.String())
	})

	t.Run("yaml", func(t *testing.T) {
		var b bytes.Buffer

		p := NewSecretPrinter(ToFormat(YAML), WithWriter(&b), WithEnginePath("secret/"))
		require.NoError(t, p.Out(vs))

		assert.Equal(t, `admin:
  custom_metadata:
    key: value
  versions:
  - created_time: "2024-05-28T04:47:54Z"
    deletion_time: "2024-05-28T05:00:00Z"
    destroyed: false
    version: 2
  - created_time: "2024-05-28T04:47:54Z"
    data:
      user: admin
    destroyed: false
    version: 1
`, b.String())
	})
}
