#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -v

# Root password for interactive debugging
echo 'root:bydtyn22'|chpasswd # TODO: remove debug password

# Enable cloud-init unconditionally
cp -v /etc/provision/cloud.cfg.d/* /etc/cloud/cloud.cfg.d
ln -vsf /lib/systemd/system/cloud-init.target /etc/systemd/system/multi-user.target.wants/

# Generate temporary ssh host keys (for Ansible)
ssh-keygen -A

# Install Ansible prerequisites
apt update
apt-get install -y --no-install-recommends \
    coreutils \
    curl \
    git \
    make \
    python3 \
    python3-apt \
    python3-pip \
    python3-venv \
    sudo \

# Launch Ansible. Even though common.yml is a very simple playbook we need to
# make sure we have a working Ansible install in the resulting image
make -C /etc/provision common

# Remove temporary host keys
rm -vf /etc/ssh/*host*key*

# Reduce image size
apt-get clean
rm -rf /var/lib/apt/lists/*
