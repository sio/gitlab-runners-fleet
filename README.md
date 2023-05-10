# Auto scaling fleet of GitLab CI runners

This repo contains infrastructure definitions and configuration for a fleet of
GitLab CI runners. Provided tools allow to scale the number of runner
instances up and down (to zero) as demand changes.

I decided not to use GitLab's suggested `docker-machine` approach because:

- `docker-machine` is mostly abandoned, life support from GitLab is barely
  enough to keep it alive. GitLab themselves
  [are looking](https://docs.gitlab.com/ee/architecture/blueprints/runner_scaling/)
  for an alternative solution.
- I want to be able to choose from a wider set of cloud providers than
  `docker-machine` supports. Coding new bindings for `docker-machine` seems
  to be a pointless endeavor.
- I wanted a project to learn Terraform/Pulumi. Previous two iterations of this
  project were created with Pulumi and Pulumi Automation API. Pulumi stopped
  developing bindings for the cloud I use (Yandex Cloud) and my Python code was
  not as clean as I would like, hence this (third) rewrite to Terraform & Go.

I intend to deploy the runners only for personal use and I aim to
architect my infra to incur (next to) zero costs when no CI jobs are running.

## Project status

Ready for deployment.
Deployed and used regularly by author.

Rewritten to Terraform after
[Pulumi Automation API](https://github.com/sio/gitlab-runners-fleet/tree/legacy/02-pulumi-automation-api)
and [plain Pulumi](https://github.com/sio/gitlab-runners-fleet/tree/legacy/01-pulumi-plain).


## Infrastructure overview

- Persistent resources
    - S3 bucket that holds base image for cloud VMs
- Scaling down to zero on demand
    - 0 to N runner hosts: Debian hosts with GitLab runner daemon + Docker
      executor
    - 0 to 1 gateway: a simple cloud VM which has a public IPv4 and acts as a
      router, firewall, reverse proxy for HTTP API and as a caching proxy for
      Docker Hub.
    - Other required resources (networks, IP addresses, VM images)


## Usage

- Prepare S3 bucket with a prebaked VM image: [build/README.md](build/README.md)
- Launch fleet manager container: [container/README.md](container/README.md)


## Underlying technology

It's incredible how much power and how easily can a single individual wield
thanks to modern tech! This project is made possible by standing on shoulders
of giants:

<table><tr><td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</td><td align="center">

[Ansible](https://docs.ansible.com/) &nbsp;&nbsp;·&nbsp;&nbsp;
[cloud-init](https://cloudinit.readthedocs.io/) &nbsp;&nbsp;·&nbsp;&nbsp;
[Debian](https://debian.org) &nbsp;&nbsp;·&nbsp;&nbsp;
[Docker](https://docs.docker.com/) &nbsp;&nbsp;·&nbsp;&nbsp;
[GitHub Actions](https://docs.github.com/actions) &nbsp;&nbsp;·&nbsp;&nbsp;
[GitLab CI](https://docs.gitlab.com/ee/ci/) &nbsp;&nbsp;·&nbsp;&nbsp;
[GNU Make](https://www.gnu.org/software/make/) &nbsp;&nbsp;·&nbsp;&nbsp;
[Golang](https://go.dev) &nbsp;&nbsp;·&nbsp;&nbsp;
[GraphQL](https://graphql.org/) &nbsp;&nbsp;·&nbsp;&nbsp;
[nftables](https://netfilter.org/projects/nftables/) &nbsp;&nbsp;·&nbsp;&nbsp;
[Python](https://python.org) &nbsp;&nbsp;·&nbsp;&nbsp;
[Qemu (qemu-utils)](https://www.qemu.org/) &nbsp;&nbsp;·&nbsp;&nbsp;
[systemd](https://systemd.io) &nbsp;&nbsp;·&nbsp;&nbsp;
[Terraform](https://www.terraform.io/) &nbsp;&nbsp;·&nbsp;&nbsp;
[Yandex Cloud](https://cloud.yandex.com) &nbsp;&nbsp;·&nbsp;&nbsp;

</td><td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</td></tr></table>


## License and copyright

Copyright 2021-2023 Vitaly Potyarkin

```
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
