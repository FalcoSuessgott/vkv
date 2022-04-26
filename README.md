<div align="center">
  <h1> vkv </h1>
  <img src="assets/base.svg" alt="drawing" height="400" width="550">

  [![Test](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml) [![golangci-lint](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/FalcoSuessgott/vkv)](https://goreportcard.com/report/github.com/FalcoSuessgott/vkv) [![codecov](https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg?token=UYVZ8LTA45)](https://codecov.io/gh/FalcoSuessgott/vkv)
[![Github all releases](https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg)](https://GitHub.com/FalcoSuessgott/vkv/releases/)
</div>


# Description
`vkv` recursively list you all key-value entries from Vaults KV2 secret engine in various formats. `vkv` flags can be devided into input, modifying and output format flags.

So far `vkv` offers:

### Input flags
* `-p | --paths`: KV mount paths (comma separated list for multiple paths) (env: `VKV_PATHS`, default: `kv`)

### Modifying flags
* `--only-keys`: show only keys (env: `VKV_ONLY_KEYS`, default: `false`)
* `--only-paths`: show only paths (env: `VKV_ONLY_PATHS`, default: `false`)
* `--show-values`: dont mask values (env: `VKV_SHOW_VALUES`, default: `false`)
* `--max-value-length`: maximum char length of values (set to `-1` for disabling) (env: `VKV_MAX_VALUE_LENGTH`, default: `12`)
* `--template-file`: path to a file containing Go-template syntax to render the KV entries (env: `VKV_TEMPLATE_FILE`)
* `--template-string`: string containting Go-template syntax to render KV entries (env: `VKV_TEMPLATE_STRING`)

### Output Flags (see [Supported Formats](https://github.com/FalcoSuessgott/vkv/tree/template#supported-formats))
* `-f | --format`: output format (options: `base`, `yaml`, `json`, `export`, `markdown`, `template`)  (env: `"VKV_FORMAT"`, default: `"base"`)

⚠️ **A flag always preceed its environment variable**

You can combine most of those flags in order to receive the desired output.
For examples see the [Examples](https://github.com/FalcoSuessgott/vkv#examples)

# Installation
Find the corresponding binaries, `.rpm` and `.deb` packages in the [release](https://github.com/FalcoSuessgott/vkv/releases) section.

# Supported OS and Vault Versions
`vkv` is being tested on `Windows`, `MacOS` and `Linux` and also against Vault Version >= `v1.8.0` (but it also may work with lower versions).

# Authentication
`vkv` supports token based authentication. It is clear that you can only see the secrets that are allowed by your token policy.

All of vaults [environment variables](https://www.vaultproject.io/docs/commands#environment-variables) are supported. In order to authenticate to a Vault instance you have to set atleast `VAULT_ADDR` and `VAULT_TOKEN`:

```bash
# on linux/macos
VAULT_ADDR="http://127.0.0.1:8200" VAULT_TOKEN="s.XXX" vkv -p <kv-path>

# on windows
SET VAULT_ADDR=http://127.0.0.1:8200
SET VAULT_TOKEN=s.XXX
vkv.exe -p <kv-path>
```

# Supported Formats
|  |                          |
|:-------------:|:-------------------------------:|
| `base`<br> ![](assets/base.svg)| `markdown`<br> ![](assets/markdown.svg) |
| `json`<br> ![](assets/json.svg)| `yaml`<br> ![](assets/yaml.svg) |

| |
|:---:|
| `template`<br> <img src="assets/template.svg" width="600" /> |


