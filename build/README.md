# Build VM image for Yandex Compute Cloud

This directory contains scripts and configuration for preparing a VM image
that will be deployed as Yandex Compute instance later.

## Overview

- `make bucket` will create an S3 bucket in Yandex Object Storage (uses
  Terraform under the hood; `make destroy` will destroy the bucket if it's
  empty)
- `make image` will build the VM image using following steps (*requires root*):
    - Download latest [Debian Official Cloud image](https://cloud.debian.org/images/cloud/)
    - Use qemu-utils to convert image to raw format and to grow root partition
    - Execute `template/prepare.sh` in systemd-nspawn container. Heavy lifting is
      delegated then to Ansible
    - Use qemu-utils again to compress resulting image back to qcow2
- `make upload` will upload the image to S3 bucket using awscli

## Deployment status

This VM image is continuously delivered to pre-existing S3 bucket
using GitHub Actions (see [.github/workflows/server.yml](../.github/workflows/server.yml))
