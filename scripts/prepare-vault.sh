#!/usr/bin/bash

export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_SKIP_VERIFY="true"
export VAULT_TOKEN="root"

vault kv put secret/demo foo=bar
vault kv destroy -versions=1 secret/demo
vault kv put secret/admin sub=password
vault kv put secret/sub/demo demo="hello world" user=admin password='s3cre5<'
vault kv put secret/sub/sub2/demo foo=bar user=user password=password
vault kv patch secret/sub/sub2/demo new=test
vault kv patch secret/sub/sub2/demo user=admin
vault kv patch secret/sub/sub2/demo env=dev sub=test
vault kv put secret/sub/sub2/demo test=thisisaverylongsecretwhichshouldbetruncated

vault kv metadata put -mount=secret -custom-metadata=key=value admin
vault kv metadata put -mount=secret -custom-metadata=key=value -custom-metadata=admin=false sub/sub2/demo
