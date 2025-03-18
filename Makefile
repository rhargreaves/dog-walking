export PROJECT_ROOT=$(shell pwd)
ifeq ($(shell uname -s),Darwin)
    export CONTAINER_HOST=host.docker.internal
else
    export CONTAINER_HOST=172.17.0.1
endif

SHOW_LOGS_ON_FAILURE := false
GO_IMAGE := golang:1.23.4-alpine
TTY_ARG := $(shell if [ -t 0 ]; then echo "-t"; else echo ""; fi)
GO_CMD := docker run -i $(TTY_ARG) --rm \
	-v $(shell pwd):/app \
	-v go-cache:/go \
	-e GOPATH=/go \
	-e GOCACHE=/go/cache \
	-w /app \
	$(GO_IMAGE) \
	sh -ec

export AWS_REGION=eu-west-1
ifeq ($(ENV),uat)
	export DOG_IMAGES_BUCKET=uat-dog-images
	export API_BASE_URL=https://api.uat.dog-walking.roberthargreaves.com
	export DOGS_TABLE_NAME=uat-dogs
	export COGNITO_USER_POOL_NAME=uat-dog-walking
	export COGNITO_CLIENT_NAME=uat-dog-walking-client
endif

create-go-cache:
	-docker volume create go-cache
.PHONY: create-go-cache

build: create-go-cache lint test-unit compile
.PHONY: build

compile: create-go-cache lint test-unit
	docker compose down
	$(GO_CMD) "cd api; \
		rm -rf build; \
		mkdir build; \
		GOOS=linux GOARCH=arm64 go build -o build/bootstrap ./cmd/api"
.PHONY: compile

compile-local-auth:
	$(GO_CMD) "cd local-auth; \
		rm -rf build; \
		mkdir build; \
		go mod download; \
		go install gotest.tools/gotestsum@latest; \
		LOCAL_JWT_SECRET=1234567890 gotestsum --format testname ./...; \
		GOOS=linux GOARCH=arm64 go build -o build/bootstrap ."
.PHONY: compile-local-auth

lint:
	$(GO_CMD) "cd api && go fmt ./... && go mod tidy"
.PHONY: lint

test-unit: create-go-cache
	$(GO_CMD) "cd api; \
		go mod download; \
		go install github.com/vektra/mockery/v2@latest; \
		mockery; \
		go install gotest.tools/gotestsum@latest; \
		gotestsum --format testname ./..."
.PHONY: test-unit

test-local: build compile-local-auth
	docker compose build
	docker compose run --rm e2e-test-local \
		|| (if [ "$(SHOW_LOGS_ON_FAILURE)" = "true" ]; then \
			docker compose logs; \
		fi; \ exit 1)
	docker compose down
.PHONY: test-local

test:
	docker compose run --build --rm e2e-test
	docker compose down
.PHONY: test

start-local-api:
	docker compose up -d sam
.PHONY: start-local-api

stop-local-api:
	docker compose down
.PHONY: stop-local-api

test-e2e-local:
	docker compose run --rm e2e-test-local
.PHONY: test-e2e-local

clean-cache:
	-docker volume rm go-mod-cache
.PHONY: clean-cache

clean: clean-cache
	-rm -rf api/build
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean