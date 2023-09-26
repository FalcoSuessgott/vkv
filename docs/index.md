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


### Features
* recursively print secrets of any KVv2 Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/05_export/formats/)
* engine export shows the secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/05_export/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: `policy`)
* print secrets in `export <key>=<value>` format for variable exporting (format: `export`)
* [import](https://falcosuessgott.github.io/vkv/06_import/) secrets back to Vault from `vkv`'s `json` or `yaml` format 
* save and restore KVv2 snapshots (including namespaces) and running on [kubernetes](https://falcosuessgott.github.io/vkv/09_advanced_examples/kubernetes/)
* list all engines or namespaces for scripting purposes
* handy [snippets](https://falcosuessgott.github.io/vkv/09_advanced_examples/) for managing KVv2 engines using `fzf`, `sops` & `diff`

Checkout the [Quickstart](https://falcosuessgott.github.io/vkv/01_quickstart) Guide to learn more about `vkv`
