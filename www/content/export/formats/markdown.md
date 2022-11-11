---
title: "Markdown"
weight: 4
---

Display the secrets recursively as a markdown table

### Required flags

```bash
vkv --path <path> --format=markdown
```

### Optional flags
| Flag                  | Description                                                                       | Env                    | Default |
|-----------------------|-----------------------------------------------------------------------------------|------------------------|---------|
| `--only-keys`         | show only keys                                                                    | `VKV_ONLY_KEYS`        | `false` |
| `--only-paths`        | show only paths                                                                   | `VKV_ONLY_PATHS`       | `false` |
| `--show-values`       | don't mask values                                                                 | `VKV_SHOW_VALUES`      | `false` |
| `--max-value-length`  | maximum char length of values (set to `-1` for disabling)                         | `VKV_MAX_VALUE_LENGTH` | `12`    |

### Demo
<div align="center">
<br>
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/markdown.gif" alt="drawing" width="1000"/>
</div>
