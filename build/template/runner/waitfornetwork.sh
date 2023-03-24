#!/bin/sh
echo 'Waiting for Internet connection...'
REMOTE_IP=8.8.8.8
timeout 10m sh <<EOF
    while true
    do
        ping -c1 -w3 "$REMOTE_IP" && break
    done
EOF
