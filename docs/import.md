# Import

`vkv import` requires an engine path (`--path`) and will accepts `vkv`s YAML or JSON output (`vkv export -f=yaml|json`) either by invoking `vkv import -` (for STDIN) or by specifying a file (`--file`). `vkv` will create the specified path if the engine does not exist yet and will error if it does, unless `--force` is specified. 

See the [CLI Reference](https://github.com/FalcoSuessgott/vkv/cmd/vkv_import/) for more details on the supported flags and env vars.

## Example Usage
```bash
> vkv export -p secret -f=yaml > secret_export.yaml
> vkv import -p copy --file=secret_export.yaml
reading secrets from secret_export.yaml
parsing secrets from YAML
writing secret "copy/admin" 
writing secret "copy/demo" 
writing secret "copy/sub/demo" 
writing secret "copy/sub/sub2/demo" 
successfully imported all secrets

result:

copy/ [type=kv2]
├── admin [v=1] [key=value]
│   └── sub=********
├── demo [v=1]
│   └── foo=***
└── sub
    ├── demo [v=1]
    │   ├── demo=***********
    │   ├── password=******
    │   └── user=*****
    └── sub2
        └── demo [v=2] [admin=false key=value]
            ├── admin=***
            ├── foo=***
            ├── password=********
            └── user=****
```

## Reading secrets from STDIN 

The `-` in `vkv import -`, tells `vkv` do read data via STDIN. The idea of `vkv import -` is, in order to copy/mirror KV-v2 secrets or complete engines across different Vault Servers or Namespaces, you can simply pipe 
`vkv`s output into the `vkv import` command:

```bash
# dont forget to use --show-values, otherwise the secrets will be uploaded masked.
vkv export -p <source> --show-values -f=yaml | vkv import - -p <destination>
```

### A few notes:
* `<source>` and `<destination>` dont have to be the root path of a secret engine, you also specify subpaths and copy them another secret engine.
* `vkv` will error if the enabled secret engine already exist, you can use `--force` to overwrite the destination engine, if the destination path contains a subpath (`root/sub`), `vkv` will then insert the secrets to that specific directory

**⚠️ `vkv import` can overwrite important secrets, always double check the command by using the dry-run mode (`--dry-run`) first**
