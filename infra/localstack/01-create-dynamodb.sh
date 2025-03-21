#!/bin/bash
set -euo pipefail

TABLE_NAME=local-dogs
AWS_REGION=eu-west-1

awslocal dynamodb create-table \
    --table-name ${TABLE_NAME} \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region ${AWS_REGION}
echo "DynamoDB table '${TABLE_NAME}' created."