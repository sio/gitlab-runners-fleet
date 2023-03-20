#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -v

apt update
apt-get install -y \
    python3 \
    python3-venv \

apt-get clean
rm -rf /var/lib/apt/lists/*
