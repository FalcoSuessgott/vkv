---
title: browse all KVv2 engines
weight: 2
---

using `vault secrets list` and a little bit of `jq`-logic, we can get a list of all KV-engines visible for the token. 

If we pipe this into `fzf` we can get a handy little  preview-app:

```bash
vkv list engines --all --include-ns-prefix | fzf --preview 'vkv export -e ${}'
```

### Demo
<div align="center">
<br>
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/fzf.gif" alt="drawing" width="1000"/>
</div>