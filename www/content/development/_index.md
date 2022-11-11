---
title: "Development"
weight: 40
---

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
go run main.go   
secret/
├── demo
│   └── foo=***
├── sub
│   └── sub=********
├── sub/
│   └── demo
│       ├── demo=***********
│       ├── password=******
│       └── user=*****
└── sub/
    └── sub2/
        └── demo
            ├── password=*******
            ├── user=********
            └── value=*********
```