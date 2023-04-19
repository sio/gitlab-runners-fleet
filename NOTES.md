# Assorted development notes

It is often difficult to remember down the road why this or that decision was
taken. These notes are mostly intended to be consumed by my future self.


## Pulumi vs Terraform

- [Plain Pulumi][v01] was not very interesting, running it in a loop involved a lot
  of shell/Makefile glue code and was not elegant
- [Pulumi Automation API][v02] was very cool! Unfortunately, Pulumi had stopped
  developing the plugin for Yandex Cloud
- Terraform-CDK is as inelegant as raw Pulumi, there is nothing like
  Automation API (yet)
- Plain Terraform requires the same glue code to run in a loop as plain
  Pulumi, but at this point it appears to be the least worst option.
  Declarative code is nice enough for my simple infra though.

[v01]: https://github.com/sio/gitlab-runners-fleet/tree/legacy/01-pulumi-plain
[v02]: https://github.com/sio/gitlab-runners-fleet/tree/legacy/02-pulumi-automation-api


## Yandex Cloud

- Bucket creation takes a long time for both `terraform plan` and `terraform
  apply` ("Refreshing state..."). This should not be an ephemeral resource
- S3 object does not get updated by `yandex_storage_object` when its
  `source` file is modified. There are no etag/checksum parameters to trigger
  an update.
- `yandex_compute_image` `source_url` MUST point to Yandex Cloud object
  storage.
  As far as I know there is no way to provide just `(bucket_name, object_key)`
  tuple like in `aws_ebs_snapshot_import`, only full URL. This is
  inconvenitent because presigned URLs are by their nature unstable and would
  trigger image reupload after each expiration. I just make the bucket public
  for now, there is no confidential information in the image. The only harm an
  attacker can inflict is to grow my egress bill by downloading the image in a
  loop. Budgets take care of that for now, but
  I should look into limiting bucket access based on IP address also (TODO).
- Application load balancer is too complex to configure and much more
  expensive than a single VM instance for a gateway.


## GitLab

- GitLab API has matured significantly since [previous iteration][v02] of this
  project. Now all runner operations are exposed via GraphQL API
  and there is no more need for the REST API. I'm glad these changes happened
  because I find GraphQL API to be a lot more convenient.


## Building VM image

- Packer does not appear to provide an easy way to modify qcow2 image on a
  host without virtualization support (qemu without kvm is painfully slow),
  hence we use a [bespoke script] which relies on systemd-nspawn
  This still requires root access to the build host - but that's not a problem
  in GitHub Actions environment.
- mkosi seems nice, but it can only build from scratch via debootstrap.
  Upstream [Debian images] are rather good, there is no need to redo the work
  of Debian Cloud Team

[bespoke script]: build/makefiles/image.mk
[Debian images]: https://cloud.debian.org


## Bringup sequence

- Create S3 bucket: `make -C build bucket` (once)
- Build base VM image and upload to S3: `make -C build image compact upload`
  (regularly in CI)
- Create/update the rest of the infra: `make -C deploy loop`
  (continuously on fleet manager)
