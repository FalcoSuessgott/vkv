Output assets/diff.gif

Set FontSize 30
Set Width 2000
Set Height 1000

Hide
Type "export VAULT_ADDR='http://127.0.0.1:8200'"
Enter

Hide
Type "export VAULT_TOKEN='root'"
Enter

Show
Type "diff -ty <(VKV_PATH='secret' vkv --show-values) <(VKV_PATH='secret_2' vkv --show-values)"
Enter
Sleep 5s