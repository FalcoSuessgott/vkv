# Export
`vkv export` requires an engine path (`--path` or `--engine-path`) and supports the following export formats (specify via `--format` flag). 

See the [CLI Reference](https://github.com/FalcoSuessgott/vkv/cmd/vkv_export/) for more details on the supported flags and env vars.

## base
```bash
> vkv export -p secret -f=base               
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

## yaml
```bash
> vkv export -p secret -f=yaml                       
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

## json
```bash
> vkv export -p secret -f=json
{
  "secret/": {
    "admin": {
      "sub": "********"
    },
    "demo": {
      "foo": "***"
    },
    "sub/": {
      "demo": {
        "demo": "***********",
        "password": "******",
        "user": "*****"
      },
      "sub2/": {
        "demo": {
          "admin": "***",
          "foo": "***",
          "password": "********",
          "user": "****"
        }
      }
    }
  }
}
```

## export
```bash
> vkv export -p secret -f=export
export admin='key'
export demo='hello world'
export foo='bar'
export password='password'
export sub='password'

eval $(vkv export -p secret -f=export)
echo $admin
key
```

### policy
```bash
> vkv export -p secret -f=policy 
PATH                    CREATE  READ    UPDATE  DELETE  LIST    ROOT
secret/sub/sub2/demo    ✖       ✖       ✖       ✖       ✖       ✔
secret/admin            ✖       ✖       ✖       ✖       ✖       ✔
secret/demo             ✖       ✖       ✖       ✖       ✖       ✔
secret/sub/demo         ✖       ✖       ✖       ✖       ✖       ✔
```

### markdown
```bash
> vkv export -p secret -f=markdown      
|         PATH         |   KEY    |    VALUE    | VERSION |       METADATA        |
|----------------------|----------|-------------|---------|-----------------------|
| secret/admin         | sub      | ********    |       1 | key=value             |
| secret/demo          | foo      | ***         |       1 |                       |
| secret/sub/demo      | demo     | *********** |       1 |                       |
|                      | password | ******      |         |                       |
|                      | user     | *****       |         |                       |
| secret/sub/sub2/demo | admin    | ***         |       2 | admin=false key=value |
|                      | foo      | ***         |         |                       |
|                      | password | ********    |         |                       |
|                      | user     | ****        |         |                       |
```
