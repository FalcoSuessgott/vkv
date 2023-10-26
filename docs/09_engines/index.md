## vkv list engines

list all KVv2 engines

```
vkv list engines [flags]
```

### Options

```
  -n, --namespace string    specify the namespace (env: VKV_LIST_ENGINES_NS)
  -p, --include-ns-prefix   prepend the namespaces (env: VKV_LIST_ENGINES_NS_PREFIX)
  -r, --regex string        filter engines by the specified regex pattern (env: VKV_LIST_ENGINES_REGEX
  -a, --all                 list all KV engines recursively from the specified namespaces (env: VKV_LIST_ENGINES_ALL)
  -f, --format string       available output formats: "base", "json", "yaml" (env: VKV_LIST_ENGINES_FORMAT) (default "base")
  -h, --help                help for engines
```