# server
`vkv server` is a subcommand that starts simple http server that accepts `GET` request `/export` on port `8080` (change using `--port`).

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
you can speciy the output format by adding a `format`-URL Query Parameter:

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
