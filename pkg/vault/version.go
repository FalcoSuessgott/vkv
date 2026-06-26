package vault

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

// SecretVersion represents a single version of a KVv2 secret.
type SecretVersion struct {
	Version      int        `json:"version"`
	CreatedTime  time.Time  `json:"created_time"`
	DeletionTime *time.Time `json:"deletion_time,omitempty"` // nil means the version has not been deleted
	Destroyed    bool       `json:"destroyed"`
	// Data holds the secrets key-value pairs. It is nil for deleted or
	// destroyed versions, since their data is no longer retrievable.
	Data map[string]interface{} `json:"data,omitempty"`
}

// VersionedSecret holds all versions of a single KVv2 secret plus its custom metadata.
type VersionedSecret struct {
	// CustomMetadata is shared across all versions of the secret (may be nil).
	CustomMetadata map[string]interface{} `json:"custom_metadata,omitempty"`
	// Versions are ordered newest first.
	Versions []*SecretVersion `json:"versions"`
}

// VersionedSecrets maps a secret subPath to its versioned secret.
type VersionedSecrets map[string]*VersionedSecret

// ReadCurrentVersionCreatedTime returns the creation time of a secret's current (latest) version.
func (v *Vault) ReadCurrentVersionCreatedTime(ctx context.Context, rootPath, subPath string) (time.Time, error) {
	data, err := v.Client.Logical().ReadWithContext(ctx, fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath))
	if err != nil {
		return time.Time{}, err
	}

	if data == nil || data.Data == nil {
		return time.Time{}, fmt.Errorf("could not read secret %s metadata", path.Join(rootPath, subPath))
	}

	// prefer the current version's created_time, fall back to the secret's updated_time
	if cur, ok := data.Data["current_version"]; ok {
		if versions, ok := data.Data["versions"].(map[string]interface{}); ok {
			if m, ok := versions[fmt.Sprintf("%v", cur)].(map[string]interface{}); ok {
				if t := parseVaultTime(m["created_time"]); !t.IsZero() {
					return t, nil
				}
			}
		}
	}

	return parseVaultTime(data.Data["updated_time"]), nil
}

// ReadAllVersions returns all versions of a single KVv2 secret, newest first.
func (v *Vault) ReadAllVersions(ctx context.Context, rootPath, subPath string) (*VersionedSecret, error) {
	metadata, err := v.Client.Logical().ReadWithContext(ctx, fmt.Sprintf(kvv2ListSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if metadata == nil || metadata.Data == nil {
		return nil, fmt.Errorf("could not read secret %s metadata", path.Join(rootPath, subPath))
	}

	versions, ok := metadata.Data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not read versions of secret %s", path.Join(rootPath, subPath))
	}

	secret := &VersionedSecret{
		Versions: make([]*SecretVersion, 0, len(versions)),
	}

	if cm, ok := metadata.Data["custom_metadata"].(map[string]interface{}); ok {
		secret.CustomMetadata = cm
	}

	for vStr, meta := range versions {
		versionNr, err := strconv.Atoi(vStr)
		if err != nil {
			return nil, fmt.Errorf("invalid version %q for secret %s: %w", vStr, path.Join(rootPath, subPath), err)
		}

		m, ok := meta.(map[string]interface{})
		if !ok {
			continue
		}

		sv := &SecretVersion{
			Version:     versionNr,
			CreatedTime: parseVaultTime(m["created_time"]),
		}

		if dt := parseVaultTime(m["deletion_time"]); !dt.IsZero() {
			sv.DeletionTime = &dt
		}

		if destroyed, ok := m["destroyed"].(bool); ok {
			sv.Destroyed = destroyed
		}

		// only retrievable versions have data
		if !sv.Destroyed && sv.DeletionTime == nil {
			data, err := v.readSecretVersionData(ctx, rootPath, subPath, versionNr)
			if err != nil {
				return nil, err
			}

			sv.Data = data
		}

		secret.Versions = append(secret.Versions, sv)
	}

	// newest version first
	sort.Slice(secret.Versions, func(i, j int) bool {
		return secret.Versions[i].Version > secret.Versions[j].Version
	})

	return secret, nil
}

// readSecretVersionData reads the key-value data of a specific KVv2 secret version.
func (v *Vault) readSecretVersionData(ctx context.Context, rootPath, subPath string, version int) (map[string]interface{}, error) {
	apiPath := fmt.Sprintf(kvv2ReadWriteSecretsPath, rootPath, subPath)

	data, err := v.Client.Logical().ReadWithDataWithContext(ctx, apiPath, map[string][]string{
		"version": {strconv.Itoa(version)},
	})
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets in %s (version %d) found", path.Join(rootPath, subPath), version)
	}

	if d, ok := data.Data["data"].(map[string]interface{}); ok {
		return d, nil
	}

	return nil, fmt.Errorf("no secrets in %s (version %d) found", path.Join(rootPath, subPath), version)
}

// ListRecursiveAllVersions recursively reads all versions of every KVv2 secret under subPath.
func (v *Vault) ListRecursiveAllVersions(ctx context.Context, rootPath, subPath string, skipErrors bool) (VersionedSecrets, error) {
	isV1, err := v.IsKVv1(ctx, rootPath)
	if err != nil {
		return nil, err
	}

	if isV1 {
		return nil, fmt.Errorf("--all-versions is only supported for KVv2 engines, %q is a KVv1 engine", rootPath)
	}

	result := make(VersionedSecrets)
	if err := v.listRecursiveAllVersions(ctx, rootPath, subPath, skipErrors, result); err != nil {
		return nil, err
	}

	return result, nil
}

// nolint: cyclop
func (v *Vault) listRecursiveAllVersions(ctx context.Context, rootPath, subPath string, skipErrors bool, acc VersionedSecrets) error {
	keys, err := v.ListKeys(ctx, rootPath, subPath)
	if err != nil {
		// no sub directories, treat subPath as a leaf secret
		secret, err := v.ReadAllVersions(ctx, rootPath, subPath)
		if err != nil {
			if skipErrors {
				return nil
			}

			return fmt.Errorf("could not read secret versions from %s: %w.\n\nYou can skip this error using --skip-errors", path.Join(rootPath, subPath), err)
		}

		acc[strings.TrimSuffix(subPath, utils.Delimiter)] = secret

		return nil
	}

	for _, k := range keys {
		nextPath := path.Join(subPath, k)

		if strings.HasSuffix(k, utils.Delimiter) {
			if err := v.listRecursiveAllVersions(ctx, rootPath, nextPath, skipErrors, acc); err != nil {
				return err
			}

			continue
		}

		secret, err := v.ReadAllVersions(ctx, rootPath, nextPath)
		if err != nil {
			if skipErrors {
				continue
			}

			return err
		}

		acc[nextPath] = secret
	}

	return nil
}

// parseVaultTime parses a Vault RFC3339 timestamp, returning the zero time on empty/invalid input.
func parseVaultTime(v interface{}) time.Time {
	s, ok := v.(string)
	if !ok || s == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}

	return t
}
