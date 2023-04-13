#!/bin/bash
set -euo pipefail

while true
do
    make scale apply
    sleep ${SCALE_DELAY:-1m}
done
