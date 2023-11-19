<div align="center">
<h1> vkv </h1>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vkv" alt="drawing"/>
</div>

![gif](assets/demo.gif)


`vkv` is a little CLI tool written in Go, which enables you to list, compare, import, document, backup & encrypt secrets from a [HashiCorp Vault KV-v2 engine](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2):


## Features
* **CI/CD Integrations for [Gitlab, GitHub, Azure Devops](https://falcosuessgott.github.io/vkv/cicd/gitlab/)**
* support all Vault Auth Env Vars and `VKV_LOGIN_COMMAND` for avoiding having to hardcode the `VAULT_TOKEN` ([example](https://falcosuessgott.github.io/vkv/authentication/))
* recursively print secrets of any KVv2 Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/export/formats/)
* engine export shows the secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/export/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: `policy`)
* print secrets in `export <key>=<value>` format for env var exporting (format: `export`)
* [import](https://falcosuessgott.github.io/vkv/06_import/) secrets back to Vault from `vkv`'s `json` or `yaml` format output
* save and restore KVv2 snapshots (including namespaces) ([kubernetes](https://falcosuessgott.github.io/vkv/advanced_examples/kubernetes/) example)
* list all KVv2-engines or namespaces for scripting purposes ([fzf](https://falcosuessgott.github.io/vkv/advanced_examples/fzf/) example)
* more handy [snippets](https://falcosuessgott.github.io/vkv/advanced_examples/diff/) using `fzf`, `sops` & `diff`,

Checkout the [Quickstart](https://falcosuessgott.github.io/vkv/quickstart/) Guide to learn more about `vkv`
