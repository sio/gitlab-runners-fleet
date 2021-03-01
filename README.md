# Auto scaling fleet of GitLab CI runners

This repo contains infrastructure definition and management scripts for a
fleet of GitLab CI runners. The tools included enable scaling the number of
runner instances up and down as demand changes.

I've decided not to use GitLab's suggested `docker-machine` approach because:

- `docker-machine` is mostly abandoned, life support from GitLab is barely
  enough to keep it alive
- I want to be able to choose from a wider set of Cloud providers than
  `docker-machine` supports. Coding new bindings for `docker-machine` seems
  to be a pointless endeavor.
- I wanted a project to try Pulumi out


## Requirements for management node

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- GNU Make
- Python3 (with venv module)
- SSH client (for executing cleanup actions)

## Some useful links

- [Data science on demand: spinning up a Wallaroo cluster is easy
  with
  Pulumi](https://www.pulumi.com/blog/data-science-on-demand-spinning-up-a-wallaroo-cluster-is-easy-with-pulumi/) -
  ([configuration](https://github.com/WallarooLabs/wallaroo_blog_examples/tree/master/provisioned-classifier/pulumi),
  [Makefile](https://github.com/WallarooLabs/wallaroo_blog_examples/blob/master/provisioned-classifier/Makefile))
- [Gitlab Runner autoscaling infrastructure on Hetzner Cloud with Terraform](https://www.stefanwienert.de/blog/2019/04/06/gitlab-runner-autoscaling-infrastructure-on-hetzner-cloud-with-terraform/)
