`vkv` and  all its subcommands are highly configurable using environment variables.

Checkout the subcommands documentation aswell as the help messages (`vkv <subcommand> --help`) to see the available environment variables.

## Mode
You can control the executed subcommand of `vkv` by setting `VKV_MODE` to either on of:

* `export`
* `import`
* `server`
* `list`
* `snapshot_restore`
* `snapshot_save`

example:

```bash
VKV_EXPORT_PATH=secret VKV_MODE=export vkv
secret/
├── v1: admin [key=value]   
│   └── sub=********        
├── v1: demo
│   └── foo=***
└── sub/
    ├── v1: demo
    │   ├── demo=***********
    │   ├── password=******
    │   └── user=*****
    └── sub2
        └── v2: demo [admin=false key=value]
            ├── admin=***
            ├── foo=***
            ├── password=********
            └── user=****
```

## Shell Completion
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

## Manpage
when installed via your systems package manager `vkv` ships manpages.

Simply run:

```bash
man vkv
```