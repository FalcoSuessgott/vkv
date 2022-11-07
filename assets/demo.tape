Output assets/demo.gif

Set FontSize 32
Set Width 2000
Set Height 1000

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Hide
Type "export VAULT_TOKEN='root'"
Enter

Show
Type "vkv -p secret --show-values"
Enter
Sleep 3s

Type "vkv -p secret --show-values -f=policy"
Enter
Sleep 3s

Type "vkv -p secret --show-values -f=export"
Enter
Sleep 3s

Type "vkv -p secret --show-values -f=markdown"
Enter
Sleep 3s