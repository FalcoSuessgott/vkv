
Output assets/template.gif

Set FontSize 32
Set Width 1500
Set Height 1200

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Hide
Type "export VAULT_TOKEN='root'"
Enter

Show
Type "cat assets/template.tmpl"
Enter

Type "vkv -f=template --template-file=assets/template.tmpl --show-values"
Enter
Sleep 5
