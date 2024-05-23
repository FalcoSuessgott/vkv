package vault

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/FalcoSuessgott/vkv/pkg/printer"
// 	"github.com/r3labs/diff/v3"
// )

// // ShowDiff is a santizierFunc that computes a changelog for each secret version with its previous version.
// func (kv *KVSecrets) ShowDiff() printer.SanitizerFunc {
// 	return func() error {
// 		for path, secrets := range kv.Secrets {
// 			// lets prepend an empty secret version as the first secret
// 			secretVersions := []*Secret{{}}
// 			secretVersions = append(secretVersions, secrets...)

// 			for i := range secretVersions {
// 				if i+1 < len(secretVersions) {
// 					log, err := diff.Diff(secretVersions[i].Data, secretVersions[i+1].Data)
// 					if err != nil {
// 						return err
// 					}

// 					kv.Secrets[path][i].Changelog = log
// 				}
// 			}
// 		}

// 		return nil
// 	}
// }

// func (kv *KVSecrets) OnlyKeys() printer.SanitizerFunc {
// 	return func() error {
// 		for _, secrets := range kv.Secrets {
// 			for _, s := range secrets {
// 				for k := range s.Data {
// 					s.Data[k] = ""
// 				}
// 			}
// 		}

// 		return nil
// 	}
// }

// func (kv *KVSecrets) OnlyPaths() printer.SanitizerFunc {
// 	return func() error {
// 		for _, secrets := range kv.Secrets {
// 			for _, s := range secrets {
// 				s.Data = map[string]interface{}{}
// 			}
// 		}

// 		return nil
// 	}
// }
