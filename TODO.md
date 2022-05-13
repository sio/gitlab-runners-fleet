# TODO list

High level tasks go here, low level items are marked as TODO in comments next
to the relevant parts of source code. GitHub issues will be used for
collaborating with other people (if anyone gets interested) and for tracking
progress of longer tasks.


## High level tasks

- Daemon mode: execute Pulumi in a loop with configurable delay
- Yandex: restart preempted instances automatically. Instance groups with
  auto-healing? Pulumi? Raw API calls?
- Yandex: recreate instance if cloud-init configuration has changed
- Figure out dind setup: https://gitlab.com/sio/ci-with-molecule-git/-/jobs/1105519117
- Investigate: new runners do not pick up old pending jobs unless kicked via webui
- Add provider: Selectel (supports nested virtualization)
- Remove instances that have reached maximum allowed age
- Add monitoring entrypoint
