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
	go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | sort -rnk3

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

.PHONY: vault
vault: clean ## set up a development vault server and write kv secrets
	nohup vault server -dev -dev-root-token-id=root 2> /dev/null &
	sleep 3

	./scripts/prepare-vault.sh

.PHONY: vault-ent
vault-ent: clean ## set up a development vault enterprise server and write kv secrets
	nohup vault server -dev -dev-root-token-id=root 2> /dev/null &
	sleep 3

	./scripts/prepare-vault.sh

	vault namespace create sub
	VAULT_NAMESPACE="sub" ./scripts/prepare-vault-ent.sh

	VAULT_NAMESPACE="sub" vault namespace create sub2
	VAULT_NAMESPACE="sub/sub2" ./scripts/prepare-vault-ent.sh

	vault namespace create test
	VAULT_NAMESPACE="test" vault namespace create test2
	VAULT_NAMESPACE="test/test2" vault namespace create test3
	VAULT_NAMESPACE="test/test2/test3" ./scripts/prepare-vault-ent.sh

.PHONY: clean
clean: ## clean the development vault
	@rm -rf coverage.out dist/ $(projectname) manpages/ dist/ completions/ || true
	@kill -9 $(shell pgrep -x vault) 2> /dev/null || true
	@kill -9 $(shell pgrep -x vault-ent) 2> /dev/null || true

ASSETS = diff demo fzf
.PHONY: assets
assets: clean vault-ent ## generate all assets
	for i in $(ASSETS); do \
		vhs < assets/tapes/$$i.tape; \
	done

.PHONY: www
www: ## build and server docs
	hugo server -s www

.PHONY: docker/build
docker/build: ## build docker images
	docker build \
		--pull \
		--no-cache \
		-t vkv:latest \
		.

.PHONY: docker/run
docker/run: ## run docker image
	@docker run \
		--network=host \
		-e VAULT_ADDR=${VAULT_ADDR} \
		-e VAULT_TOKEN=${VAULT_TOKEN} \
		vkv export -p secret
