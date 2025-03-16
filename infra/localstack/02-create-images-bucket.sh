#!/bin/bash
BUCKET_NAME=local-dog-images
awslocal s3 mb s3://${BUCKET_NAME} --region eu-west-1

echo "S3 bucket '${BUCKET_NAME}' created."