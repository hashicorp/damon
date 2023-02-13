#!/usr/bin/env sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

if [ "$1" = 'daemon' ]; then
  shift
fi

exec damon "$@"
