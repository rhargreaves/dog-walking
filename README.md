# Dog Walking Service
[![Deploy UAT and Run Tests](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml)
[![Deploy PROD & Test](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-prod.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-prod.yaml)

A service for helping dog owners find dog walkers.

## Architecture

<img src="docs/arch.png" alt="architecture diagram" />

## API & Photo Moderator S3 Trigger

### Build

```sh
make build
```

The API & Photo Moderator S3 trigger are compiled to `api/build/bootstrap` and `photo-moderator/build/bootstrap` respectively.

### Local Testing

```sh
make test-local
```

### E2E Testing

Ensure all environment variables in [.example.env](.example.env) have been set - either in a `.env` file or otherwise.

```sh
make test
```

Dependencies are swapped out as follows:

| Production/UAT | Local |
|----------------|-------|
| AWS Lambda/API Gatway | AWS SAM Local |
| AWS DynamoDB   | Localstack DynamoDB |
| AWS S3   | Localstack S3 |
| AWS Rekognition   | Lookup based on [image MD5](api/internal/rekognition_stub/hashes.go) |

In addition, locally, the Photo Moderator S3 trigger is manually ran as part of the E2E tests when a photo is uploaded. This is because neither Localstack nor AWS SAM can effectively emulate the S3 trigger behaviour themselves.

## Infrastructure

Deployed using Terraform (change to the `infra` directory).

### Operations

Ensure you've set up an `.env` file! (see [.example.env](.example.env))

* `make init`
* `make plan`
* `make apply`
* `make destroy`

## API Documentation

Docs and OpenAPI specs are provided online via Swagger (`/api-docs`)

### Photo Moderation

Uploaded photos are verified using *AWS Rekognition* for any inappropriate/unwanted content (for example nudity, guns/violence). If the image is deemed safe, then the image is checked for a dog. If a dog is found, the image is approved and made available via the CDN. The dog breed (if also detected) is updated on the dog's profile.

## CI/CD

Deployed using GitHub Actions (see status badges at the top of the repo). There are also actions for burning everything to the ground to save :moneybag:.

### Baseline Infrastructure

There's a CloudFormation template in [infra/baseline](infra/baseline) for setting up a IAM role to deploy this service. It also sets up an account-wide API Gateway CloudWatch IAM role.

## Observability

There are [alarms](infra/modules/monitoring/main.tf) for alerting on endpoint error rates but also errors from Lambdas.

There are also some basic [dashboards](infra/modules/monitoring/main.tf) defined.

## Load Testing

See [load-test/README.md](load-test/README.md)

## Proof-of-concept tradeoffs

I wanted to keep this project on a low-budget, and relatively low-complexity, so there are some decisions I made with regards to architecture:

Key differences between this project and a real-life deployment:

* UAT and PROD should ideally be in different AWS accounts.
* Resource IDs (such as ARNs and Hosted Zone IDs) are secrets here due to this being a public repo. These could be environment variables in a private repo.
* Ideally there would be some synthetic tests also running to check the service is ultimately available. Also, CloudWatch metrics and logs are fairly bare-bones. Using an observatibility platform such as DataDog would be ideal, with SLOs defined for each endpoint.
* A high-throughput service might be better hosted not on AWS Lambda (which can get expensive) but some always-on container orchestration platform such as Kubernetes or AWS ECS.
