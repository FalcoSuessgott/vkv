# fzf

using `vault secrets list` and a little bit of `jq`-logic, we can get a list of all KV-engines visible for the token. 

If we pipe this into `fzf` we can get a handy little  preview-app:

```bash
vkv list engines --all --include-ns-prefix | fzf --preview 'vkv export -e ${}'
```

## Demo
![gif](assets/fzf.gif)
