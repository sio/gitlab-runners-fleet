# Fleet manager Docker container

This container combines all components of this project into one deployable
artifact.

Resulting image is published at `ghcr.io/sio/gitlab-runners-fleet:v3`

See [compose.yml](compose.yml) for a deployment example.

## Deployment status

All changes to this repo are continuously integrated via GitHub Actions
[workflow](../.github/workflows/container.yml).
Container image is also automatically rebuilt on monthly schedule to reflect
changes in underlying images (including security updates).

The image itself has not been deployed to my homeprod (yet). `#TODO`
