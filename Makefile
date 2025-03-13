export WORKDIR=$(shell pwd)
ifeq ($(shell uname -s),Darwin)
    export CONTAINER_HOST=host.docker.internal
else
    export CONTAINER_HOST=172.17.0.1
endif

build:
	-rm bootstrap api.zip
	docker compose build
	docker compose run --rm bootstrap
	cp api.zip infra/api.zip
.PHONY: build

test-local: build
	docker compose run --rm acceptance-test-local \
		|| (docker compose logs && exit 1)
	docker compose down
.PHONY: test-local

test:
	docker compose run --build --rm acceptance-test
	docker compose down
.PHONY: test

sam-local-api: build
	docker compose run --rm --service-ports sam
.PHONY: sam-local-api

clean:
	-rm bootstrap api.zip infra/api.zip
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean