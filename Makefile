build-images:
	docker compose build
.PHONY: build-images

build: build-images
	docker compose run --rm bootstrap
	cp api.zip infra/api.zip
.PHONY: build

test: build-images
	docker compose run --rm acceptance-test
.PHONY: test

sam-local-api: build
	docker compose run --rm \
		--service-ports \
		sam local start-api \
			--host 0.0.0.0 \
			--container-host host.docker.internal \
			--docker-volume-basedir $(shell pwd)
.PHONY: sam-local-api

clean:
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean