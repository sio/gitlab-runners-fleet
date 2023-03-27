# Assorted development notes

It is often difficult to remember down the road why this or that decision was
taken. These notes are mostly intended to be consumed by my future self.


## Pulumi vs Terraform

- Plain Pulumi was not very interesting, running it in a loop involved a lot
  of shell/Makefile glue code and was not elegant
- Pulumi Automation API was very cool! Unfortunately, Pulumi had stopped
  developing the plugin for Yandex Cloud
- Terraform-CDK is as inelegant as raw Pulumi, there is nothing like
  Automation API (yet)
- Plain Terraform requires the same glue code to run in a loop as plain
  Pulumi, but at this point it appears to be the least worst option.
  Declarative code is nice enough for my simple infra though.


## Yandex Cloud

- Bucket creation takes a long time for both `terraform plan` and `terraform
  apply` ("Refreshing state..."). This should not be an ephemeral resource
- S3 object does not get updated by `yandex_storage_object` when its
  `source` file is modified. There are no etag/checksum parameters to trigger
  an update.
- `yandex_compute_image` `source_url` MUST point to Yandex Cloud object
  storage
- Application load balancer is too complex to configure and much more
  expensive than a single VM instance for a gateway.


## Building VM image

- Packer does not appear to provide an easy way to modify qcow2 image on a
  host without virtualization support (qemu without kvm is painfully slow),
  hence we use a bespoke script which relies on qemu-nbd and chroot.
  This still requires root access to the build host.
- mkosi seems nice, but it can only build from scratch via debootstrap.
  Upstream Debian images are rather good, there is no need to redo the work
  of Debian Cloud Team

## Bringup sequence

- Create S3 bucket: `make -C build bucket` (once)
- Build base VM image and upload to S3: `make -C build image compact upload`
  (regularly in CI)
- Create/update the rest of the infra: `make -C deploy`
  (regularly on fleet manager)
