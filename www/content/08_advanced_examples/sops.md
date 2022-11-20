---
title: "encrypt & decrypt using sops"
weight: 3
---

In order to store the secret export created by `vkv` [sops](https://github.com/mozilla/sops#encrypting-using-hashicorp-vault) can be used.
This example shows how to encrypt & decrypt `vkv` exported secrets using `sops` and Vaults transit engine:

### Prerequisites
* Install [sops](https://github.com/mozilla/sops/releases)

### Demo
```bash
export VAULT_ADDR="https://vault.server"
export VAULT_TOKEN="hvs.XXXX"

# enable engine and create encryption key
vault secrets enable -path=sops transit
vault write sops/keys/vkv type=rsa-4096

# export secrets as yaml and write to file
vkv -p secret --show-values -f=yaml > export.yaml

# configure sops
cat <<EOF > .sops.yaml
creation_rules:
        - path_regex: \.yaml$
          hc_vault_transit_uri: "http://$VAULT_ADDR/v1/sops/keys/vkv"
EOF

# encrypt secrets
sops -e export.yaml > encrypted_export.yaml

# decrypt secrets
sops -d encrypted_export.yaml
```

an encrypted secrets file looks like this:


```yaml
secret/:
    admin:
        sub: ENC[AES256_GCM,data:fmHQMHCBNIs=,iv:s2q/j2tYvTN+u8KOXKm+Rbt1Y3oFO0fwjYCQy3jBHEU=,tag:+f5t/2LZCkIaAaoYeiY9KA==,type:str]
    demo:
        foo: ENC[AES256_GCM,data:kZwT,iv:QkNZwsUZ4lngluUHXae7abYAjAZFbNgJ7GdgM18GlLM=,tag:74aMFxt1pmCCy/rFcD7/rw==,type:str]
    sub/:
        demo:
            demo: ENC[AES256_GCM,data:f1m2veKk4w7i2hc=,iv:PyycH5Z9TEf/9u/nm7XYcnMEBHB+AY4ARAABoV8DQ74=,tag:QsPcwlw3vQNxytrfiZ1lyg==,type:str]
            password: ENC[AES256_GCM,data:m6DXfI4r,iv:3rLoWDTRfHuGUzJjoOemYv4C89EedK+CKX+9R7QfDZI=,tag:2Ul1wt77XGzina8QyOZMjQ==,type:str]
            user: ENC[AES256_GCM,data:EYXYIU8=,iv:7ll05h50Nu0Mp+bWuIrJjEsP4KRpH8L1vn3ZvqXlEPc=,tag:zN8Pzmk4QvX1hdlb955KaQ==,type:str]
        sub2/:
            demo:
                foo: ENC[AES256_GCM,data:AnMP,iv:oyaYacdlcnInw57im4ARprWz6wkgKqguiK6IHwdwn4w=,tag:4u0JZJ1jFvevZLDe/tlmzg==,type:str]
                password: ENC[AES256_GCM,data:Q1zzktiD58Q=,iv:1OZjqPyW0MNiTcll3tXkZ4AQ9CnNqtWzYmSw6PPEYxo=,tag:Z29GfxZnrWJnnzI8dckNHw==,type:str]
                user: ENC[AES256_GCM,data:a4Ju3Q==,iv:lpMo+/5K3mwqLQSpoAKLaL/Np6KAtoJDFZEAslu6TOY=,tag:BiZ4t7muKQu+AU9JvTNc4w==,type:str]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    hc_vault:
        - vault_address: http://127.0.0.1:8200
          engine_path: sops
          key_name: vkv
          created_at: "2022-11-17T14:23:40Z"
          enc: vault:v1:LtnJjUYl/pBDSOQrhSIsLp6XW0Ng/TM26GjBYcy95Fn8qAXBqJRyhYUd9Df5HF91RIhpiV11Rgj9hKj4sg0HcZIQnuBTQo30mgVQGcUhIL2PrV1qgDB3Ezm0W90s4aKH/8fhvToGVPB5nhf2/z9hTwnmFO+39GnC1JooofRdo9+1B7DBcsvliWSc4gIu3EbJwYUTqxLu92BJYWYM5ZZNtox0snJiYK3dfI7tltD7AtCYmwSJXFlkN3/lUrBGcNFOpya00/dKR/it2qbgtIJclcqPsYx6zoJ0MG/Z1RIs+v5mQjQlDKP2Stqxj+MijwtoPeXtdhdt7wdbWTxQHe4euw==
    age: []
    lastmodified: "2022-11-17T14:23:40Z"
    mac: ENC[AES256_GCM,data:X+yzr/+S6KWYaMrb1rvDK8bS5ghXDrnSzIT0RZ1TJdXzhcJmvxyUo2rfWkTVzNKcmoJnWyPKRujGVKZWTVmqa5hC5PQbc6tOqRgNrUmo4YkUoOyCyvPhJqCrRrjrhRtrUPIvRcDWMIhZCibjlr+XSglb+lkgVWNcghLWnez9I/w=,iv:2dZUWjnsLOLHV7wfFbIBDSW3ehLtsxmNeqmMwDsDqwU=,tag:PpaLsAA25TWQnYNTQnDswQ==,type:str]
    pgp: []
    unencrypted_suffix: _unencrypted
    version: 3.7.3
```