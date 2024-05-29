package vault

import (
	"fmt"
	"path"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

// NewKVSecrets returns all KV secrets for a given kv mount.
func NewKVSecrets(v *Vault, rootPath, subPath string, skipErrors bool, allVersions bool) (*KVSecrets, error) {
	kv := &KVSecrets{
		Vault:     v,
		MountPath: rootPath,
		Secrets:   make(map[string][]*Secret),
	}

	desc, err := v.GetEngineDescription(rootPath)
	if err != nil {
		return nil, err
	}

	kv.Description = desc

	engineType, version, err := v.GetEngineTypeVersion(rootPath)
	if err != nil {
		return nil, err
	}

	kv.Type = engineType + version

	if err := kv.iterator(subPath, skipErrors, allVersions); err != nil {
		return nil, err
	}

	return kv, nil
}

func (kv *KVSecrets) iterator(subPath string, skipErrors bool, allVersions bool) error {
	// list keys for the current secret dir
	keys, err := kv.ListKeys(kv.MountPath, subPath)
	// no sub directories in here, but lets check for normal kv pairs then..
	if err != nil || len(keys) == 0 {
		if err := kv.listSecrets(subPath, skipErrors, allVersions); err != nil {
			return err
		}

		return nil
	}

	// we found keys, lets add them to the list or dig deeper
	for _, k := range keys {
		// / at the end means the secret is a dir, so we go into it ...
		if strings.HasSuffix(k, utils.Delimiter) {
			if err := kv.iterator(path.Join(subPath, k), skipErrors, allVersions); err != nil {
				return err
			}
		} else {
			if err := kv.listSecrets(path.Join(subPath, k), skipErrors, allVersions); err != nil {
				return err
			}
		}
	}

	return nil
}

func (kv *KVSecrets) listSecrets(p string, skipErrors bool, allVersions bool) error {
	versions, err := kv.GetAllVersions(kv.MountPath, p)
	if err != nil {
		return err
	}

	if !allVersions || versions == 0 {
		secrets, err := kv.ReadSecrets(kv.MountPath, p)
		if !skipErrors && err != nil {
			return fmt.Errorf("could not read secrets from %s/%s: %w.\n\nYou can skip this error using --skip-errors", kv.MountPath, p, err)
		}

		kv.Secrets[path.Join(kv.MountPath, p)] = append(kv.Secrets[path.Join(kv.MountPath, p)], secrets)
		return nil
	}

	for i := 1; i <= versions; i++ {
		secrets, err := kv.ReadSecrets(kv.MountPath, p, i)
		if !skipErrors && err != nil {
			return fmt.Errorf("could not read secrets from %s/%s: %w.\n\nYou can skip this error using --skip-errors", kv.MountPath, p, err)
		}

		kv.Secrets[path.Join(kv.MountPath, p)] = append(kv.Secrets[path.Join(kv.MountPath, p)], secrets)
	}

	return nil
}
