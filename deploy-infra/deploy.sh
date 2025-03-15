#!/bin/bash
set -euo pipefail

STACK_NAME="dog-walking-deploy-role"
TEMPLATE_FILE="deploy-role.yaml"
REGION="eu-west-1"
AWS_PROFILE="${AWS_PROFILE:-default}"
ROLE_NAME="$1"

echo "Deploying deploy role (${STACK_NAME})..."
aws cloudformation deploy \
    --template-file "${TEMPLATE_FILE}" \
    --stack-name "${STACK_NAME}" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-fail-on-empty-changeset \
    --region "${REGION}" \
    --profile "${AWS_PROFILE}" \
    --parameter-overrides \
        RoleName="${ROLE_NAME}"

echo "Deployed"