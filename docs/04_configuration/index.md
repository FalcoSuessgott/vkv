### shell completion
`vkv` offers shell completion for `zsh`, `bash` and `fish` shells:

```bash
# bash
source <(vkv completion bash)

# systemwide
vkv completion bash > /etc/bash_completion.d/vkv

# zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
source <(vkv completion zsh); compdef _vkv vkv

# systemwide
vkv completion zsh > "${fpath[1]}/_vkv"

# fish
vkv completion fish | source

# systemwide
vkv completion fish > ~/.config/fish/completions/vkv.fish
```

### manpage
when installed via your systems package manager `vkv` ships manpages.

Simply run:

```bash
man vkv
```