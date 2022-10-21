projectname?=vkv

default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build golang binary
	@go build -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags)" -o $(projectname)

.PHONY: install
install: ## install golang binary
	@go install -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags)"

.PHONY: run
run: ## run the app
	@go run -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags)"  main.go

PHONY: test
test: clean ## display test coverage
	go test -json -v ./... | gotestfmt

PHONY: fmt
fmt: ## format go files
	gofumpt -w .
	gci write .

PHONY: lint
lint: ## lint go files
	golangci-lint run -c .golang-ci.yml

.PHONY: pre-commit
pre-commit:	## run pre-commit hooks
	pre-commit run

.PHONY: bootstrap
bootstrap: ## install build deps
	go generate -tags tools tools/tools.go

vault: export VAULT_ADDR = http://127.0.0.1:8200
vault: export VAULT_SKIP_VERIFY = true
vault: export VAULT_TOKEN = root

.PHONY: vault
vault: clean ## set up a development vault server and write kv secrets
	nohup vault server -dev -dev-root-token-id=root 2> /dev/null &
	sleep 5
	vault kv put secret/demo foo=bar
	vault kv put secret/sub sub=password
	vault kv put secret/sub/demo1 demo="hello world" user=admin password=s3cre5
	vault kv put secret/sub/sub2/demo value="nevermind" user="database" password=secret2

	vault secrets enable -path secret_2 -version=2 kv
	vault kv put secret_2/demo foo=bar
	vault kv put secret_2/sub sub=password
	vault kv put secret_2/sub/demo foo=bar user=user password=password
	vault kv put secret_2/sub/sub2/demo foo=bar user=user password=password

.PHONY: clean
clean: ## clean the development vault
	@rm -rf coverage.out dist/ $(projectname)
	@kill -9 $(shell pgrep -x vault) 2> /dev/null || true