# vkv

[![Test](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml) [![golangci-lint](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/FalcoSuessgott/vkv)](https://goreportcard.com/report/github.com/FalcoSuessgott/vkv) [![codecov](https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg?token=UYVZ8LTA45)](https://codecov.io/gh/FalcoSuessgott/vkv)
[![Github all releases](https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg)](https://GitHub.com/FalcoSuessgott/vkv/releases/)

![img](assets/demo.gif)

# Description
`vkv` recursively list you all key-value entries from Vaults KV2 secret engine in various formats. `vkv` flags can be devided into input, modifying and output format flags.

So far `vkv` offers:

### Input flags
* `-p | --paths`: KV mount paths (comma separated list for multiple paths) (env: `VKV_PATHS`, default: `kv`)

### Modifying flags
* `--only-keys`: show only keys (env: `VKV_ONLY_KEYS`, default: `false`)
* `--only-paths`: show only paths (env: `VKV_ONLY_PATHS`, default: `false`)
* `-show-values`: dont mask values (env: `VKV_SHOW_VALUES`, default: `false`)
* `--max-value-length`: maximum char length of values (set to `-1` for disabling) (env: `VKV_MAX_VALUE_LENGTH`, default: `12`)

### Output Flags
* `-f | --format`: output format (options: `base`, `yaml`, `json`, `export`, `markdown`)  (env: `"VKV_FORMAT"`, default: `"base"`)

⚠️ **A flag always preceed its environment variable**

You can combine most of those flags in order to receive the desired output.
For examples see the [Examples](https://github.com/FalcoSuessgott/vkv#examples)

# Installation
Find the corresponding binaries, `.rpm` and `.deb` packages in the [release](https://github.com/FalcoSuessgott/vkv/releases) section.

# Supported OS and Vault Versions
`vkv` is being tested on `Windows`, `MacOS` and `Linux` and also against Vault Version < `v1.8.0` (but it also may work with lower versions).

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

# Examples
Imagine you have the following KV2 structure mounted at path `secret/`:

```
secret/
  demo
    foo=bar

  sub
    sub=passw0rd

  sub/demo
    demo="hello world"
    password=s3cre5
    user=admin

  sub/sub2/demo
    value=nevermind
    password=secret2
    user=database
```

## Input
### list secrets (`--path` | `-p` | `VKV_PATHS="kv1:kv2"`)
You can list all secrets recursively by running:

```bash
$> vkv --path secret
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
            ├── user=************
            └── value=*********
```

You can also specifiy a specific subpaths:

```bash
$> vkv --path secret/sub/sub2
secret/sub/sub2/
└── sub/
    └── sub2/
        └── demo
            ├── user=************
            └── value=*********
```

and list as much paths as you want:

```bash
# or as comma separated with no spaces!
$> vkv -p secret -p secret2
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
            ├── user=************
            └── value=*********
secret_2/
├── demo
│   └── foo=***
├── sub
│   └── sub=********
├── sub/
│   └── demo
│       ├── foo=***
│       ├── password=********
│       └── user=****
└── sub/
    └── sub2/
        └── demo
            ├── foo=***
            ├── password=********
            └── user=****
```

## Modifying
### list only paths (`--only-paths` | `VKV_ONLY_PATHS=true`)
We can receive only the paths by running

```bash
$> vkv  -p secret --only-paths
secret/
├── demo
├── sub
├── sub/
│   └── demo
└── sub/
    └── sub2/
        └── demo
```

### list only secret keys  (`--only-keys` | `VKV_ONLY_KEYS=true`)
If we want to know just the keys in every directory we can run

```bash
$> vkv -p secret --only-keys
secret/
├── demo
│   └── foo
├── sub
│   └── sub
├── sub/
│   └── demo
│       ├── demo
│       ├── password
│       └── user
└── sub/
    └── sub2/
        └── demo
            ├── user
            └── value
```

### show values  (`--show-values` | `VKV_SHOW_VALUES=true`)
Per default values are masked. Using `--show-values` shows the values. **Use with Caution**

```bash
$> vkv -p secret --show-values
secret/
├── demo
│   └── foo=bar
├── sub
│   └── sub=password
├── sub/
│   └── demo
│       ├── demo=hello world
│       ├── password=s3cre5
│       └── user=admin
└── sub/
    └── sub2/
        └── demo
            ├── user=databasepassword=secret2
            └── value=nevermind
```

## Output Format
### export format (`--format=export` | `VKV_FORMAT=export`)
You can print out the entries in `export key=value` format for further processing:

```bash
$> vkv --path secret/sub/sub2 --format=export
export demo="hello world"
export password="s3cre5"
export user="admin"
export user="databasepassword=secret2"
export value="nevermind"
export foo="bar"
export sub="password
```

You can then use `eval` to source those env vars:

```bash
echo $foo # not defined
eval $(vkv -f=export --path secret/sub/sub2)
echo $foo
"bar" # value under the specific key exported
```

## markdown (`--format=markdown` | `VKV_FORMAT=markdown`)
```bash
vkv -p secret --format=markdown
```

returns:

| MOUNT  |        PATHS         |   KEYS   |    VALUES    |
|--------|----------------------|----------|--------------|
| secret | secret/demo          | foo      | ***          |
|        | secret/sub           | sub      | ********     |
|        | secret/sub/demo      | demo     | ***********  |
|        |                      | password | ******       |
|        |                      | user     | *****        |
|        | secret/sub/sub2/demo | user     | ************ |
|        |                      | value    | *********    |


### json (`--format=json` | `VKV_FORMAT=json`)
You can combine all flags and export the result to json by running:

```bash
vkv -p secret --show-values --format=json
```

```json
{
  "secret": {
    "secret/demo": {
      "foo": "***"
    },
    "secret/sub": {
      "sub": "********"
    },
    "secret/sub/demo": {
      "demo": "***********",
      "password": "******",
      "user": "*****"
    },
    "secret/sub/sub2/demo": {
      "user": "************",
      "value": "*********"
    }
  }
}%
```

### yaml (`--format=yaml` | `VKV_FORMAT=yaml`)
Same applies for yaml:

```bash
vkv --path secret --show-values --format=yaml
```

```yaml
secret:
  secret/demo:
    foo: '***'
  secret/sub:
    sub: '********'
  secret/sub/demo:
    demo: '***********'
    password: '******'
    user: '*****'
  secret/sub/sub2/demo:
    user: '************'
    value: '*********'
```

# Acknowledgements / Similar tools
`vkv` is inspired by:
* https://github.com/jonasvinther/medusa

Similar tools are:
* https://github.com/kir4h/rvault
