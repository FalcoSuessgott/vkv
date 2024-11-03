# Server
`vkv server` starts a **unauthenticated http server on port `127.0.0.01:8080` per default, that returns the kv secrets. This is useful during CI/CD setups.

See the [CLI Reference](https://github.com/FalcoSuessgott/vkv/cmd/vkv_server/) for more details on the supported flags and env vars.

## Server side
```bash
export VAULT_ADDR="..."
export VAULT_TOKEN="..."
> vkv server --path secret
listening on 127.0.0.1:8080
```

## Client side
```bash
$> curl localhost:8080/export
export admin='key'
export demo='hello world'
export foo='bar'
export password='password'
export sub='password'
export user='user'
```

## Output Format
you can specify the output format by adding a `format`-URL Query Parameter:

```bash
$> curl localhost:8080/export?format=yaml
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
