# Auto scaling fleet of GitLab runners

This repo contains infrastructure definition and management scripts for a
fleet of GitLab runners. The tools included enable scaling the number of
runner instances up and down as demand changes.

I've decided not to use GitLab's suggested `docker-machine` approach because:

- `docker-machine` is mostly abandoned, life support from GitLab is barely
  enough to keep it alive
- I want to be able to choose from a wider set of Cloud providers than
  `docker-machine` supports. Coding new bindings for `docker-machine` seems
  to be a pointless endeavor.
