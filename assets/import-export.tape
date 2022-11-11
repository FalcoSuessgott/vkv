Output www/static/images/export-import.gif

Set FontSize 32
Set Width 2000
Set Height 1000

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Type "export VAULT_TOKEN='root'"
Enter

Show
Type "vkv -p secret/sub"
Enter
Sleep 2s

Type "vkv -p secret/sub --show-values -f=yaml | vkv import - -p new_path"
Enter
Sleep 3s

Type "vkv -p secret/admin --show-values -f=yaml | vkv import - -p secret/dev"
Enter
Sleep 3s
