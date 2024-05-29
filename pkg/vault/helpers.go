package vault

import (
	"errors"
	"fmt"
	"path"
)

func (v *Vault) ReadSecrets(rootPath, subPath string, version ...int) (*Secret, error) {
	// error if more than 1 version specified
	if len(version) > 1 {
		return nil, fmt.Errorf("multiple versions specified")
	}

	v1, err := v.IsKVv1(rootPath)
	if err != nil {
		return nil, err
	}

	// return kv1 secret
	if v1 {
		secret, err := v.Client.KVv1(rootPath).Get(v.Context, subPath)
		if err != nil {
			return nil, err
		}

		return &Secret{
			Data: secret.Data,
		}, nil
	}

	// if version specified, return specific secret version
	if len(version) == 1 {
		secret, err := v.Client.KVv2(rootPath).GetVersion(v.Context, subPath, version[0])
		if err != nil {
			return nil, err
		}

		return &Secret{
			Data:               secret.Data,
			CustomMetadata:     secret.CustomMetadata,
			Version:            secret.VersionMetadata.Version,
			VersionCreatedTime: secret.VersionMetadata.CreatedTime,
			Destroyed:          secret.VersionMetadata.Destroyed,
			DeletionTime:       secret.VersionMetadata.DeletionTime,
		}, nil
	}

	// return latest version
	secret, err := v.Client.KVv2(rootPath).Get(v.Context, subPath)
	if err != nil {
		return nil, err
	}

	return &Secret{
		Data:               secret.Data,
		CustomMetadata:     secret.CustomMetadata,
		Version:            secret.VersionMetadata.Version,
		VersionCreatedTime: secret.VersionMetadata.CreatedTime,
		Destroyed:          secret.VersionMetadata.Destroyed,
		DeletionTime:       secret.VersionMetadata.DeletionTime,
	}, nil
}

// ListKeys returns all keys from vault kv secret path.
func (v *Vault) ListKeys(rootPath, subPath string) ([]string, error) {
	apiPath := fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath)

	isV1, err := v.IsKVv1(rootPath)
	if err != nil {
		return nil, err
	}

	if isV1 {
		apiPath = fmt.Sprintf(kvv1ListSecretsPath, rootPath, subPath)
	}

	data, err := v.Client.Logical().ListWithContext(v.Context, apiPath)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no keys found in \"%s\"", path.Join(rootPath, subPath))
	}

	keys := []string{}

	k, ok := data.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response \"%v\"", data.Data["keys"])
	}

	for _, e := range k {
		keys = append(keys, fmt.Sprintf("%v", e))
	}

	return keys, nil
}

// IsKVv1 returns true if the current path is a KVv1 Engine.
func (v *Vault) IsKVv1(path string) (bool, error) {
	data, err := v.Client.Logical().ReadWithContext(v.Context, fmt.Sprintf(mountDetailsPath, path))
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, errors.New("cannot lookup mount type")
	}

	// early versions of Vaults KV engine are of type "generic"
	if data.Data["type"] == "generic" {
		return true, nil
	}

	if opts, ok := data.Data["options"].(map[string]interface{}); ok {
		//nolint: forcetypeassert
		if opts["version"].(string) == "1" {
			return true, nil
		}
	}

	return false, nil
}

// WriteSecrets writes kv secrets to a specified path.
func (v *Vault) WriteSecrets(rootPath, subPath string, secrets map[string]interface{}) error {
	apiPath := fmt.Sprintf(kvv2ReadWriteSecretsPath, rootPath, subPath)
	options := map[string]interface{}{}

	isV1, err := v.IsKVv1(rootPath)
	if err != nil {
		return err
	}

	if isV1 {
		apiPath = fmt.Sprintf(kvv1ReadWriteSecretsPath, rootPath, subPath)
		options = secrets
	} else {
		options["data"] = secrets
	}

	_, err = v.Client.Logical().WriteWithContext(v.Context, apiPath, options)
	if err != nil {
		return err
	}

	return nil
}

// DisableKV2Engine disables the kv2 engine at a specified path.
func (v *Vault) DisableKV2Engine(rootPath string) error {
	_, err := v.Client.Logical().DeleteWithContext(v.Context, fmt.Sprintf(mountEnginePath, rootPath))
	if err != nil {
		return err
	}

	return nil
}

// GetAllVersions returns the number of versions for a kv2 secret, returns 0 if no KVv2 engine.
func (v *Vault) GetAllVersions(rootPath, subPath string) (int, error) {
	v1, err := v.IsKVv1(rootPath)
	if err != nil {
		return 0, err
	}

	if v1 {
		return 0, nil
	}

	versions, err := v.Client.KVv2(rootPath).GetVersionsAsList(v.Context, subPath)
	if err != nil {
		return 0, fmt.Errorf("cannot list versions for %s/%s: %w", rootPath, subPath, err)
	}

	return len(versions), nil
}
