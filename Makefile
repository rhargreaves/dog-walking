build-images:
	docker compose build
.PHONY: build-images

build: build-images
	docker compose run bootstrap
	cp api.zip infra/api.zip
.PHONY: build

test: build-images
	docker compose run acceptance-test
	docker compose down --remove-orphans
.PHONY: test

clean:
	docker compose down --rmi all --volumes --remove-orphans
.PHONY: clean