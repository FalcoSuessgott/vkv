package vault

// import (
// 	"fmt"
// )

// func (s *VaultSuite) TestSanitizerOnlyPaths() {
// 	testCases := []struct {
// 		name string
// 		kv   *KVSecrets
// 		exp  string
// 		err  bool
// 	}{
// 		{
// 			name: "default",
// 			kv: &KVSecrets{
// 				Secrets: map[string][]*Secret{
// 					"test/admin": {
// 						{
// 							Data: map[string]interface{}{},
// 						},
// 					},
// 				},
// 			},
// 			exp: "",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			// exec sanitizerFunc
// 			err := tc.kv.OnlyPaths()()

// 			if tc.err {
// 				s.Require().Error(err, tc.name)
// 			} else {
// 				s.Require().NoError(err, tc.name)
// 				s.Require().Equal(tc.exp, tc.kv.Secrets["test/admin"][0].String(), tc.name)
// 			}
// 		})
// 	}
// }

// func (s *VaultSuite) TestSanitizerOnlyKeys() {
// 	testCases := []struct {
// 		name string
// 		kv   *KVSecrets
// 		exp  string
// 		err  bool
// 	}{
// 		{
// 			name: "default",
// 			kv: &KVSecrets{
// 				Secrets: map[string][]*Secret{
// 					"test/admin": {
// 						{Data: map[string]interface{}{
// 							"foo":  "bar",
// 							"test": "thisisaverylongsecret",
// 						}},
// 					},
// 				},
// 			},
// 			exp: "foo\ntest\n",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			// exec sanitizerFunc
// 			err := tc.kv.OnlyKeys()()

// 			if tc.err {
// 				s.Require().Error(err, tc.name)
// 			} else {
// 				s.Require().NoError(err, tc.name)
// 				s.Require().Equal(tc.exp, tc.kv.Secrets["test/admin"][0].String(), tc.name)
// 			}
// 		})
// 	}
// }

// func (s *VaultSuite) TestSanitizerMaskSecrets() {
// 	testCases := []struct {
// 		name   string
// 		length int
// 		kv     *KVSecrets
// 		exp    string
// 		err    bool
// 	}{
// 		{
// 			name:   "default",
// 			length: 12,
// 			kv: &KVSecrets{
// 				Secrets: map[string][]*Secret{
// 					"test/admin": {
// 						{Data: map[string]interface{}{
// 							"foo":  "bar",
// 							"test": "thisisaverylongsecret",
// 						}},
// 					},
// 				},
// 			},
// 			exp: "foo\t= \"***\"\ntest\t= \"************\"\n",
// 		},
// 		{
// 			name:   "disabled",
// 			length: -1,
// 			kv: &KVSecrets{
// 				Secrets: map[string][]*Secret{
// 					"test/admin": {
// 						{Data: map[string]interface{}{
// 							"foo":  "bar",
// 							"test": "thisisaverylongsecret",
// 						}},
// 					},
// 				},
// 			},
// 			exp: "foo\t= \"***\"\ntest\t= \"*********************\"\n",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			// exec sanitizerFunc
// 			err := tc.kv.MaskSecrets(tc.length)()

// 			if tc.err {
// 				s.Require().Error(err, tc.name)
// 			} else {
// 				s.Require().NoError(err, tc.name)
// 				s.Require().Equal(tc.exp, tc.kv.Secrets["test/admin"][0].String(), tc.name)
// 			}
// 		})
// 	}
// }

// func (s *VaultSuite) TestSanitizerShowDiff() {
// 	testCases := []struct {
// 		name string
// 		kv   *KVSecrets
// 		exp  []string
// 		err  bool
// 	}{
// 		{
// 			name: "default",
// 			kv: &KVSecrets{
// 				MountPath:   "secret",
// 				Description: "test",
// 				Type:        "kv2",
// 				Secrets: map[string][]*Secret{
// 					"test/admin": {
// 						// Version 0 (empty secret) is added by default
// 						// version 1
// 						{Data: map[string]interface{}{
// 							"foo": "bar",
// 						}},
// 						// version 2 (added)
// 						{Data: map[string]interface{}{
// 							"foo": "bar",
// 							"new": "element",
// 						}},
// 						// version 3 (changed)
// 						{Data: map[string]interface{}{
// 							"foo": "change",
// 							"new": "element",
// 						}},
// 						// version 4 (removed)
// 						{Data: map[string]interface{}{}},
// 					},
// 				},
// 			},
// 			exp: []string{
// 				// v1
// 				"[+] foo\t= \"bar\"\n",
// 				// v2
// 				"foo\t= \"bar\"\n[+] new\t= \"element\"\n",
// 				// v3
// 				"[~] foo\t= \"bar\" -> \"change\"\nnew\t= \"element\"\n",
// 				// v4
// 				"[-] foo\t= \"change\"\n[-] new\t= \"element\"\n",
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			// inject test client
// 			tc.kv.Vault = s.client

// 			// disable colored output for test purposes
// 			s.Suite.T().Setenv("NO_COLOR", "true")

// 			// write secrets
// 			for path, secret := range tc.kv.Secrets {
// 				for _, version := range secret {
// 					s.Require().NoError(tc.kv.WriteSecrets(tc.kv.MountPath, path, version.Data))
// 				}
// 			}

// 			// exec sanitizerFunc
// 			err := tc.kv.ShowDiff()()

// 			if tc.err {
// 				s.Require().Error(err, tc.name)
// 			} else {
// 				s.Require().NoError(err, tc.name)
// 				for path, secret := range tc.kv.Secrets {
// 					for i, version := range secret {
// 						s.Require().Equal(tc.exp[i], version.DiffString(), fmt.Sprintf("%s %s@v%d", tc.name, path, i+1))
// 					}
// 				}
// 			}
// 		})
// 	}
// }
