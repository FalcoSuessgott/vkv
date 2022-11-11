---
title: "Template"
weight: 8
---

Template the secrets using [Go-Template-Syntax](https://pkg.go.dev/text/template).

When using the template output format, all the data is passed to STDOUT as a

```go
map[string][]entry
```

where entry is a struct of

```go
type entry struct {
  Key   string
  Value interface{}
}
```

Also see [Generate Vault Policies from a KV engine](https://falcosuessgott.github.io/vkv/export/advanced_examples/vault_policies/).

### Required flags

```bash
vkv --path <path> --format=template
```


### Optional flags
| Flag                  | Description                                                                       | Env                    | Default |
|-----------------------|-----------------------------------------------------------------------------------|------------------------|---------|
| `--template-file`     | path to a file containing Go-template syntax to render the KV entries             | `VKV_TEMPLATE_FILE`    |         |
| `--template-string`   | string containing Go-template syntax to render KV entries                         | `VKV_TEMPLATE_STRING`  |         |

### Demo
<div align="center">
<br>
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/template.gif" alt="drawing" width="1000"/>
</div>