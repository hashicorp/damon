#!/usr/bin/env sh

set -e

if [ "$1" = 'daemon' ]; then
  shift
fi

exec damon "$@"
