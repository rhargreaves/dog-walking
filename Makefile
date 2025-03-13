export WORKDIR=$(shell pwd)

build:
	rm bootstrap api.zip
	docker compose build
	docker compose run --rm bootstrap
	cp api.zip infra/api.zip
.PHONY: build

test: build
	docker compose run --rm acceptance-test
.PHONY: test

sam-local-api: build
	docker compose run sam --rm --service-ports
.PHONY: sam-local-api

clean:
	-rm bootstrap api.zip infra/api.zip
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean