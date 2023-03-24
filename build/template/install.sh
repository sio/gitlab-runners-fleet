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
    python3 \
    python3-venv \

apt-get clean
rm -rf /var/lib/apt/lists/*
