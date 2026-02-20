#!/usr/bin/env sh
# Copyright IBM Corp. 2021, 2023
# SPDX-License-Identifier: MPL-2.0


set -e

if [ "$1" = 'daemon' ]; then
  shift
fi

exec damon "$@"
