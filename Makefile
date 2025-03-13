export WORKDIR=$(shell pwd)
ifeq ($(shell uname -s),Darwin)
    export CONTAINER_HOST=host.docker.internal
else
    export CONTAINER_HOST=172.17.0.1
endif

GO_IMAGE := golang:1.23.4-alpine
GO_CMD := docker run --rm -v $(shell pwd):/app -w /app $(GO_IMAGE) sh -c

build: lint
	-rm -rf out
	mkdir -p out
	$(GO_CMD) "go mod download && \
		GOOS=linux GOARCH=arm64 go build -o out/bootstrap main.go"
.PHONY: build

lint:
	$(GO_CMD) "go fmt ./... && go mod tidy"
.PHONY: lint

test-local: build
	docker compose build
	docker compose run --rm acceptance-test-local \
		|| (docker compose logs && exit 1)
	docker compose down
.PHONY: test-local

test:
	docker compose run --build --rm acceptance-test
	docker compose down
.PHONY: test

start-local-api: build
	docker compose run --rm --service-ports sam
.PHONY: start-local-api

clean:
	-rm -rf out
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean