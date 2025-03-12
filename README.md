# Dog Walking Service
[![Deploy UAT and Run Tests](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml/badge.svg)](https://github.com/rhargreaves/dog-walking/actions/workflows/deploy-uat.yaml)

A service for helping dog owners find dog walkers

## API

### Build

```sh
make build
```

The API is packaged into `api.zip` ready to be uploaded to a lambda function.

### Test

```sh
make test
```

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

## Proof-of-concept tradeoffs

Key differences between this project and a real-life deployment:

* UAT and PROD should ideally be in different AWS accounts.
* Resource IDs (such as ARNs and Hosted Zone IDs) are secrets here due to this being a public repo. These could be environment variables in a private repo.
