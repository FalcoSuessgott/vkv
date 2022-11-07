package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
				WithVaultClient(nil),
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
	}

	for _, tc := range testCases {
		var b bytes.Buffer
		tc.opts = append(tc.opts, WithWriter(&b))

		p := NewPrinter(tc.opts...)

		m := map[string]interface{}{}

		m[tc.rootPath+"/"] = tc.s
		assert.NoError(t, p.Out(m))
		assert.Equal(t, tc.output, b.String(), tc.name)
	}
}
