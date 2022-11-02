Output assets/export.gif

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
Type "vkv -p secret/sub/demo -f=export"
Enter
Sleep 3

Type "eval $(vkv -p secret/sub/demo -f=export)"
Enter
Sleep 1

Type "echo $password"
Enter
Sleep 3