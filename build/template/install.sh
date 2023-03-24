#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -v

# Enable cloud-init unconditionally
cp -v /etc/provision/cloud.cfg.d/* /etc/cloud/cloud.cfg.d
ln -vsf /lib/systemd/system/cloud-init.target /etc/systemd/system/multi-user.target.wants/

# Root password for interactive debugging
echo 'root:bydtyn22'|chpasswd # TODO: remove debug password

# Install required packages
apt update
apt-get install -y \
    chrony \
    coreutils \
    curl \
    docker.io \
    downtimed \
    git \
    iputils-ping \
    make \
    nftables \
    python3 \
    python3-apt \
    python3-pip \
    python3-venv \
    sudo \
    uptimed \

apt-get clean
rm -rf /var/lib/apt/lists/*

# TODO: disable services which are not universally needed
# - docker
# - gitlab-runner

# Launch Ansible
exec make common
