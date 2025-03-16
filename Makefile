export HOST_API_PATH=$(shell pwd)/api
ifeq ($(shell uname -s),Darwin)
    export CONTAINER_HOST=host.docker.internal
else
    export CONTAINER_HOST=172.17.0.1
endif

SHOW_LOGS_ON_FAILURE := false
GO_IMAGE := golang:1.23.4-alpine
GO_CMD := docker run --rm \
	-v $(shell pwd):/app \
	-v go-mod-cache:/gomodcache \
	-e GOMODCACHE=/gomodcache \
	-w /app \
	$(GO_IMAGE) \
	sh -ec

create-mod-cache:
	-docker volume create go-mod-cache
.PHONY: create-mod-cache

build: create-mod-cache lint
	docker compose down
	$(GO_CMD) "cd api; \
		rm -rf build; \
		mkdir build; \
		go mod download; \
		GOOS=linux GOARCH=arm64 go build -o build/bootstrap ./cmd/api"
.PHONY: build

lint:
	$(GO_CMD) "cd api && go fmt ./... && go mod tidy"
.PHONY: lint

test-local: build
	docker compose build
	docker compose run --rm acceptance-test-local \
		|| (if [ "$(SHOW_LOGS_ON_FAILURE)" = "true" ]; then \
			docker compose logs; \
		fi; \ exit 1)
	docker compose down
.PHONY: test-local

test:
	docker compose run --build --rm acceptance-test
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