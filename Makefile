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

.PHONY: fmtcheck
fmtcheck: ## run gofmt and print detected files
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

PHONY: test
test: ## run go tests
	go test -v ./...

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

PHONY: lint-fix
lint-fix: ## fix
	golangci-lint run -c .golang-ci.yml --fix

.PHONY: docker-build
docker-build: ## dockerize golang application
	@docker build --tag $(projectname) .

.PHONY: docker-run
docker-run:
	@docker run $(projectname)

.PHONY: pre-commit
pre-commit:	## run pre-commit hooks
	pre-commit run

