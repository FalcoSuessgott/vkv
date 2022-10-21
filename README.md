<div align="center">
  <h1> vkv </h1>
  <img src="assets/base.svg" alt="drawing" height="400" width="550">

  [![Test](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml) [![golangci-lint](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/FalcoSuessgott/vkv)](https://goreportcard.com/report/github.com/FalcoSuessgott/vkv) [![codecov](https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg?token=UYVZ8LTA45)](https://codecov.io/gh/FalcoSuessgott/vkv)
[![Github all releases](https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg)](https://GitHub.com/FalcoSuessgott/vkv/releases/)
</div>


# Description
`vkv` recursively list you all key-value entries from Vaults KV2 secret engine in various formats. `vkv` flags can be devided into input, modifying and output format flags.

So far `vkv` offers:

| Flag                 | Description                                                                       | Env Var                | Default |
|----------------------|-----------------------------------------------------------------------------------|------------------------|---------|
| `-p`, `--paths`      | KV mount paths (comma separated list for multiple paths)                          | `VKV_PATHS`            | `kv`    |
| `-f`, `--format`     | output format (options: `base`, `yaml`, `json`, `export`, `markdown`, `template`) | `VKV_FORMAT`           | `base`  |
| `--only-keys`        | show only keys                                                                    | `VKV_ONLY_KEYS`        | `false` |
| `--only-paths`       | show only paths                                                                   | `VKV_ONLY_PATHS`       | `false` |
| `--show-values`      | dont mask values                                                                  | `VKV_SHOW_VALUES`      | `false` |
| `--max-value-length` | maximum char length of values (set to `-1` for disabling)                         | `VKV_MAX_VALUE_LENGTH` | `12`    |
| `--template-file`    | path to a file containing Go-template syntax to render the KV entries             | `VKV_TEMPLATE_FILE`    |         |
| `--template-string`  | string containing Go-template syntax to render KV entries                         | `VKV_TEMPLATE_STRING`  |         |


⚠️ **A flag always preceed its environment variable**

You can combine most of those flags in order to receive the desired output.

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
### `base`
![](assets/base.svg)

### `markdown`
![](assets/markdown.svg)

### `json`
![](assets/json.svg)

### `yaml`
![](assets/yaml.svg)

### `template`
![](assets/template.svg)



# Development
Clone this repository and run:

```sh
make bootstrap
```

in order to have all used build dependencies

You can spin up a development vault for local testing by running:

```sh
make vault
```

The following environment variables are required:

```sh
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"
export VKV_PATHS="secret"
```

If everthing worked fine, you should be able to run:

```sh
go run main.go   
secret/
├── demo
│   └── foo=***
├── sub
│   └── sub=********
├── sub/
│   └── demo
│       ├── demo=***********
│       ├── password=******
│       └── user=*****
└── sub/
    └── sub2/
        └── demo
            ├── password=*******
            ├── user=********
            └── value=*********
```

