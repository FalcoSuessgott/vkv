<div align="center">
<h1>vkv</h1>
<img src="docs/assets/logo.png" alt="drawing" width="200"/>
<br>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vkv" alt="drawing"/>

`vkv` is a little CLI tool written in Go, which enables you to list, compare, import, document, backup & encrypt secrets from a [HashiCorp Vault KV engine](https://developer.hashicorp.com/vault/docs/secrets/kv):

<img src="docs/assets/demo.gif" alt="drawing"/>
</div>

## Features
* Support KV version 1 & version 2 (no need to specify the version `vkv` will automatically detect the engines version)
* **CI/CD Integrations for [Gitlab, GitHub, Azure Devops](https://falcosuessgott.github.io/vkv/cicd_gitlab)**
* support all Vault Auth Env Vars and `VKV_LOGIN_COMMAND` for avoiding having to hardcode the `VAULT_TOKEN` ([example](https://falcosuessgott.github.io/vkv/authentication/))
* recursively print secrets of any KV  Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/export_formats/)
* engine export shows the secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/export/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: `policy`)
* print secrets in `export <key>=<value>` format for env var exporting (format: `export`)
* move or migrate secrets from KVV1 to a KVV2 Engine or any subpath [example](https://falcosuessgott.github.io/vkv//example_migrate_kvv1_kvv2/)
* [import](https://falcosuessgott.github.io/vkv/import/) secrets back to Vault from `vkv`'s `json` or `yaml` format output
* save and restore KVv2 snapshots (including namespaces) ([kubernetes](https://falcosuessgott.github.io/vkv/example_kubernetes/) example)
* list all KVv2-engines or namespaces for scripting purposes ([fzf](https://falcosuessgott.github.io/vkv/example_fzf/) example)
* more handy [snippets](https://falcosuessgott.github.io/vkv/example_diff/) using `fzf`, `sops` & `diff`

Checkout the [Quickstart](https://falcosuessgott.github.io/vkv/quickstart/) Guide to learn more about `vkv`

## Quickstart

```bash
# Installation
version=$(curl https://api.github.com/repos/falcosuessgott/vkv/releases/latest -s | jq .name -r)
curl -OL "https://github.com/FalcoSuessgott/vkv/releases/download/${version}/vkv_$(uname)_$(uname -m).tar.gz"
tar xzf vkv_$(uname)_$(uname -m).tar.gz
chmod u+x vkv
./vkv version

# set required env vars
export VAULT_ADDR=https://vault-server:8200
export VAULT_TOKEN=<your-vault-token>

# verify connection
vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.12.1
Build Date      2022-10-27T12:32:05Z
Storage Type    inmem
Cluster Name    vault-cluster-ffd05212
Cluster ID      42ef92d5-eb21-0cb5-dd0b-804dac04e505
HA Enabled      false

# list secrets recursively of a KVv2 engine
vkv export --path <KVv2-engine path>
secret/
├── v1: admin [key=value]   # v1 -> secret version; "admin" -> secrets name; "[key=value]" -> secrets custom metadata
│   └── sub=********        # "sub" -> key; "*****" -> masked value (disable with --show-values)
├── v1: demo
│   └── foo=***
└── sub/
    ├── v1: demo
    │   ├── demo=***********
    │   ├── password=******
    │   └── user=*****
    └── sub2
        └── v2: demo [admin=false key=value]
            ├── admin=***
            ├── foo=***
            ├── password=********
            └── user=****
```