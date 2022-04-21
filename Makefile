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

PHONY: clean
clean: ## clean up environment
	@rm -rf coverage.out dist/ $(projectname)

PHONY: cover
cover: ## display test coverage
	go test -v -race $(shell go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

PHONY: fmt
fmt: ## format go files
	gofumpt -w -s  .

PHONY: lint
lint: ## lint go files
	golangci-lint run -c .golang-ci.yml

.PHONY: pre-commit
pre-commit:	## run pre-commit hooks
	pre-commit run

vault: export VAULT_ADDR = http://127.0.0.1:8200
vault: export VAULT_SKIP_VERIFY = true
vault: export VAULT_TOKEN = root

.PHONY: vault
vault:
	echo hallo
	nohup vault server -dev -dev-root-token-id=root 2> /dev/null &
	vault kv put secret/demo foo=bar
	vault kv put secret/sub sub=password
	vault kv put secret/sub/demo foo=bar user=user password=password
	vault kv put secret/sub/sub2/demo foo=bar user=user password=password

	vault secrets enable -path secret_2 -version=2 kv
	vault kv put secret_2/demo foo=bar
	vault kv put secret_2/sub sub=password
	vault kv put secret_2/sub/demo foo=bar user=user password=password
	vault kv put secret_2/sub/sub2/demo foo=bar user=user password=password

.PHONY: kill
kill:
	@kill -9 $(shell pgrep -x vault) 2> /dev/null || true
