package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/r3labs/diff/v3"
)

// nolint: gosec
const (
	kvv1ReadWriteSecretsPath = "%s/%s"
	kvv1ListSecretsPath      = "%s/%s"

	kvv2ReadWriteSecretsPath = "%s/data/%s"
	kvv2ListSecretsPath      = "%s/metadata/%s"

	mountDetailsPath  = "sys/internal/ui/mounts/%s"
	mountEnginePath   = "sys/mounts/%s"
	listSecretEngines = "sys/mounts"

	capabilities = "sys/capabilities-self"

	listNamespaces  = "sys/namespaces"
	createNamespace = "sys/namespaces/%s"

	defaultTimestamp = "00010101000000"
	dateFormat       = "Monday, 02-Jan-06 15:04:05"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client

	Context context.Context
}

// KVSecrets struct for kv secrets.
type KVSecrets struct {
	*Vault `json:"-"`

	MountPath   string               `json:"mount_path"`
	Type        string               `json:"type"`
	Description string               `json:"description"`
	Secrets     map[string][]*Secret `json:"secrets"`
}

// Secret is a single KV secret
type Secret struct {
	Data               map[string]interface{} `json:"data"`
	Changelog          diff.Changelog         `json:"-"`
	CustomMetadata     map[string]interface{} `json:"custom_metadata"`
	Version            int                    `json:"version"`
	VersionCreatedTime time.Time              `json:"version_created_time"`
	Destroyed          bool                   `json:"destroyed"`
	Deleted            bool                   `json:"deleted"`
	DeletionTime       time.Time              `json:"deletion_time"`
}

// Engines struct that hols all engines key is the namespace.
type Engines map[string][]string

// Namespaces represents vault hierarchical namespaces.
type Namespaces map[string][]string
