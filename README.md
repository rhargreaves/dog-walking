# Dog Walking Service

A service for helping dog owners find dog walkers

## Infrastructure

Deployed using Terraform in a Docker container.

### Plan

`ENV=uat` is the default.

```bash
make plan
```

### Apply

```bash
make apply
```

### Destroy

```bash
make destroy
```

## Proof-of-concept tradeoffs

Key differences between this project and a real-life deployment:

* UAT and PROD should ideally be in different AWS accounts.
