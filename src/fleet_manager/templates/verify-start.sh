#!/bin/bash
set -euo pipefail

echo Waiting for gitlab-runner to come online...
timeout 1m sh <<EOF
    while true
    do
        sleep 2s
        gitlab-runner verify && break
    done
EOF
echo Runner verification succeeded
