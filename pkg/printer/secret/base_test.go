package secret

import (
	"bytes"
	"testing"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintBase(t *testing.T) {
	testCases := []struct {
		name     string
		s        map[string]interface{}
		rootPath string
		opts     []Option
		output   string
		err      bool
	}{
		{
			name:     "test: default options",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(false),
			},
			output: `root/
└── secret
    ├── key=*****
    └── user=********
`,
		},
		{
			name:     "test: show secrets",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(true),
			},
			output: `root/
└── secret
    ├── key=value
    └── user=password
`,
		},
		{
			name:     "test: only paths",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyPaths(true),
			},
			output: `root/
└── secret
`,
		},
		{
			name:     "test: only keys",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "password",
				},
			},
			opts: []Option{
				ToFormat(Base),
				OnlyKeys(true),
			},
			output: `root/
└── secret
    ├── key
    └── user
`,
		},
		{
			name:     "test: multiple lines",
			rootPath: "root",
			s: map[string]interface{}{
				"secret": map[string]interface{}{
					"key":  "value",
					"user": "value",
				},
			},
			opts: []Option{
				ToFormat(Base),
				ShowValues(true),
			},
			output: `root/
└── secret
    ├── key=value
    └── user=value
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

		m := map[string]interface{}{
			utils.NormalizePath(tc.rootPath): tc.s,
		}

		require.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}
