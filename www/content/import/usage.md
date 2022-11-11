---
title: "Usage"
weight: 1
---

The `-` in `vkv import -`, tells `vkv` do read data via STDIN. The idea of `vkv import -` is, in order to copy/mirror KV-v2 secrets or complete engines across diferrent Vault Servers or Namespaces, you can simply pipe 
`vkv`s output into the `vkv import` command:

```bash
 # dont forget to use --show-values, otherwise the secrets will be uploaded masked.
vkv -p <source> --show-values -f=yaml | vkv import - -p <destination>
```

### Demo
<div align="center">
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/export-import.gif" alt="drawing" width="1000"/>
</div>

### A few notes:
* `<source>` and `<destination>` dont have to be the root path of a secret engine, you also specify subpaths and copy them another secret engine.
* `vkv` will error if the enabled secret engine already exist, you can use `--force` to overwrite the destination engine, if the destination path contains a subpath (`root/sub`), `vkv` will then insert the secrets to that specific directory

**⚠️ `vkv import` can overwrite important secrets, always double check the commmand by using the dry-run mode (`--dry-run`) first**

## Input flags
| Flag                  | Description                                                                       | Env                    | Default |
|-----------------------|-----------------------------------------------------------------------------------|------------------------|---------|
| `-p`, `--path`        | Destination KV-v2 path                                                            | `VKV_IMPORT_PATH`      |         |
| `-f`, `--file`        | path to a file containing vkv JSON or YAML output                                 | `VKV_IMPORT_FILE`      |         |
| `-d`, `--dry-run`     | Dry-run does not upload any secrets and just prints a preview                     | `VKV_IMPORT_DRY_RUN`   |         |
| `-s`, `--silent`      | do not show the resulting secret engine                                           | `VKV_IMPORT_SILENT`    | `false` |
| `--force`             | force overwrites the specified secret engine. **Use with Caution**                | `VKV_IMPORT_FORCE`     | `false` |
| `---max-value-length` | maximum char length of values (set to `-1` for disabling)                         | `VKV_IMPORT_MAX_VALUE_LENGTH` | `12`    | 
| `--show-values`       | dont mask values when printing the resulting secret engine                        | `VKV_IMPORT_SHOW_VALUES`| `false`|

⚠️ **A flag always precede its environment variable**