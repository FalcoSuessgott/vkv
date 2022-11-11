Output www/static/images/demo.gif

Set FontSize 32
Set Width 2000
Set Height 1000

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Type "export VAULT_TOKEN='root'"
Enter

Type "echo 'recursively list all secrets from KV-engine 'secret' in various output formats'"
Enter
Sleep 2s

Show
Type "vkv --path secret"
Enter
Sleep 3s

Type "export secrets and copy them to another secret engine using vkv import"
Enter
Sleep 2s

Type "vkv -p secret/sub --show-values --format=json | vkv import - -p copy_path"
Enter
Sleep 3s