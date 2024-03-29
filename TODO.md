# TODO list

High level tasks go here, low level items are marked as TОDO in comments next
to the relevant parts of source code. GitHub issues will be used for
collaborating with other people (if anyone gets interested) and for tracking
progress of longer tasks.


<!--
## High priority: blocks deployment
-->


## Medium priority: quality of life

**build/**

- Public SSH keys are hardcoded into VM image: `build/template/common.yml`.
  That's less than ideal, and should be made configurable


## Low priority: nice to have

**build/**

- Switch reverse proxy from Caddy to Nginx: Caddy is not available in Debian repos
- `terraform destroy` triggers graceful Linux shutdown. Does gitlab-runner get
  gracefully unregistered by systemd? Looks like no. Do we want it to?

**scale/**

- App ignores Ctrl+C interrupt


## Lowest priority: maybe sometime (if ever)

**deploy/**

- Look into generating a pre-signed URL for VM image on fleet manager.
  That would allow to make the S3 bucket private.
  Be careful: changing URL (GET params) would trigger terraform to rebuild the image,
  which in turn could(?) trigger VM rebuilds, use `lifecycle_rule.ignore_changes`
- Add provider: Selectel (supports nested virtualization)
- Yandex: recreate instance if cloud-init configuration has changed
- Yandex: restart preempted instances automatically. Instance groups with
  auto-healing? Pulumi? Raw API calls?
- S3 bucket for shared runner cache. Private runners do not share cache
  between hosts by default

**global**

- Add monitoring entrypoint for fleet manager. Simple implementation: write
  metrics to file after each `terraform apply` to be consumed by Prometheus
  (node_exporter / textfile collector)
