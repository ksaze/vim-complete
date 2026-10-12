#!/usr/bin/env bash

set -euo pipefail

trap 'rm -f compiler.vcx' EXIT

echo "[1/2] Building compiler using current compiler"

build/vc build compiler -o compiler.vcx -m compiler/main.vc -f >/dev/null

echo "[2/2] Comparing compiler executables"

diff -u build/toolchain/compiler.vcx compiler.vcx

echo "PASS Bootstrap reached a fixed point."

rm -f compiler.vcx
