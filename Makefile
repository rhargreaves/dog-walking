export HOST_API_PATH=$(shell pwd)/api
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
	-v go-mod-cache:/gomodcache \
	-e GOMODCACHE=/gomodcache \
	-w /app \
	$(GO_IMAGE) \
	sh -ec

export AWS_REGION=eu-west-1
ifeq ($(ENV),uat)
	export DOG_IMAGES_BUCKET=uat-dog-images
	export API_BASE_URL=https://api.uat.dog-walking.roberthargreaves.com
	export DOGS_TABLE_NAME=uat-dogs
endif

create-mod-cache:
	-docker volume create go-mod-cache
.PHONY: create-mod-cache

build: create-mod-cache lint test-unit compile
.PHONY: build

compile: create-mod-cache lint test-unit
	docker compose down
	$(GO_CMD) "cd api; \
		rm -rf build; \
		mkdir build; \
		GOOS=linux GOARCH=arm64 go build -o build/bootstrap ./cmd/api"
.PHONY: compile

lint:
	$(GO_CMD) "cd api && go fmt ./... && go mod tidy"
.PHONY: lint

test-unit: create-mod-cache
	$(GO_CMD) "cd api; \
		go mod download; \
		go install github.com/vektra/mockery/v2@latest; \
		go generate ./...; \
		go install gotest.tools/gotestsum@latest; \
		gotestsum --format testname ./..."
.PHONY: test-unit

test-local: build
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

start-local-api: build
	docker compose run --rm --service-ports sam
.PHONY: start-local-api

clean-cache:
	-docker volume rm go-mod-cache
.PHONY: clean-cache

clean: clean-cache
	-rm -rf api/build
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean