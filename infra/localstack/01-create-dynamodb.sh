#!/bin/bash
TABLE_NAME=local-dogs
awslocal dynamodb create-table \
    --table-name ${TABLE_NAME} \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST

echo "DynamoDB table '${TABLE_NAME}' created."