---
title: "Quickstart"
weight: 1
---

This guide will run you through some of the features of `vkv`.

### 0. Prerequsuites
In order to perform all of the described tasks, you will need the following tools:

* a Linux/MacOS Shell
* `docker` installed and running (alternatively `vault` CLI can be used)
* `vkv` installed (follow https://falcosuessgott.github.io/vkv/installation/)

### 01. Spin up a development Vault server
First, we setup a development Vault server. 

Open a terminal and run:

```bash
docker run -p 8200:8200 hashicorp/vault server -dev -dev-root-token-id=root
```

You should then be able to visit `http://127.0.0.1:8200` in your browser and see a Vault login page.

The `root` token is `root`.

### 02. Verify connection
Once you have exported the required environment variables, you can verify your connection with the vault CLI by running:

```bash
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"
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
Cluster Name    vault-cluster-1bbeabe8
Cluster ID      405e99ab-5f8c-18ca-dafb-228c91add2f4
HA Enabled      false
```

If `vault status` returned an output like this you good to go to the next step


### 03. Write secrets to Vault using `vault` 
In a development Vault server a `KVv2` under `secret/` is enabled by default.
We want to write some secrets using `vault`:

```bash
vault kv put -mount=secret admin username=user password=passw0rd
vault kv metadata put -mount=secret -custom-metadata=key=value admin
vault kv put -mount=secret db/prod env=prod username=user password=passw0rd-prod
vault kv put -mount=secret db/dev env=dev username=user password=passw0rd-dev
```

 
### 04. List secrets using `vkv`
`vkv` requires atleast `VAULT_ADDR` and `VAULT_TOKEN` if the `vault status` command works, `vkv` will also work.

We can now use `vkv` to list all of our secrets recursively:

```bash
vkv export --path secret                                            
secret/
├── v1: admin [key=value]
│   ├── password=********
│   └── username=****
└── db/
    ├── v1: dev
    │   ├── env=***
    │   ├── password=************
    │   └── username=****
    └── v1: prod
        ├── env=****
        ├── password=************
        └── username=****
```

Here are some explanations:
* `vkv` masks the secrets per default, you can disable this by using `--show-values` or `VKV_EXPORT_SHOW_VALUES=true`
* `vkv` limits the length of the secrets per default to `12` for readability purposes (You can set you own value length by using `--max-value-length=XX` or `VKV_EXPORT_MAX_VALUE_LENGTH=XX`)
* `v1` indicates the secret version (disable by using `--show-version` or `VKV_EXPORT_SHOW_VERSION=false` 
* `[key=value]` represents the custom metadata that we added to the secret in step 3. (disable by `--show-metadata` or (`VKV_EXPORT_SHOW_METADATA=false`)


This output format is the default format called `base`. `vkv` has many other useful output formats. 

You can see them all using this onliner:

```bash
for f in base yaml json markdown policy export; do
echo -n "\n===> Output Format: $f <===\n"
vkv export -p secret --format=$f;
done 

===> Output Format: base <===
secret/
├── v1: admin [key=value]
│   ├── password=********
│   └── username=****
└── db/
    ├── v1: dev
    │   ├── env=***
    │   ├── password=************
    │   └── username=****
    └── v1: prod
        ├── env=****
        ├── password=************
        └── username=****

===> Output Format: yaml <===
secret/:
  admin:
    password: '********'
    username: '****'
  db/:
    dev:
      env: '***'
      password: '************'
      username: '****'
    prod:
      env: '****'
      password: '*************'
      username: '****'


===> Output Format: json <===
{
  "secret/": {
    "admin": {
      "password": "********",
      "username": "****"
    },
    "db/": {
      "dev": {
        "env": "***",
        "password": "************",
        "username": "****"
      },
      "prod": {
        "env": "****",
        "password": "*************",
        "username": "****"
      }
    }
  }
}

===> Output Format: markdown <===
|      PATH      |   KEY    |    VALUE     | VERSION | METADATA  |
|----------------|----------|--------------|---------|-----------|
| secret/admin   | password | ********     |       1 | key=value |
|                | username | ****         |         |           |
| secret/db/dev  | env      | ***          |       1 |           |
|                | password | ************ |         |           |
|                | username | ****         |         |           |
| secret/db/prod | env      | ****         |       1 |           |
|                | password | ************ |         |           |
|                | username | ****         |         |           |

===> Output Format: policy <===
PATH            CREATE  READ    UPDATE  DELETE  LIST    ROOT
secret/admin    ✖       ✖       ✖       ✖       ✖       ✔
secret/db/dev   ✖       ✖       ✖       ✖       ✖       ✔
secret/db/prod  ✖       ✖       ✖       ✖       ✖       ✔

===> Output Format: export <===
export env='prod'
export password='passw0rd-prod'
export username='user'
```

Most of these formats, offer various commandline flags, such as `--show-secrets`, `--only-paths`, `--only-keys`, `--max-value-length` to modify the output. These flags can also be set through environment variables:

```bash
VKV_EXPORT_FORMAT=JSON VKV_EXPORT_SHOW_VALUES=true vkv export -p secret
{
  "secret/": {
    "admin": {
      "password": "passw0rd",
      "username": "user"
    },
    "db/": {
      "dev": {
        "env": "dev",
        "password": "passw0rd-dev",
        "username": "user"
      },
      "prod": {
        "env": "prod",
        "password": "passw0rd-prod",
        "username": "user"
      }
    }
  }
}
```

### 05. Import secrets using `vkv`
Meanwhile `vkv export` can be used to store secrets, `vkv import` is used to import secrets from a `vkv export` command (either `yaml` or `json` format is accepted).

Knowing this, we can copy a secret engine to another secret engine:

```bash
vkv export -p secret -f=yaml --show-values | vkv import - -p copy
reading secrets from STDIN
parsing secrets from YAML
writing secret "copy/db/dev" 
writing secret "copy/db/prod" 
writing secret "copy/admin" 
successfully imported all secrets

result:

copy/
├── v1: admin
│   ├── password=********
│   └── username=****
└── db/
    ├── v1: dev
    │   ├── env=***
    │   ├── password=************
    │   └── username=****
    └── v1: prod
        ├── env=****
        ├── password=************
        └── username=****
```

Or even to another Vault instance:

```bash
vkv export -p secret -f=yaml --show-values| VAULT_ADDR="..." VAULT_TOKEN="..." vkv import - -p copy
[...]
```

Dont forget to set `--show-values` otherwise `vkv` will import masked `secrets`.

The `-` tells `vkv` to read the secrets from STDIN. You cal also specifiy a file using the `--file` parameter.

`vkv` will create the `KVv2` engine if it doesn't exist. If the engine indeed exists, `vkv` will error unless `--force` is used.

You can also copy subpaths to other engines:

```bash
vkv export -p secret/admin -f=yaml --show-values| vkv import - -p admin
reading secrets from STDIN
parsing secrets from YAML
writing secret "admin/admin" 
successfully imported all secrets

result:

admin/
└── v1: admin
    ├── password=********
    └── username=****
```

### 06. Create KVv2 Snapshots using `vkv`
`vkv` enables you to create and restore snapshots of all KVv2 engines in all namespaces of a Vault instance (requires an appropiate token + policy):

Consider the following namespaces and KVv2 engines on a Vault Enterprise instance:

```bash
# list all namespaces
vkv list namespaces --all
sub
sub/sub2
test
test/test2
test/test2/test3

# list all engines with their respective namespace as the prefix
vkv list engines --all --include-ns-prefix
secret
secret_2
sub/sub2/sub_sub2_secret
sub/sub2/sub_sub2_secret_2
sub/sub_secret
sub/sub_secret_2
test/test2/test3/test_test2_test3_secret
test/test2/test3/test_test2_test3_secret_2
```

You can create a snapshot of those KVv2 engines by running:

```bash
vkv snapshot save --destionation vkv-export-$(date '+%Y-%m-%d')
created vkv-export-2022-12-29
created vkv-export-2022-12-29/secret.yaml
created vkv-export-2022-12-29/secret_2.yaml
created vkv-export-2022-12-29/sub
created vkv-export-2022-12-29/sub/sub_secret_2.yaml
created vkv-export-2022-12-29/sub/sub_secret.yaml
created vkv-export-2022-12-29/sub/sub2
created vkv-export-2022-12-29/sub/sub2/sub_sub2_secret.yaml
created vkv-export-2022-12-29/sub/sub2/sub_sub2_secret_2.yaml
created vkv-export-2022-12-29/test
created vkv-export-2022-12-29/test/test2
created vkv-export-2022-12-29/test/test2/test3
created vkv-export-2022-12-29/test/test2/test3/test_test2_test3_secret.yaml
created vkv-export-2022-12-29/test/test2/test3/test_test2_test3_secret_2.yaml
```

As you cann see: `vkv` exported all engines and wrote them to the specified directory:

```bash
vkv-export-2022-12-29/
├── secret_2.yaml
├── secret.yaml
├── sub
│   ├── sub2
│   │   ├── sub_sub2_secret_2.yaml
│   │   └── sub_sub2_secret.yaml
│   ├── sub_secret_2.yaml
│   └── sub_secret.yaml
└── test
    └── test2
        └── test3
            ├── test_test2_test3_secret_2.yaml
            └── test_test2_test3_secret.yaml

5 directories, 8 files
```

whereas one file is the JSON output of a single KVv2 engine:

```bash
cat vkv-export-2022-12-29/secret.yaml  
{
  "admin": {
    "sub": "password"
  },
  "demo": {
    "foo": "bar"
  },
  "sub/": {
    "demo": {
      "demo": "hello world",
      "password": "s3cre5",
      "user": "admin"
    },
    "sub2/": {
      "demo": {
        "admin": "key",
        "foo": "bar",
        "password": "password",
        "user": "user"
      }
    }
  }
}
```

You could `.tar.gz` those directories and save those encrypted files in a secure fashion.

### 07. Restore vkv snapshots

In order to restore a `vkv` snapshot the `snapshot restore` command is invoked:

```bash
# no KVv2 engines configured
vkv list engines --all --include-ns-prefix                      
[ERROR] no engines found.

# restore a snapshot
vkv snapshot restore --source vkv-export-2022-12-29
[root] restore engine: secret
[root] writing secret "secret/admin" 
[root] writing secret "secret/demo" 
[root] writing secret "secret/sub/demo" 
[root] writing secret "secret/sub/sub2/demo" 
[root] restore engine: secret_2
[root] writing secret "secret_2/admin" 
[root] writing secret "secret_2/demo" 
[root] writing secret "secret_2/sub/demo" 
[root] writing secret "secret_2/sub/sub2/demo" 
[root] restore namespace: "sub"
[sub] restore namespace: "sub2"
[sub/sub2] restore engine: sub_sub2_secret
[sub/sub2] writing secret "sub_sub2_secret/admin" 
[sub/sub2] writing secret "sub_sub2_secret/demo" 
[sub/sub2] writing secret "sub_sub2_secret/sub/demo" 
[sub/sub2] writing secret "sub_sub2_secret/sub/sub2/demo" 
[sub/sub2] restore engine: sub_sub2_secret_2
[sub/sub2] writing secret "sub_sub2_secret_2/admin" 
[sub/sub2] writing secret "sub_sub2_secret_2/demo" 
[sub/sub2] writing secret "sub_sub2_secret_2/sub/sub2/demo" 
[sub/sub2] writing secret "sub_sub2_secret_2/sub/demo" 
[sub] restore engine: sub_secret
[sub] writing secret "sub_secret/admin" 
[sub] writing secret "sub_secret/demo" 
[sub] writing secret "sub_secret/sub/demo" 
[sub] writing secret "sub_secret/sub/sub2/demo" 
[sub] restore engine: sub_secret_2
[sub] writing secret "sub_secret_2/sub/demo" 
[sub] writing secret "sub_secret_2/sub/sub2/demo" 
[sub] writing secret "sub_secret_2/admin" 
[sub] writing secret "sub_secret_2/demo" 
[root] restore namespace: "test"
[test] restore namespace: "test2"
[test/test2] restore namespace: "test3"
[test/test2/test3] restore engine: test_test2_test3_secret
[test/test2/test3] writing secret "test_test2_test3_secret/sub/sub2/demo" 
[test/test2/test3] writing secret "test_test2_test3_secret/admin" 
[test/test2/test3] writing secret "test_test2_test3_secret/demo" 
[test/test2/test3] writing secret "test_test2_test3_secret/sub/demo" 
[test/test2/test3] restore engine: test_test2_test3_secret_2
[test/test2/test3] writing secret "test_test2_test3_secret_2/admin" 
[test/test2/test3] writing secret "test_test2_test3_secret_2/demo" 
[test/test2/test3] writing secret "test_test2_test3_secret_2/sub/demo" 
[test/test2/test3] writing secret "test_test2_test3_secret_2/sub/sub2/demo" 

# verify engines have been created
vkv list engines --all --include-ns-prefix
secret
secret_2
sub/sub2/sub_sub2_secret
sub/sub2/sub_sub2_secret_2
sub/sub_secret
sub/sub_secret_2
test/test2/test3/test_test2_test3_secret
test/test2/test3/test_test2_test3_secret_2
```


**Please note, that `vkv snapshot restore` is an experimental feature, you should always double check `vkv`s beahviour on a development Vault before running it against a production Vault**

