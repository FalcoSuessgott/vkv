# Snapshots

`vkv` enables you to create and restore snapshots of all KVv2 engines in all namespaces of a Vault instance (requires an appropiate token + policy):

See the [CLI Reference](https://github.com/FalcoSuessgott/vkv/cmd/vkv_snapshort/) for more details on the supported flags and env vars.

## Example Usage
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
vkv snapshot save --destination vkv-export-$(date '+%Y-%m-%d')
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

As you can see: `vkv` exported all engines and wrote them to the specified directory:

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

## Restore vkv snapshots

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
