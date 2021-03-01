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


## Requirements for management node

- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- GNU Make
- Python3 (with venv module)
- SSH client (for executing cleanup actions)


## Layers upon layers... How this works

1. **Systemd service** is triggered on timer to execute a single maintenance
cycle
2. **Makefile** is providing a user-friendly interface to orchestrate all the
cli tools
3. **Pulumi** takes care of executing scaling algorithms and provisioning
required number of cloud instances
4. **Cloud-init** does the most basic steps to prepare the instance to be
further provisioned
5. **Ansible** executed on the instance against itself is configuring all the
required services


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
