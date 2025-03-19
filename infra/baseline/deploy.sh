#!/bin/bash
set -euo pipefail

REGION="eu-west-1"
AWS_PROFILE="${AWS_PROFILE:-default}"
ROLE_NAME="$1"

deploy_stack() {
    local stack_name="$1"
    local template_file="$2"
    shift 2
    local additional_args=()
    if [ $# -gt 0 ]; then
        additional_args=("$@")
    fi

    echo "Deploying ${stack_name}..."
    aws cloudformation deploy \
        --template-file "${template_file}" \
        --stack-name "${stack_name}" \
        --capabilities CAPABILITY_NAMED_IAM \
        --no-fail-on-empty-changeset \
        --region "${REGION}" \
        --profile "${AWS_PROFILE}" \
        ${additional_args[@]+"${additional_args[@]}"}
}

deploy_stack "dog-walking-deploy-role" "deploy-role.yaml" --parameter-overrides "RoleName=${ROLE_NAME}"
deploy_stack "api-gateway-logging" "api-gateway-logging.yaml"

echo "Deployed successfully"