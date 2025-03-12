TERRAFORM_IMAGE := hashicorp/terraform:latest
TF_ENV_VARS := -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN \
	-e TF_VAR_hosted_zone_id
TERRAFORM_WORKDIR := /workspace
ENV ?= uat

DOCKER_TF := docker run --rm \
	-v $(shell pwd):$(TERRAFORM_WORKDIR) \
	-w $(TERRAFORM_WORKDIR) \
	$(TF_ENV_VARS) \
	$(TERRAFORM_IMAGE)

deploy: init plan apply
.PHONY: deploy

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
	touch api.zip # dummy deploy artifact
	$(DOCKER_TF) -chdir=environments/$(ENV) destroy -auto-approve
.PHONY: destroy

validate:
	$(DOCKER_TF) -chdir=environments/$(ENV) validate
.PHONY: validate

unlock:
	$(DOCKER_TF) -chdir=environments/$(ENV) force-unlock -force $(TF_LOCK_ID)
.PHONY: unlock

fmt:
	$(DOCKER_TF) fmt
.PHONY: fmt