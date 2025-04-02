TERRAFORM_IMAGE := hashicorp/terraform:latest
TERRAFORM_WORKDIR := /workspace

ifneq ($(wildcard ../.env),)
include ../.env
endif

ifeq ($(ENV),)
$(error ENV is not set.)
endif

ifeq ($(HOSTED_ZONE_ID),)
$(error HOSTED_ZONE_ID is not set.)
endif

ifeq ($(SYSADMIN_USERNAME),)
$(error SYSADMIN_USERNAME is not set.)
endif

ifeq ($(SYSADMIN_PASSWORD),)
$(error SYSADMIN_PASSWORD is not set.)
endif

ifeq ($(CORS_ALLOWED_ORIGIN),)
$(error CORS_ALLOWED_ORIGIN is not set.)
endif

DOCKER_TF := docker run --rm \
	-v $(shell pwd):$(TERRAFORM_WORKDIR) \
	-v $(shell pwd)/../api:/api \
	-v $(shell pwd)/../photo-moderator:/photo-moderator \
	-w $(TERRAFORM_WORKDIR) \
	-e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN \
	-e TF_VAR_hosted_zone_id=$(HOSTED_ZONE_ID) \
	-e TF_VAR_sysadmin_username=$(SYSADMIN_USERNAME) \
	-e TF_VAR_sysadmin_password=$(SYSADMIN_PASSWORD) \
	-e TF_VAR_cors_allowed_origin=$(CORS_ALLOWED_ORIGIN) \
	-e TF_VAR_api_bootstrap_path=/api/build/bootstrap \
	-e TF_VAR_photo_moderation_bootstrap_path=/photo-moderator/build/bootstrap \
	$(TERRAFORM_IMAGE)

deploy: init plan apply
.PHONY: deploy

init:
	$(DOCKER_TF) -chdir=environments/$(ENV) init
.PHONY: init

plan: fmt
	$(DOCKER_TF) -chdir=environments/$(ENV) plan -out=plan.tfplan
.PHONY: plan

apply:
	$(DOCKER_TF) -chdir=environments/$(ENV) apply plan.tfplan
.PHONY: apply

destroy:
	mkdir -p ../api/build
	touch ../api/build/bootstrap # need some artifact even if empty
	$(DOCKER_TF) -chdir=environments/$(ENV) destroy -auto-approve
.PHONY: destroy

validate:
	$(DOCKER_TF) -chdir=environments/$(ENV) validate
.PHONY: validate

unlock:
	$(DOCKER_TF) -chdir=environments/$(ENV) force-unlock -force $(TF_LOCK_ID)
.PHONY: unlock

fmt:
	$(DOCKER_TF) fmt -recursive
.PHONY: fmt

import:
	$(DOCKER_TF) -chdir=environments/$(ENV) import $(IMPORT_RES) $(IMPORT_ID)
.PHONY: import
