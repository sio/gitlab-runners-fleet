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


# nss_wrapper <https://cwrap.org/nss_wrapper.html>
fakeid() {
    if [[ "$UID" == 0 || "$EUID" == 0 ]]; then return; fi
    if getent passwd "$UID" &>/dev/null; then return; fi

    export NSS_WRAPPER_PASSWD=$(mktemp)
    export NSS_WRAPPER_GROUP=$(mktemp)

    getent passwd > "$NSS_WRAPPER_PASSWD"
    getent group > "$NSS_WRAPPER_GROUP"

    sed -i '/^fleetmanager:/d' "$NSS_WRAPPER_PASSWD" "$NSS_WRAPPER_GROUP"
    echo "fleetmanager:!:$UID:$(id -g)::/home/fleetmanager:/usr/sbin/nologin" >> "$NSS_WRAPPER_PASSWD"
    echo "fleetmanager:!:$(id -g):fleetmanager" >> "$NSS_WRAPPER_GROUP"

    export LD_PRELOAD=libnss_wrapper.so
    export USER=fleetmanager
}; fakeid


# Log some troubleshooting information
id
pulumi version
pulumi plugin ls
python3 --version
python3 -m pip freeze


# Initialize state backend
pulumi login


# Invoke Pulumi program
default=(up --daemon)
exec fleet-manager ${@:-${default[*]}}
