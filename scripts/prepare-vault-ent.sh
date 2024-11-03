#!/usr/bin/bash

export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_SKIP_VERIFY="true"
export VAULT_TOKEN="root"

vault secrets enable -path ${VAULT_NAMESPACE//\//_}_secret -version=2 kv
vault kv put ${VAULT_NAMESPACE//\//_}_secret/demo foo=bar
vault kv put ${VAULT_NAMESPACE//\//_}_secret/admin sub=password
vault kv put ${VAULT_NAMESPACE//\//_}_secret/sub/demo demo="hello world" user=admin password=s3cre5
vault kv put ${VAULT_NAMESPACE//\//_}_secret/sub/sub2/demo foo=bar user=user password=password
vault kv put ${VAULT_NAMESPACE//\//_}_secret/sub/sub2/demo admin=key foo=bar user=user password=password
vault kv metadata put -mount=${VAULT_NAMESPACE//\//_}_secret -custom-metadata=key=value admin
vault kv metadata put -mount=${VAULT_NAMESPACE//\//_}_secret -custom-metadata=key=value -custom-metadata=admin=false sub/sub2/demo
vault policy write kv assets/kv-policy.hcl

vault secrets enable -path ${VAULT_NAMESPACE//\//_}_secret_2 -version=2 kv
vault kv put ${VAULT_NAMESPACE//\//_}_secret_2/demo foo=bar
vault kv put ${VAULT_NAMESPACE//\//_}_secret_2/admin sub=password
vault kv put ${VAULT_NAMESPACE//\//_}_secret_2/sub/demo demo="hello world" user=admin password=s3cre5
vault kv put ${VAULT_NAMESPACE//\//_}_secret_2/sub/sub2/demo foo=bar-updated user=user password=password
