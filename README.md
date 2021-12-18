# vkv

[![Test](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml) [![golangci-lint](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg)](https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/FalcoSuessgott/vkv)](https://goreportcard.com/report/github.com/FalcoSuessgott/vkv) [![Go Reference](https://pkg.go.dev/badge/github.com/FalcoSuessgott/vkv.svg)](https://pkg.go.dev/github.com/FalcoSuessgott/vkv) [![codecov](https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg?token=UYVZ8LTA45)](https://codecov.io/gh/FalcoSuessgott/vkv)


> recursively list secrets from Vaults KV2 engine

![img](assets/example.png)

# Installation
Find the corresponding binaries, `.rpm` and `.deb` packages in the [release](https://github.com/FalcoSuessgott/vkv/releases) section.

# Authentication
`vkv` supports token based authentication. It is clear that you can only see the secrets that are allowed by your token policy.

## Required Environment Variables
In order to authenticate to a Vault instance you have to export `VAULT_ADDR` and `VAULT_TOKEN`.

```bash
VAULT_ADDR="http://127.0.0.1:8200" VAULT_TOKEN="root" vkv
```

## Optional  Environment Variables
Furthermore you can export:

* `VAULT_NAMESPACE` for namespace login
* `VAULT_SKIP_VERIFY` for insecure HTTPS connection
* `HTTP_PROXY` and `HTTPS_PROXY` for proxy connections.

# Usage
```bash
$> vkv -h
recursively list secrets from Vaults KV2 engine

Usage:
  vkv [flags]

Flags:
  -h, --help           help for vkv
      --only-keys      print only keys
      --only-paths     print only paths
  -p, --path string    path (default "kv")
      --show-secrets   print out secrets
  -j, --to-json        print secrets in json format
  -y, --to-yaml        print secrets in yaml format
  -v, --version        display version
```

# Walkthrough
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

You can specifiy specific subpaths:

```bash
vkv --path secret/sub/sub2
secret/sub/sub2/demo
        foo=***
        password=********
        user=****
```

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

### show secrets  `--show-secrets`
Per default secret values are masked. Using `--show-secrets` shows the secrets. **Use with Caution**

We can get the secrets of a certain sub path, by running

```bash
vkv -p secret --show-secrets
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

### export to json `--to-json | -j`
You can combine all flags and export the result to json by running:

```bash
vkv -p secret --sub-path sub --show-secrets --to-json | jq .
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

### export to yaml  `--to-yaml | -y`
Same applies for yaml:

```bash
vkv --path secret --sub-path sub --show-secrets --to-yaml
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
