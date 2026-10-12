#!/usr/bin/env bash

set -euo pipefail

trap 'rm -f stage1.vcx stage2.vcx' EXIT

echo "[1/3] Building stage1 compiler"

nvim --headless \
  -u runtime/runtime.lua \
  -S bootstrap/stage0.lua \
  -- bootstrap/stage1.vc stage1.vcx \
  >/dev/null

echo "[2/3] Building stage2 compiler"

nvim --headless \
  -u runtime/runtime.lua \
  -S stage1.vcx \
  -- bootstrap/stage1.vc stage2.vcx \
  2>/dev/null

echo "[3/3] Comparing stage1 and stage2 executables"

diff -u stage1.vcx stage2.vcx

rm -f stage1.vcx stage2.vcx

echo "PASS stage0 and stage1 are equivalent."
