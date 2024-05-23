package vault

import (
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/r3labs/diff/v3"
)

func (s *VaultSuite) TestString() {
	testCases := []struct {
		name   string
		secret *Secret
		exp    string
	}{
		{
			name: "simple",
			secret: &Secret{
				Data: map[string]interface{}{
					"this": "one",
					"key":  "value",
					"foo":  "12",
					"bar":  "false",
				},
			},
			exp: "bar\t= \"false\"\nfoo\t= \"12\"\nkey\t= \"value\"\nthis\t= \"one\"\n",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(tc.exp, tc.secret.String(false, -1), tc.name)
		})
	}
}

func (s *VaultSuite) TestDiffString() {
	testCases := []struct {
		name          string
		previous      *Secret
		currentSecret *Secret
		exp           string
	}{
		{
			name: "equal",
			previous: &Secret{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			currentSecret: &Secret{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			exp: "key\t= \"value\"\n",
		},
		{
			name: "added",
			previous: &Secret{
				Data: map[string]interface{}{},
			},
			currentSecret: &Secret{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			exp: "[+] key\t= \"value\"\n",
		},
		{
			name: "changed",
			previous: &Secret{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			currentSecret: &Secret{
				Data: map[string]interface{}{
					"key": "changed",
				},
			},
			exp: "[~] key\t= \"value\" -> \"changed\"\n",
		},
		{
			name: "deleted",
			previous: &Secret{
				Data: map[string]interface{}{
					"key": "value",
				},
			},
			currentSecret: &Secret{
				Data: map[string]interface{}{},
			},
			exp: "[-] key\t= \"value\"\n",
		},
		{
			name: "complex",
			previous: &Secret{
				Data: map[string]interface{}{
					"this": "one",
					"key":  "value",
					"foo":  "12",
					"bar":  "false",
				},
			},
			currentSecret: &Secret{
				Data: map[string]interface{}{
					"this":    "one",
					"key":     "changed",
					"another": "one",
					"bar":     "false",
				},
			},
			exp: `[+] another	= "one"
bar	= "false"
[-] foo	= "12"
[~] key	= "value" -> "changed"
this	= "one"
`,
		},
	}

	for _, tc := range testCases {
		// disable colored output for test purposes
		s.Suite.T().Setenv(utils.NoColorEnv, "true")

		log, err := diff.Diff(tc.previous.Data, tc.currentSecret.Data)
		if err != nil {
			s.Require().NoError(err, tc.name)
		}

		tc.currentSecret.Changelog = log

		s.Run(tc.name, func() {
			s.Require().Equal(tc.exp, tc.currentSecret.DiffString(false, false, -1), tc.name)
		})
	}
}
