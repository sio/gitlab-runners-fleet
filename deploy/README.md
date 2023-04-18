# Deploy GitLab runners to Yandex Cloud

This directory contains Terraform configuration for deploying multiple Yandex
Compute instances as GitLab runners.

This deployment depends on:

- [VM image] must be prepared and uploaded to S3 bucket
- [scale/] app must be compiled and available in $PATH or pointed to by an
  environment variable

## Usage

- `make scale apply` - a single scaling iteration and a single (static) infra
  deployment
- `make loop` - continuously reevaluate scaling decision and resize infra if
  needed

## Deployment status

This configuration is deployed together with [scale/] app
as part of the Docker [container](../container/README.md)

[scale/]: ../scale/README.md
[VM image]: ../build/README.md
