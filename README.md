# Dog Walking Service
[![Deploy UAT and Run Tests](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml)
[![Deploy PROD & Test](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-prod.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-prod.yaml)

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

### Operations

Configure environment: `export ENV=uat`

* `make init`
* `make plan`
* `make apply`
* `make destroy`

## CI/CD

Deployed using GitHub Actions (see status badges at the top of the repo). There are also actions for burning everything to the ground to save :moneybag:. They are scheduled to run every midnight.

### Baseline Infrastructure

There's a CloudFormation template in [infra/baseline](infra/baseline) for setting up a IAM role to deploy this service. It also sets up an account-wide API Gateway CloudWatch IAM role.

## Observability

There are [alarms](infra/modules/monitoring/main.tf) for alerting on endpoint error rates but also errors from Lambdas.

There are also some basic [dashboards](infra/modules/monitoring/main.tf) defined.

## Proof-of-concept tradeoffs

I wanted to keep this project on a low-budget, and relatively low-complexity, so there are some decisions I made with regards to architecture:

Key differences between this project and a real-life deployment:

* UAT and PROD should ideally be in different AWS accounts.
* Resource IDs (such as ARNs and Hosted Zone IDs) are secrets here due to this being a public repo. These could be environment variables in a private repo.
* Ideally there would be some synthetic tests also running to check the service is ultimately available. Also, CloudWatch metrics and logs are fairly bare-bones. Using an observatibility platform such as DataDog would be ideal, with SLOs defined for each endpoint.
* A high-throughput service might be better hosted not on AWS Lambda (which can get expensive) but some always-on container orchestration platform such as Kubernetes or AWS ECS.
