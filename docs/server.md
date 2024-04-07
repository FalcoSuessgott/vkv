# Server
`vkv server` is a subcommand that starts simple http server that accepts `GET` request `/export` on port `8080` (change using `--port`).

## Options

```
  -P, --port string          HTTP Server Port (env: VKV_SERVER_PORT) (default "8080")
  -p, --path string          KVv2 Engine path (env: VKV_SERVER_PATH)
  -e, --engine-path string   engine path in case your KV-engine contains special characters such as "/", the path value will then be appended if specified ("<engine-path>/<path>") (env: VKV_SERVER_ENGINE_PATH)
      --skip-errors          dont exit on errors (permission denied, deleted secrets) (env: VKV_SERVER_SKIP_ERRORS)
  -h, --help                 help for server
```

This is helps using `vkv` as a service container for usage during CI:


## Server side
```bash
export VAULT_ADDR="..."
export VAULT_TOKEN="..." 
vkv server --path secret
```

## Client side
```bash
$> curl localhost:88080/export
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

## Output Format
you can specify the output format by adding a `format`-URL Query Parameter:

```bash
$> curl localhost:88080/export?format=yaml
secret/:
  admin:
    sub: '********'
  demo:
    foo: '***'
  sub/:
    demo:
      demo: '***********'
      password: '******'
      user: '*****'
    sub2/:
      demo:
        admin: '***'
        foo: '***'
        password: '********'
        user: '****'
```
