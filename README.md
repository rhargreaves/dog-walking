# Dog Walking Service
[![Deploy UAT and Run Tests](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml)

A service for helping dog owners find dog walkers

## API

### Build

```sh
make build
```

The API is compiled to `out/bootstrap` ready to be deployed as a Lambda function.

### Local Testing

```sh
make test-local
```

Dependencies are swapped out as follows:

| Production/UAT | Local |
|----------------|-------|
| AWS Lambda/API Gatway | AWS SAM Local |
| AWS DynamoDB   | Localstack DynamoDB |
| AWS S3   | Localstack S3 |
| AWS Rekognition   | Lookup based on [image MD5](api/internal/rekognition_stub/hashes.go) |

## Infrastructure

Deployed using Terraform.


```sh
cd infra
```

### Plan

`ENV=uat` is the default.

```sh
make plan
```

### Apply

```sh
make apply
```

### Destroy

```sh
make destroy
```

## CI/CD

Deployed using GitHub Actions. There's a CloudFormation template in [deploy-infra](deploy-infra) for setting up a IAM role to deploy this service.

## Proof-of-concept tradeoffs

Key differences between this project and a real-life deployment:

* UAT and PROD should ideally be in different AWS accounts.
* Resource IDs (such as ARNs and Hosted Zone IDs) are secrets here due to this being a public repo. These could be environment variables in a private repo.
