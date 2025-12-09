# Copyright IBM Corp. 2021, 2023
# SPDX-License-Identifier: MPL-2.0

container {
  dependencies = true
  alpine_secdb = true
  secrets      = true
}

binary {
  secrets    = true
  go_modules = true
  osv        = true
  oss_index  = false
  nvd        = false
}
