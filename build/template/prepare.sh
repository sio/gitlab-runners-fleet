#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -v

# Debug information
tail -v -n+0 /etc/*release*
networkctl
ip addr
ip route

# Copy provisioning config from ephemeral bind mount
mkdir -p /etc/provision
cp -avr /tmp/prepare/template/* /etc/provision
chown -R root:root /etc/provision

# Configure DNS name resolution
mv -v /etc/resolv.conf{,.orig}
echo 'nameserver 8.8.8.8' > /etc/resolv.conf

# Create temporary directory for sshd config validation (cleaned up automatically - tmpfs)
[[ $(df -Th /run | tail -n1 | cut -d\   -f1) == "tmpfs" ]]
mkdir -p /run/sshd

# Enable cloud-init unconditionally
cp -v /etc/provision/cloud.cfg.d/* /etc/cloud/cloud.cfg.d
ln -vsf /lib/systemd/system/cloud-init.target /etc/systemd/system/multi-user.target.wants/

# Generate temporary ssh host keys (for Ansible)
ssh-keygen -A

# Wait for container network to come online
/usr/lib/systemd/systemd-networkd-wait-online --interface=host0 --timeout=120 # seconds

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

# Restore /etc/resolv.conf
mv -v /etc/resolv.conf{.orig,}
