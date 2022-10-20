# TODO list

High level tasks go here, low level items are marked as TODO in comments next
to the relevant parts of source code. GitHub issues will be used for
collaborating with other people (if anyone gets interested) and for tracking
progress of longer tasks.


## High prioriy: blocks deployment

- Systemd in docker in docker (dind) does not work on my private runners (but
  did work on shared GitLab ones)
- Create a docker image for fleet-manager

## Medium priority: quality of life

- Remove instances that have reached maximum allowed age
- Deal with DockerHub rate limits
- Add provider: Selectel (supports nested virtualization)

## Low priority: nice to have

- S3 endpoint for shared runner cache
- Yandex: restart preempted instances automatically. Instance groups with
  auto-healing? Pulumi? Raw API calls?
- Yandex: recreate instance if cloud-init configuration has changed
- Add monitoring entrypoint
- Write better README, --help and some documentation (maybe)

## Lowest priority: maybe sometime (if ever)

- Get rid of magic values (see `legacy:MAGIC.md`)
