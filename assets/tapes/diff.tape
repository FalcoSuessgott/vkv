Output www/static/images/diff.gif

Set FontSize 30
Set Width 2000
Set Height 1000

Show
Type "diff -ty <(vkv export -p secret --show-values) <(vkv export -p secret_2 --show-values)"
Enter
Sleep 5s
