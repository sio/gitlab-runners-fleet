# Calculate which runner instances to create/keep running

This directory contains a command-line Golang application that decides how
many GitLab runner instances should be running at any given moment.

To make that decision the application uses following data sources:

- GitLab GraphQL API: how many CI jobs are pending now?
- Runner instance HTTP API: how may runners are currently healthy and how many
  jobs can they take?
  This API is provided by a bespoke Python [script](../build/template/runner/control.py)
- JSON configuration (read from stdin).
  See [example](sample_config.json) and [config.go](app/config.go)
- Previous scaler state (read from and saved to a JSON file on local filesystem)

In addition to making scaling decisions this app does some housekeeping:

- Assign runners to all projects owned by current user. This is a workaround
  to match functionality provided by [group runners]
  (not available to individual gitlab.com users)
- Remove dead runners from GitLab UI

[group runners]: https://docs.gitlab.com/ee/ci/runners/runners_scope.html#group-runners


## Deployment status

This app is deployed together with Terraform as part of the
Docker [container](../container/README.md)
