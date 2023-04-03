# TODO list

## build/

- Create CI pipeline with GitHub actions to push updated images to S3 (on schedule)
- Figure out Docker-in-Docker CI problems
- Switch reverse proxy from Caddy to Nginx: Caddy is not available in Debian repos


## deploy/

- Package into a Docker container
- Deploy fleet manager to home server
- Look into generating a pre-signed URL for VM image on fleet manager.
  That would allow to make the S3 bucket private.
  Be careful: changing URL (GET params) would trigger tf to rebuild the image,
  which in turn could(?) trigger VM rebuilds.
- `tf destroy` triggers graceful Linux shutdown. Does gitlab-runner get
  gracefully unregistered by systemd?
- Add wrapper that runs `tf apply` in a loop with a small delay (~1min)
- Package scaler app + terraform into Docker container for deployment


## scale/

- Rewrite in Go: calculate scaling actions and populate tfvars
    - Optional: use external data source in terraform to call scaler automatically
    - Pass cloud IP address as a parameter to scaler app (via `external` data source)
    - Read all configuration from JSON on stdin, add sample config for manual
      invocation during testing. Print usage message on stderr and fail after a
      timeout if no stding was received
    - Use custom datatype for unmarshalling strings that may come from env vars
- GitLab API
    - Add a delay after assigning new runners to GitLab projects - to let them
      pick up pending jobs. Or better yet, take new assignments into account when
      calculating existing jobs capacity
- Instance management
    - Calculate instance status
    - Calculate scaling actions
    - Call cleanup before destroying an instance
    - Restore instances from TF state, from GitLab API (maybe)
- Update fleet_manager Ansible role
- Update statuses asynchronously
- App ignores Ctrl+C interrupt


## Global

- Write minimal documentation for each stage
- Review todo lists in legacy branches (+MAGIC.md)
