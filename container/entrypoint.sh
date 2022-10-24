#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


# Environment variables
export PATH="/opt/pulumi:$PATH"
export PULUMI_BACKEND_URL="${PULUMI_BACKEND_URL:-file:///state}"
export PULUMI_SKIP_UPDATE_CHECK="${PULUMI_SKIP_UPDATE_CHECK:-true}"
: "${PULUMI_CONFIG_PASSPHRASE:?required for Pulumi state backend}"

: "${GITLAB_RUNNER_TOKEN:?required unless explicitly provided in config.toml}"
: "${GITLAB_API_TOKEN:?required unless explicitly provided in config.toml}"

: "${YC_TOKEN:?required when using default cloud provider (YandexCloud)}"
: "${YC_CLOUD_ID:?required when using default cloud provider (YandexCloud)}"
: "${YC_FOLDER_ID:?required when using default cloud provider (YandexCloud)}"



# Log dependency versions
pulumi version
pulumi plugin ls
python3 --version
python3 -m pip freeze


# Initialize state backend
pulumi login


# Invoke Pulumi program
default=(up --daemon)
exec fleet-manager ${@:-${default[*]}}
