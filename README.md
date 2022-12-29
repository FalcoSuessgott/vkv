<div align="center">

<div align="center">
<img src="images/logo.png" alt="drawing" width="200"/>
<br>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vkv" alt="drawing"/>

`vkv` is a little CLI tool written in Go, which enables you to list, compare, import, document, backup & encrypt secrets from a [HashiCorp Vault KV-v2 engine](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2):

<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/demo.gif" alt="drawing" width="1000" />

</div>

### Features
* recursively print secrets of any KVv2 Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/04_export/formats/)
* engine export shows the secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/04_export/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: `policy`)
* print secrets in `export <key>=<value>` format for variable exporting (format: `export`)
* [import](https://falcosuessgott.github.io/vkv/05_import/) secrets back to Vault from `vkv`'s `json` or `yaml` format 
* save and restore KVv2 snapshots (including namespaces)
* list all engines or namespaces for scripting purposes
* handy [snippets](https://falcosuessgott.github.io/vkv/08_advanced_examples/) for managing KVv2 engines using `fzf`, `sops` & `diff`



Checkout the [Quickstart](https://falcosuessgott.github.io/vkv/01_quickstart) Guide to learn more about `vkv`
