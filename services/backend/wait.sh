#!/bin/sh

set -e

cmd="$@"

sleep 5


>&2 echo "Backernd is run"
exec $cmd