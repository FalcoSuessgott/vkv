recursively list secrets from Vaults KV2 engine in various formats

```
vkv export [flags]
```

### Options

```
  -p, --path string              KVv2 Engine path (env: VKV_EXPORT_PATH)
  -e, --engine-path string       Specify the engine path using this flag in case your kv-engine contains special characters such as "/".
                                 vkv will then append the values of the path-flag to the engine path, if specified (<engine-path>/<path>)(env: VKV_EXPORT_ENGINE_PATH)
      --only-keys                show only keys (env: VKV_EXPORT_ONLY_KEYS)
      --only-paths               show only paths (env: VKV_EXPORT_ONLY_PATHS)
      --show-version             show the secret version (env: VKV_EXPORT_VERSION) (default true)
      --show-metadata            show the secrets metadata (env: VKV_EXPORT_METADATA) (default true)
      --show-values              don't mask values (env: VKV_EXPORT_SHOW_VALUES)
      --max-value-length int     maximum char length of values. Set to "-1" for disabling (env: VKV_EXPORT_MAX_VALUE_LENGTH) (default 12)
      --template-file string     path to a file containing Go-template syntax to render the KV entries (env: VKV_EXPORT_TEMPLATE_FILE)
      --template-string string   template string containing Go-template syntax to render KV entries (env: VKV_EXPORT_TEMPLATE_STRING)
  -f, --format string            available output formats: "base", "json", "yaml", "export", "policy", "markdown", "template" (env: VKV_EXPORT_FORMAT) (default "base")
  -h, --help                     help for export
```