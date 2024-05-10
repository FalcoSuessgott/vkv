# Development
Clone this repository and run:

```sh
make bootstrap
```

in order to have all used build dependencies

You can spin up a development vault for local testing by running:

```sh
make vault
```

The following environment variables are required:

```sh
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="root"
export VKV_PATH="secret"
```

If everything worked fine, you should be able to run:

```sh
go run main.go export -p secret
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