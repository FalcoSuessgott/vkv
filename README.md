# vkv

[![Test](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml) [![golangci-lint](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/FalcoSuessgott/vkv)](https://goreportcard.com/report/github.com/FalcoSuessgott/vkv) [![Go Reference](https://pkg.go.dev/badge/github.com/FalcoSuessgott/vkv.svg)](https://pkg.go.dev/github.com/FalcoSuessgott/vkv) [![codecov](https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg?token=UYVZ8LTA45)](https://codecov.io/gh/FalcoSuessgott/vkv)


![img](assets/example.png)

# Description
`vkv` recursively list you all key-value entries from Vaults KV2 secret engine in various formats. `vkv` flags can be devided into input, modifying and output format flags.

So far `vkv` offers:

### Input flags
* `-p | --paths` (default: `kv`): Comma separated list of KVv2 Engine Paths

### Modifying flags
* `--only-keys`: show only keys
* `--only-paths`: show only paths
* `-show-values`: dont mask values
* `--max-value-length` (default: `12`): maximum char length of values (precedes `VKV_MAX_PASSWORD_LENGTH`). Set to `-1` for disabling

### Output Flags
* `-e | --export`: print entries in export format (export "key=value")
* `-j | --json`: print entries in json format
* `-y | --yaml`: print entries in yaml format
* `-m | --markdown`: print entries in markdown table format

You can combine most of those flags in order to receive the desired output.
For examples see the [Examples](https://github.com/FalcoSuessgott/vkv#examples)

# Installation
Find the corresponding binaries, `.rpm` and `.deb` packages in the [release](https://github.com/FalcoSuessgott/vkv/releases) section.

# Supported OS and Vault Versions
`vkv` is being tested on `Windows`, `MacOS` and `Ubuntu` and also against Vault Version < `v1.8.0` (but it also may work with lower versions).

# Authentication
`vkv` supports token based authentication. It is clear that you can only see the secrets that are allowed by your token policy.

### Required Environment Variables
In order to authenticate to a Vault instance you have to export `VAULT_ADDR` and `VAULT_TOKEN`.

```bash
# on linux/macos
VAULT_ADDR="http://127.0.0.1:8200" VAULT_TOKEN="s.XXX" vkv -p <kv-path>

# on windows
SET VAULT_ADDR=http://127.0.0.1:8200
SET VAULT_TOKEN=s.XXX
vkv.exe -p <kv-path>
```

### Optional Environment Variables
Furthermore you can export:

* `VAULT_NAMESPACE` for namespace login
* `VAULT_SKIP_VERIFY` for insecure HTTPS connection
* `HTTP_PROXY` and `HTTPS_PROXY` for proxy connections.

# Examples
Imagine you have the following KV2 structure mounted at path `secret/`:

```
secret/
  demo
    foo=bar

  sub
    sub=passw0rd

  sub/demo
    foo=bar
    password=passw0rd
    user=user

  sub/sub2/demo
    foo=bar
    password=passw0rd
    user=user
```

## Input
### list secrets `--path | -p (default "kv")`
You can list all secrets recursively by running:

```bash
vkv --path secret
secret/demo
        foo=***
secret/sub
        sub=********
secret/sub/demo
        foo=***
        password=********
        user=****
secret/sub/sub2/demo
        foo=***
        password=********
        user=****
```

You can also specifiy a specific subpaths:

```bash
vkv --path secret/sub/sub2
secret/sub/sub2/demo
        foo=***
        password=********
        user=****
```

and list as much paths as you want:

```bash
# comma separated and no spaces!
vkv -p secret,secret2
secret/demo
        foo=***
secret/sub
        sub=********
secret/sub/demo
        foo=***
        password=********
        user=****
secret/sub/sub2/demo
        foo=***
        password=********
        user=****
secret2/demo
        user=********
```

## Modifying
### list only paths `--only-paths`
We can receive only the paths by running

```bash
vkv  -p secret --only-paths
secret/demo
secret/sub
secret/sub/demo
secret/sub/sub2/demo
```

### list only secret keys  `--only-keys`
If we want to know just the keys in every directory we can run

```bash
vkv -p secret --only-keys
secret/demo
        foo
secret/sub
        sub
secret/sub/demo
        foo
        password
        user
secret/sub/sub2/demo
        foo
        password
        user
```

### show values  `--show-values`
Per default values are masked. Using `--show-values` shows the values. **Use with Caution**

We can get the secrets of a certain sub path, by running

```bash
vkv -p secret --show-values
secret/demo
        foo=bar
secret/sub
        sub=password
secret/sub/demo
        foo=bar
        password=password
        user=user
secret/sub/sub2/demo
        foo=bar
        password=password
        user=user
```

## Output Format
### export format `--export | -e`
You can print out the entries in `export key=value` format for further processing:

```bash
vkv --path secret/sub/sub2 --export
export foo=secret1
export password=secret2
export user=secret3
```

You can then use `eval` to source those env vars:

```bash
echo $foo # not defined
eval $(vkv --export --path secret/sub/sub2)
echo $foo
"secret1" # value under the specific key exported
```

## markdown `--markdown | -m`
```bash
vkv -p secret --markdown
```

returns:

|        PATHS         |   KEYS   |  VALUES  |
|----------------------|----------|----------|
| secret/demo          | foo      | ***      |
| secret/sub           | sub      | ******** |
| secret/sub/demo      | foo      | ***      |
|                      | password | ******** |
|                      | user     | ****     |
| secret/sub/sub2/demo | foo      | ***      |
|                      | password | ******** |
|                      | user     | ****     |

In combination with:

`--only-paths`:
|        PATHS         |
|----------------------|
| secret/demo          |
| secret/sub           |
| secret/sub/demo      |
| secret/sub/sub2/demo |

`--only-keys`:
|        PATHS         |   KEYS   |
|----------------------|----------|
| secret/demo          | foo      |
| secret/sub           | sub      |
| secret/sub/demo      | foo      |
|                      | password |
|                      | user     |
| secret/sub/sub2/demo | user     |
|                      | foo      |
|                      | password |


### json `--json | -j`
You can combine all flags and export the result to json by running:

```bash
vkv -p secret --show-values --json | jq .
```

```json
{
  "secret/demo": {
    "foo": "bar"
  },
  "secret/sub": {
    "sub": "password"
  },
  "secret/sub/demo": {
    "foo": "bar",
    "password": "password",
    "user": "user"
  },
  "secret/sub/sub2/demo": {
    "foo": "bar",
    "password": "password",
    "user": "user"
  }
}
```

### yaml  `--yaml | -y`
Same applies for yaml:

```bash
vkv --path secret --show-values --yaml
```

```yaml
secret/demo:
  foo: bar
secret/sub:
  sub: password
secret/sub/demo:
  foo: bar
  password: password
  user: user
secret/sub/sub2/demo:
  foo: bar
  password: password
  user: user
```

# Acknowledgements / Similar tools
`vkv` is inspired by:
* https://github.com/jonasvinther/medusa

Similar tools are:
* https://github.com/kir4h/rvault
