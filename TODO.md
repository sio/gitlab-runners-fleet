# TODO list

High level tasks go here, low level items are marked as TОDO in comments next
to the relevant parts of source code. GitHub issues will be used for
collaborating with other people (if anyone gets interested) and for tracking
progress of longer tasks.


## High priority: blocks deployment

**build/**

- Create CI pipeline with GitHub actions to push updated images to S3 (on schedule)

**deploy/**

- Add wrapper that runs `tf apply` in a loop with a small delay (~1min)
- Package into a Docker container: terraform + scaler app
- Deploy fleet manager to home server


## Medium priority: quality of life

**deploy/**

- Deal with DockerHub rate limits: set up Docker proxy on the gateway?


## Low priority: nice to have

**build/**

- Switch reverse proxy from Caddy to Nginx: Caddy is not available in Debian repos
- `tf destroy` triggers graceful Linux shutdown. Does gitlab-runner get
  gracefully unregistered by systemd? Looks like no. Do we want it to?
- Read `tf output` via subprocess and pipe instead of accessing tfstate file
  directly. Currently S3 backend for Terraform is not supported

**scale/**

- App ignores Ctrl+C interrupt

**global**

- Write minimal documentation for each stage


## Lowest priority: maybe sometime (if ever)

**deploy/**

- Look into generating a pre-signed URL for VM image on fleet manager.
  That would allow to make the S3 bucket private.
  Be careful: changing URL (GET params) would trigger tf to rebuild the image,
  which in turn could(?) trigger VM rebuilds.
- Add provider: Selectel (supports nested virtualization)
- Yandex: recreate instance if cloud-init configuration has changed
- Yandex: restart preempted instances automatically. Instance groups with
  auto-healing? Pulumi? Raw API calls?
- S3 bucket for shared runner cache. Private runners do not share cache
  between hosts by default

**global**

- Add monitoring entrypoint for fleet manager
