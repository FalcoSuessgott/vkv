---
title: "vkv snapshot save"
weight: 6
---

create a snapshot of all visible KV engines recursively for all namespaces

```
vkv snapshot save [flags]
```

### Options

```
  -d, --destination string   vkv snapshot destination path (env: VKV_SNAPSHOT_SAVE_DESTINATION) (default "./vkv-snapshot-export")
  -h, --help                 help for save
  -n, --namespace string     namespaces from which to save recursively all visible KV engines (env: VKV_SNAPSHOT_SAVE_NS
```