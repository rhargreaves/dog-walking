TERRAFORM_IMAGE := hashicorp/terraform:latest
TF_ENV_VARS := -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN -e TF_VAR_*
TERRAFORM_WORKDIR := /workspace
ENV ?= uat

DOCKER_TF := docker run --rm \
	-v $(shell pwd):$(TERRAFORM_WORKDIR) \
	-w $(TERRAFORM_WORKDIR) \
	$(TF_ENV_VARS) \
	$(TERRAFORM_IMAGE)

init:
	$(DOCKER_TF) -chdir=environments/$(ENV) init
.PHONY: init

plan:
	$(DOCKER_TF) -chdir=environments/$(ENV) plan -out=plan.tfplan
.PHONY: plan

apply:
	$(DOCKER_TF) -chdir=environments/$(ENV) apply plan.tfplan
.PHONY: apply

destroy:
	$(DOCKER_TF) -chdir=environments/$(ENV) destroy -auto-approve
.PHONY: destroy

validate:
	$(DOCKER_TF) -chdir=environments/$(ENV) validate
.PHONY: validate

fmt:
	$(DOCKER_TF) fmt
.PHONY: fmt