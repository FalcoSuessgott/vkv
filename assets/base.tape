Output assets/base.gif

Set FontSize 36
Set Width 1100
Set Height 1000

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Hide
Type "export VAULT_TOKEN='root'"
Enter

Show
Type@200ms "vkv -p secret --show-values"
Enter
Sleep 5s