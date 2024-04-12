# Migrate Secrets from KVv1 to KVv2
`vkv` enables you to quickly migrate KVv1 secrets KVv2:

```bash
# list all secret engines
$> vkv list engines
kvv1/

# list kvv1 secrets
$>  vkv export -p kvv1
kvv1/
└── dev
    ├── admin=****

# move secrets to kvv2 engine
$> vkv export -p kvv1 -f=json | vkv import -p kv
v2
reading secrets from STDIN
parsing secrets from JSON
writing secret "kvv2/dev" 
successfully imported all secrets

result:

kvv2/
└── v1: dev
    ├── admin=****

# verify
$> vkv export -p kvv2 --show-values
kvv2/
└── v1: dev
    ├── admin=user
    └── password=ok
```

You can also move a KV mount within another engine:

```bash
$> vkv export -p kvv1 -f=json | vkv import -p engine/subpath --force 
reading secrets from STDIN
parsing secrets from JSON
writing secret "engine/subpath/dev" 
successfully imported all secrets

result:

engine/subpath/
└── subpath/
    └── dev
        ├── admin=****

# verify
$> vkv export -p engine
engine/
└── subpath/
    └── v1: dev
        ├── admin=****
        └── password=**
```