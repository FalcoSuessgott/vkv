Output www/static/images/demo.gif

Set FontSize 32
Set Width 2000
Set Height 1000

Show
Type "vkv export --path secret"
Enter
Sleep 3s

Type "vkv export -p secret --show-values --format=json | vkv import - -p secret_copy"
Enter
Sleep 3s

Type "vkv list engines --all --include-ns-prefix | fzf --preview 'vkv export -e ${}'"
Enter
Sleep 3

Up
Sleep 2

Down
Sleep 2