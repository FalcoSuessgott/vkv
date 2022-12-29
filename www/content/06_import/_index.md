---
title: "vkv import "
weight: 6
---

import secrets from vkv's json or yaml output

```
vkv import [flags]
```

### Options

```
  -d, --dry-run                print resulting KV engine (env: VKV_IMPORT_DRY_RUN)
  -f, --file string            path to a file containing vkv yaml or json output (env: VKV_IMPORT_FILE)
      --force                  overwrite existing kv entries (env: VKV_IMPORT_FORCE)
  -h, --help                   help for import
      --max-value-length int   maximum char length of values. Set to "-1" for disabling (env: VKV_IMPORT_MAX_VALUE_LENGTH) (default 12)
  -p, --path string            KVv2 Engine path (env: VKV_IMPORT_PATH)
      --show-values            don't mask values (env: VKV_IMPORT_SHOW_VALUES)
  -s, --silent                 do not output secrets (env: VKV_IMPORT_SILENT)
```


# read secrets from STDIN 

The `-` in `vkv import -`, tells `vkv` do read data via STDIN. The idea of `vkv import -` is, in order to copy/mirror KV-v2 secrets or complete engines across diferrent Vault Servers or Namespaces, you can simply pipe 
`vkv`s output into the `vkv import` command:

```bash
# dont forget to use --show-values, otherwise the secrets will be uploaded masked.
vkv -p <source> --show-values -f=yaml | vkv import - -p <destination>
```

### A few notes:
* `<source>` and `<destination>` dont have to be the root path of a secret engine, you also specify subpaths and copy them another secret engine.
* `vkv` will error if the enabled secret engine already exist, you can use `--force` to overwrite the destination engine, if the destination path contains a subpath (`root/sub`), `vkv` will then insert the secrets to that specific directory

**⚠️ `vkv import` can overwrite important secrets, always double check the commmand by using the dry-run mode (`--dry-run`) first**