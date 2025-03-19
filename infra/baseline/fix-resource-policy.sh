#!/bin/bash

# from https://stackoverflow.com/a/65239057/2323497
# fixes: Error: creating API Gateway v2 Stage ($default):
# operation error ApiGatewayV2: CreateStage, https response error StatusCode: 400,
# RequestID: x, BadRequestException: Insufficient permissions to enable logging
aws logs put-resource-policy \
	--policy-name AWSLogDeliveryWrite20150319 \
	--policy-document "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"AWSLogDeliveryWrite\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"delivery.logs.amazonaws.com\"},\"Action\":[\"logs:CreateLogStream\",\"logs:PutLogEvents\"],\"Resource\":[\"*\"]}]}"