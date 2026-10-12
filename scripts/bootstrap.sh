#!/usr/bin/env bash

set -euo pipefail

trap 'rm -f stage1.vcx' EXIT

nvim --headless \
  -u runtime/runtime.lua \
  -S bootstrap/stage0.lua \
  -- bootstrap/stage1.vc stage1.vcx \

nvim --headless \
  -u runtime/runtime.lua \
  -S stage1.vcb \
  -- bootstrap/stage1.vc build/toolchain/compiler.vcx \
  2>/dev/null
