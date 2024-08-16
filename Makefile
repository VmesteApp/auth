include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

run-app: ### Run app (after `make compose-up`)
	go run -tags migrate ./cmd/app
.PHONY: run-app

migrate-create:  ### create new migration
	./bin/migrate create -ext sql -dir migrations $(name)
.PHONY: migrate-create

migrate-up: ### migration up
	./bin/migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

linter-check: ### check by golangci linter
	./bin/golangci-lint run
.PHONY: linter-golangci

linter-fix: ### fix by golangci linter
	./bin/golangci-lint run --fix
.PHONY: linter-golangci

bin-deps: ### install deps
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: bin-deps