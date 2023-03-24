#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

#
# Wait for a specific command on stdin
#

while read LINE
do
    if [[ "$LINE" == "UNREGISTER" ]]
    then
        echo "INFO: received $LINE command" >&2
        exit 0
    else
        echo "WARNING: invalid command: $LINE" >&2
    fi
done
echo "ERROR: infinite loop broken" >&2
exit 1
