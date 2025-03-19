TERRAFORM_IMAGE := hashicorp/terraform:latest
TERRAFORM_WORKDIR := /workspace
ENV ?= uat

DOCKER_TF := docker run --rm \
	-v $(shell pwd):$(TERRAFORM_WORKDIR) \
	-v $(shell pwd)/../api:/api \
	-w $(TERRAFORM_WORKDIR) \
	-e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN \
	-e TF_VAR_hosted_zone_id \
	-e TF_VAR_sysadmin_password \
	-e TF_VAR_bootstrap_path=/api/build/bootstrap \
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
