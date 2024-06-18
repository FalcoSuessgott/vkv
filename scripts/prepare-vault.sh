#!/usr/bin/bash

export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_SKIP_VERIFY="true"
export VAULT_TOKEN="root"

vault secrets enable -path secret -version=2 kv
vault kv put secret/demo foo=bar
vault kv put secret/admin sub=password
vault kv put secret/sub/demo demo="hello world" user=admin password='s3cre5<'
vault kv put secret/sub/sub2/demo foo=bar user=user password=password
vault kv put secret/sub/sub2/demo foo=bar user=user password=password 
vault kv metadata put -mount=secret -custom-metadata=key=value admin
vault kv metadata put -mount=secret -custom-metadata=key=value -custom-metadata=admin=false sub/sub2/demo
vault policy write kv assets/kv-policy.hcl

vault secrets enable -path secret_2 -version=2 kv
vault kv put secret_2/demo foo=bar
vault kv put secret_2/admin sub=password
vault kv put secret_2/sub/demo demo="hello world" user=admin password='s3cre5<'
vault kv put secret_2/sub/sub2/demo foo=bar-updated user=user password=password


<<<<<<< HEAD
# # test cases
# 1. rootpath -> rootpath
# 2. subpath -> rootpath
# 3. subpath -> subpath
# 4. root path -> enginepath
# 5. enginepath -> enginepath
# 6. engie path + sub > engine path + sub
# 7. engine path -> root path
# 8. engine path -> subpath
=======
# test cases
1. rootpath -> rootpath
2. subpath -> rootpath
3. subpath -> subpath
4. root path -> enginepath
5. enginepath -> enginepath
6. engie path + sub > engine path + sub
7. engine path -> root path
8. engine path -> subpath
>>>>>>> 5e1e421 (feat(import): support engine path)
