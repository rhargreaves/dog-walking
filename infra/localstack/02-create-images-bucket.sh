#!/bin/bash
set -euo pipefail

BUCKET_NAME=local-dog-images
PENDING_BUCKET_NAME=local-pending-dog-images
AWS_REGION=eu-west-1

awslocal s3 mb s3://${BUCKET_NAME} --region ${AWS_REGION}
echo "S3 bucket '${BUCKET_NAME}' created."

awslocal s3 mb s3://${PENDING_BUCKET_NAME} --region ${AWS_REGION}
echo "S3 bucket '${PENDING_BUCKET_NAME}' created."