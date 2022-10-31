# Auto scaling fleet of GitLab CI runners

This repo contains infrastructure definitions and management scripts for a
fleet of GitLab CI runners. The tools included enable scaling the number of
runner instances up and down as demand changes.

I've decided not to use GitLab's suggested `docker-machine` approach because:

- `docker-machine` is mostly abandoned, life support from GitLab is barely
  enough to keep it alive
- I want to be able to choose from a wider set of cloud providers than
  `docker-machine` supports. Coding new bindings for `docker-machine` seems
  to be a pointless endeavor.
- I wanted a project to try Pulumi out. Plain Pulumi (see `legacy` branch) was
  not as convenient as I assumed it would be and I postponed this project
  for a while. Later I've learned about Pulumi Automation API and rewrote the
  project from scratch - the experience was totally worth it!

Because I intend to deploy the runners strictly for personal use, I've also
added a requirement which might seem arbitrary and overly restrictive in
enterprise environment: _Cloud services should incur zero costs when no CI
jobs are running_. Because of this I provision runners with Ansible from a
generic Debian image instead of baking the image with Packer or similar
tools.


## Project status

Rough around the edges but usable. Deployed and used regularly by author.


## Installation and usage

Documentation is not written yet. See Makefile for pointers, `make up` will
yell at you until you provide all required environment variables :)


## License and copyright

Copyright 2021-2022 Vitaly Potyarkin

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
