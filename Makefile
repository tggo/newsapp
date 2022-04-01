APP_NAME := booster-news
BUILD_NUMBER = -1 ## rewrite in ci/cd
BIN ?= ./bin
PKG_LIST ?= $(shell go list ./... | grep -v /vendor/ | grep -v /tmp/ )

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG := $(shell git describe --exact-match --abbrev=0 --tags 2> /dev/null)
GIT_HASH := $(shell git log --format="%h" -n 1)
GIT_DIRTY := $(shell git diff-index --quiet HEAD -- || echo "-dirty")
GIT_COMMIT_HASH := $(shell git rev-parse @)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
APP_VERSION ?= $(if $(GIT_TAG),$(GIT_TAG)-$(GIT_HASH)$(GIT_DIRTY),$(GIT_BRANCH)$(GIT_DIRTY))

DOCKER_BIN = $(shell command -v docker 2> /dev/null)
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)"


LDFLAGS := -ldflags " -X main.appName=$(APP_NAME) -X main.release=$(APP_VERSION)  -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(BUILD_DATE) -X main.buildNumber=$(BUILD_NUMBER)"
CI_COMMIT:=$(GIT_COMMIT_HASH)

__dummy := $(shell touch .env)
include .env
export

.PHONY : help build run fmt lint gotest test test_ci cover shell redis-cli image clean
.DEFAULT_GOAL : help
.SILENT : test shell

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

image: ## Build docker image with app
	$(DOCKER_BIN) build -f .cicd/build/Dockerfile -t $(APP_NAME):local .
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Now you can use image like `docker run --rm $(APP_NAME):local ...`';

build: ## Build app binary file
	$(DC_BIN) build

run: build ## Start app
	$(DC_BIN) up

shell: ## Start shell into container with golang
	$(DC_BIN) run $(DC_RUN_ARGS) app sh

clean: ## Make clean
	$(DC_BIN) down -v -t 1
	$(DOCKER_BIN) rmi $(APP_NAME):local -f
	@echo "-> stoping postgres"
	$(DOCKER_BIN) rm -f boostersnews-postgres || true

local_run:  local_build ## Run server on host
	$(BIN)/$(APP_NAME)

local_build: local_build_server ## Build binaries on host
	@echo "end build binaries"

local_build_server:
	@echo "start build server"
	@CGO_ENABLED=0 go build  -a -ldflags="-w -s" -trimpath -o $(BIN)/$(APP_NAME) $(LDFLAGS) ./cmd/server

lint: ## Run GoLangCI Lint
	golangci-lint run ./...

test: ## Run unit tests
	@go test -cover -count=1 -timeout=10s -short ${PKG_LIST}

test_ci:
	go get -u github.com/jstemmer/go-junit-report
	go test -count=1 -timeout=160s -short -v ${PKG_LIST} 2>&1 | $(GOPATH)/bin/go-junit-report -set-exit-code > report.xml
	go test -count=1  ${PKG_LIST} -short -v -coverprofile .testCoverage.txt

generate_mocks: ## generate specific mocks via go-mockery
	@echo "use go-mockery: brew install mockery"
	mockery \
		--name Service --dir ./internal/news/service   --output ./internal/news/mocks --case snake
	mockery \
		--name Repository --dir ./internal/news   --output ./internal/news/mocks --case snake

generate_swagger_api:
	## via docker
	$(DOCKER_BIN) run --rm -v "${PWD}/api:/api" openapitools/openapi-generator-cli generate \
        -i /api/swagger.yaml \
        -g go-gin-server \
        -o /api/gen
	## via local installed
	#openapi-generator generate -g go-gin-server -i ./api/swagger.yaml -o ./gen
	find $(shell pwd)/api/gen/go -type f -name "model_post*.go" | sed -e "p;s/model_/response_/" | xargs -n2 cp
	find $(shell pwd)/api/gen/go -type f -name "model_one*.go" | sed -e "p;s/model_/response_/" | xargs -n2 cp
	find $(shell pwd)/api/gen/go -type f -name "model_list*.go" | sed -e "p;s/model_/response_/" | xargs -n2 cp
	find $(shell pwd)/api/gen/go -type f -name "model_bad*.go" | sed -e "p;s/model_/response_/" | xargs -n2 cp
	find $(shell pwd)/api/gen/go -type f -name "model_success*.go" | sed -e "p;s/model_/response_/" | xargs -n2 cp
	rm -rf ./api/openapi
	mkdir -p ./api/openapi
	mv ./api/gen/go/response_* ./api/openapi
	rm -rf ./api/gen

docker_postgres: ## Start docker postgres for local development
	$(DOCKER_BIN) run -p 5441:5432  --name boostersnews-postgres -e POSTGRES_DB=boostersnews -e POSTGRES_USER=boosterdev -e POSTGRES_PASSWORD=mysecretpassword -d postgres:12

docker_postgres_stop: ## Start docker postgres for local development
	$(DOCKER_BIN) stop boostersnews-postgres

migrate: ## get DB migrate version
	@echo "database url get from env variable, store version on 'schema_migrations'"
	@echo "documentation https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
	@echo "create migrate via cmd: migrate create -ext sql -dir .cicd/deploy/migrate -seq create_posts_table"
	@migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" version
	@#migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" up
#	@migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" down

migrate_up: ## run migration to last version
	@echo "database url get from env variable, store version on 'schema_migrations'"
	@migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" up
#	@migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" down

migrate_down: ## down migration version
	@echo "database url get from env variable, store version on 'schema_migrations'"
	@migrate -source="file://.cicd/deploy/migrate" -database="$(POSTGRES_URL)" down

