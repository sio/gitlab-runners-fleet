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


## Project status

Work in progress. Not working yet and may break things.

```
# TODO: dind connection refused: https://gitlab.com/sio/ci-with-molecule-git/-/jobs/1105519117
# TODO: new runners do not pick up old pending jobs unless kicked via webui
```


## Requirements for management node

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- GNU Make
- Python3 (with venv module)
- SSH client (for executing cleanup actions)


## Layers upon layers... How this works

1. [Systemd service] is triggered on timer to execute a single maintenance
cycle
2. [Makefile] is providing a user-friendly interface to orchestrate all the
cli tools
3. [Pulumi] takes care of executing scaling algorithms and provisioning
required number of cloud instances
4. [Cloud-init] does the most basic steps to prepare the instance to be
further provisioned
5. [Ansible] executed on the instance against itself is configuring all the
required services
6. [Python script] invokes GitLab API to assign created runners to all
   projects owned by current user. Similar functionality is provided out of
   the box for [group runners](https://docs.gitlab.com/ee/ci/runners/#types-of-runners)
   by GitLab itself.

[Systemd service]: systemd/
[Makefile]: Makefile
[Pulumi]: pulumi/
[Cloud-init]: instance/cloudinit.yml.j2
[Ansible]: instance/playbook.yml
[Python script]: pulumi/assign_runners.py


## Installation and usage

#### Simple interactive invocation

- Clone this repo
- Execute `make` and follow error messages. Makefile will yell at you if you
  don't provide the required environment variables or if some dependency
  application is not installed.
- You can use `make PULUMI_AUTO_INSTALL=yes` to install all the dependencies
  automatically on a throwaway Debian/Ubuntu machine: some system-wide
  packages will be installed, root privileges required.


#### Longterm management node

- Clone this repo / Unpack the tarball to `/etc/gitlab-runners-fleet/`
- Install `systemd/gitlab-runners-fleet.{service,timer}` into proper systemd
  locations
- Create `/etc/gitlab-runners-fleet.env` with all the required secrets. See
  `systemd/environment.sample` for inspiration.
- Reload systemd daemon, enable and start `gitlab-runners-fleet.timer`
- Check journalctl for errors in a few minutes

I wrote an [ansible role] to automate installation on my management node.
Reuse this role or use it as a template for your own role if you wish.

[ansible role]: https://gitlab.com/sio/server_common/-/tree/master/ansible/roles/ci_runners_manager


## Some useful links

- [Data science on demand: spinning up a Wallaroo cluster is easy
  with
  Pulumi](https://www.pulumi.com/blog/data-science-on-demand-spinning-up-a-wallaroo-cluster-is-easy-with-pulumi/) -
  ([configuration](https://github.com/WallarooLabs/wallaroo_blog_examples/tree/master/provisioned-classifier/pulumi),
  [Makefile](https://github.com/WallarooLabs/wallaroo_blog_examples/blob/master/provisioned-classifier/Makefile))
- [Gitlab Runner autoscaling infrastructure on Hetzner Cloud with Terraform](https://www.stefanwienert.de/blog/2019/04/06/gitlab-runner-autoscaling-infrastructure-on-hetzner-cloud-with-terraform/)


## License and copyright

Copyright 2021 Vitaly Potyarkin

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
