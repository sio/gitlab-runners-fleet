#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -v

ping 8.8.8.8 -c3

apt update
apt-get install -y \
    python3 \
    python3-venv \

apt-get clean
rm -rf /var/lib/apt/lists/*
