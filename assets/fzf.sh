vault secrets list -format=json | jq -r 'to_entries | map(select(.value.type=="kv")) | from_entries | keys[]' | fzf --preview 'vkv -p ${}'
