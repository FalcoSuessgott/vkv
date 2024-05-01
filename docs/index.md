# vkv 

<div align="center">
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vkv" alt="drawing"/>
</div>

![gif](assets/demo.gif)


`vkv` is a little CLI tool written in Go, which enables you to list, compare, import, document, backup & encrypt secrets from a [HashiCorp Vault KV engine](https://developer.hashicorp.com/vault/docs/secrets/kv):

## Features
* Support KV version 1 & version 2 (no need to specify the version `vkv` will automatically detect the engines version)
* **CI/CD Integrations for [Gitlab, GitHub, Azure Devops](https://falcosuessgott.github.io/vkv/cicd_gitlab)**
* support all Vault Auth Env Vars and `VKV_LOGIN_COMMAND` for avoiding having to hardcode the `VAULT_TOKEN` ([example](https://falcosuessgott.github.io/vkv/authentication/))
* recursively print secrets of any KV  Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/export/)
* engine export shows the secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/export/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: `policy`)
* print secrets in `export <key>=<value>` format for env var exporting (format: `export`)
* move or migrate secrets from KVV1 to a KVV2 Engine or any subpath [example](https://falcosuessgott.github.io/vkv/example_migrate_kvv1_kvv2/)
* [import](https://falcosuessgott.github.io/vkv/import/) secrets back to Vault from `vkv`'s `json` or `yaml` format output
* save and restore KVv2 snapshots (including namespaces) ([kubernetes](https://falcosuessgott.github.io/vkv/example_kubernetes/) example)
* list all KVv2-engines or namespaces for scripting purposes ([fzf](https://falcosuessgott.github.io/vkv/example_fzf/) example)
* more handy [snippets](https://falcosuessgott.github.io/vkv/example_diff/) using `fzf`, `sops` & `diff`

**Checkout the [Quickstart](https://falcosuessgott.github.io/vkv/quickstart/) Guide to learn more about `vkv` as well as the [CLI Reference](https://falcosuessgott.github.io/vkv/cmd/vkv/)**