#!/bin/bash
set -euo pipefail

BUCKET_NAME=local-dog-images
AWS_REGION=eu-west-1

awslocal s3 mb s3://${BUCKET_NAME} --region ${AWS_REGION}
echo "S3 bucket '${BUCKET_NAME}' created."