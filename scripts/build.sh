#!/usr/bin/env bash
set -euo pipefail

# Where the toolchain gets built/installed.
VC_ROOT="${VC_ROOT:-$(pwd)/build}"

# Resolve to an absolute path
case "$VC_ROOT" in
    /*) ;;
    *) VC_ROOT="$(pwd)/$VC_ROOT" ;;
esac

BUILD_MODE="${1:-dev}"

echo "Installing to: $VC_ROOT"

find . -type f -name "*.sh" -exec chmod +x {} +

rm -rf "$VC_ROOT"

mkdir -p "$VC_ROOT/toolchain"
mkdir -p "$VC_ROOT/stdlib"

trap 'rm -f stage1.vcx' EXIT

# Populate the runtime used by the toolchain.
if [[ "$BUILD_MODE" == "dev" ]]; then
    cp -as "$(realpath runtime)/." "$VC_ROOT/toolchain/"
elif [[ "$BUILD_MODE" == "release" ]]; then
    cp -a runtime/. "$VC_ROOT/toolchain/"
else
    echo "usage: $0 [dev|release]" >&2
    exit 2
fi

# bootstrap stage2 compiler
nvim --headless \
  -u runtime/runtime.lua \
  -S bootstrap/stage0.lua \
  -- bootstrap/stage1.vc stage1.vcx \

nvim --headless \
  -u runtime/runtime.lua \
  -S stage1.vcx \
  -- bootstrap/stage1.vc "$VC_ROOT/toolchain/compiler.vcx" \
  2>/dev/null

# build linker
nvim --headless \
  -u runtime/runtime.lua \
  -S "$VC_ROOT/toolchain/compiler.vcx" \
  -- linker/linker.vc "$VC_ROOT/toolchain/linker.vcx" \
  2>/dev/null

# build cli
go -C cli build -o "$VC_ROOT/vc" ./cmd/VimComplete

# build stdlib
VC_ROOT="$VC_ROOT" "$VC_ROOT/vc" build stdlib -o "$VC_ROOT/stdlib" -f

# build compiler
VC_ROOT="$VC_ROOT" "$VC_ROOT/vc" build compiler -m compiler/main.vc -o "$VC_ROOT/toolchain/compiler.vcx" -f
