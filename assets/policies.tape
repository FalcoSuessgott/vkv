Output assets/policies.gif

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
Type "cat assets/policies.tmpl"
Enter

Type "vkv -p secret -f=template --template-file=assets/policies.tmpl"
Enter
Sleep 5
