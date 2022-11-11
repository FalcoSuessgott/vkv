---
title: "Authentication"
weight: 2
---

`vkv` supports all of Vaults [environment variables](https://www.vaultproject.io/docs/commands#environment-variables). In order to authenticate you will have to set at least `VAULT_ADDR` and `VAULT_TOKEN`.

### MacOS/Linux
```
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="hvs.XXX" 
vkv -p <kv-path>
```

### Windows
```
SET VAULT_ADDR=http://127.0.0.1:8200
SET VAULT_TOKEN=s.XXX
vkv.exe -p <kv-path>
```

⚠️ **Your token policy requires `read` and `list` capabilities on every path of the specified secret engine, otherwise `vkv` errors with `403`. This behaviour is likely to change in future releases.**
