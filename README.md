<div align="center">

`vkv` is a little CLI tool written in Go, which enables you to list, compare, import, document, backup & encrypt secrets from a [HashiCorp Vault KV-v2 engine](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2):

<img src="www/static/images/demo.gif" alt="drawing" width="1000" />
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vkv/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vkv/branch/master/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vkv/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vkv" alt="drawing"/>

</div>

### Features
* recursively print secrets of any KVv2 Engine in `json`, `yaml`, `markdown` and [other formats](https://falcosuessgott.github.io/vkv/export/formats/)
* show a secret version as well as its [custom metadata](https://developer.hashicorp.com/vault/docs/commands/kv/metadata)
* customize the output (show only-keys, only-paths, mask/unmask secrets) via [flags or environment](https://falcosuessgott.github.io/vkv/export/usage/)
* print the CRUD-capabilities of the authenticated token for each KV-path (format: [`policy`](https://falcosuessgott.github.io/vkv/export/formats/token_policy/))
* print secrets in `export <key>=<value>` format for variable exporting (format: [`export`](https://falcosuessgott.github.io/vkv/export/formats/export/))
* [import](https://falcosuessgott.github.io/vkv/import/usage/) secrets back to Vault from `vkv`'s `json` or `yaml` format 
* handy [snippets](https://falcosuessgott.github.io/vkv/export/advanced_examples/) for managing KVv2 engines using `fzf`, `sops` & `diff`


Checkout the [Docs](https://falcosuessgott.github.io/vkv/) to learn more about `vkv`
