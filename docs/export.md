# Export
`vkv export` requires an engine path (`--path` or `--engine-path`) and supports the following export formats (specify via `--format` flag).

See the [CLI Reference](https://falcosuessgott.github.io/vkv/cmd/vkv_export/) for more details on the supported flags and env vars.

!!! warning
    Vault allows `/` in the name of a KV engine. This makes it difficult for `vkv` to distinguish between directories and the KV engine name..

    If your KV engine name/mount contains a `/` you have to specify it using `--engine-path|-e`, otherwise `vkv` will output the secrets wrong.

    This also applies for any `vkv import ...` operations.

!!! info
    `vkv` handles 3 different path arguments, specified using `-e|-p`

    1. `root path`: any normal KV mount. Use `-p`.
    2. `engine-path`: in case your KV mount contains a `/`. Use `-e`.
    3. `sub path`: the path to the corresponding directory within a KV mount.
    When using `-p` this is everything after the first `/`: e.g: `kv/prod/db/`; root path=`kv`, subpath=`prod/db`.
    In conjunction with a `-e` you can specify a sub-path by using -p: `-e=kv/prod -p=db`.

## base
```bash
> vkv export -p secret -f=base
secret/ [desc=key/value secret storage] [type=kv2]
├── admin [v=1] [key=value]
│   └── sub=********
├── demo [v=1]
│   └── foo=***
└── sub
    ├── demo [v=1]
    │   ├── demo=***********
    │   ├── password=******
    │   └── user=*****
    └── sub2
        └── demo [v=2] [admin=false key=value]
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

## policy
```bash
> vkv export -p secret -f=policy
PATH                    CREATE  READ    UPDATE  DELETE  LIST    ROOT
secret/sub/sub2/demo    ✖       ✖       ✖       ✖       ✖       ✔
secret/admin            ✖       ✖       ✖       ✖       ✖       ✔
secret/demo             ✖       ✖       ✖       ✖       ✖       ✔
secret/sub/demo         ✖       ✖       ✖       ✖       ✖       ✔
```

## markdown
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

## template
`template` is a special output format that allows you, render the output using Golangs template engine. Format `template` requires either a `--template-file` or a `--template-string` flag or the equivalent env vars.

The secrets are passed as map with the secret path as the key and the actual secrets as values:

```
# <PATH>              <SECRETS>
secret/admin          map[sub:password]
secret/demo           map[foo:bar]
secret/sub/demo       map[demo:hello world password:s3cre5< user:admin]
secret/sub/sub2/demo  map[foo:bar password:password user:user]
```

Here is an advanced template that renders the secrets in a special env var export format. Note that within a `--template-file` or a `--template-string` the following functions are available: [http://masterminds.github.io/sprig/](http://masterminds.github.io/sprig/):

```jinja
# export.tmpl
{{- range $path, $secrets := . }}
{{- range $key, $value := $secrets }}
export {{ list $path $key | join "/" | replace "/" "_" | upper | trimPrefix "SECRET_" }}={{ $value | squote -}}
{{ end -}}
{{- end }}
```

This would result in the following output:

```bash
> vkv export -p secret -f=template --template-file=export.tmpl
export ADMIN_SUB='password'
export DEMO_FOO='bar'
export SUB_DEMO_DEMO='hello world'
export SUB_DEMO_PASSWORD='s3cre5<'
export SUB_DEMO_USER='admin'
export SUB_SUB2_DEMO_FOO='bar'
export SUB_SUB2_DEMO_PASSWORD='password'
export SUB_SUB2_DEMO_USER='user'
```

Per default `vkv` splits the secret paths at `/`, if you prefer a non-nested output (for scripting purposes) you can enable `--merge-paths` (only works in `yaml`, `json` or `template` output format):

```bash
# YAML
> vkv export -p secret --merge-paths -f=yaml
secret/admin:
  sub: password
secret/demo:
  foo: bar
secret/sub/demo:
  demo: hello world
  password: s3cre5<
  user: admin
secret/sub/sub2/demo:
  foo: bar
  password: password
  user: user

# JSON
> vkv export -p secret --merge-paths -f=json
{
  "secret/admin": {
    "sub": "password"
  },
  "secret/demo": {
    "foo": "bar"
  },
  "secret/sub/demo": {
    "demo": "hello world",
    "password": "s3cre5<",
    "user": "admin"
  },
  "secret/sub/sub2/demo": {
    "foo": "bar",
    "password": "password",
    "user": "user"
  }
}
```