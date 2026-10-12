REPO        := $(CURDIR)
VC_ROOT     ?= $(REPO)/build
NVIM        := nvim
RUNTIME_LUA := runtime/runtime.lua

# Resolve VC_ROOT to an absolute path no matter how it was supplied —
# unset (defaults above), an environment variable (VC_ROOT=foo make cli),
# or a command-line assignment (make cli VC_ROOT=foo). This matters
# because `go -C cli build -o $(VC_ROOT)/vc` resolves a relative path
# against cli/, not the repo root. `override` is needed because a
# command-line-assigned VC_ROOT otherwise has higher precedence than any
# assignment in this file, including this one. realpath -m (not plain
# realpath) is required since the path may not exist yet.
override VC_ROOT := $(shell realpath -m "$(VC_ROOT)")

.PHONY: help dirs runtime bootstrap linker cli compiler stdlib \
        test build release clean

help:
	@echo "vim-complete dev targets"
	@echo ""
	@echo "  make cli            build the vc CLI only"
	@echo "  make bootstrap      rebuild the compiler via the bootstrap chain"
	@echo "  make linker         rebuild the linker              (needs: bootstrap)"
	@echo "  make compiler       rebuild the full compiler via vc (needs: cli, bootstrap, linker)"
	@echo "  make stdlib         rebuild the stdlib via vc        (needs: cli, compiler)"
	@echo ""
	@echo "  make test           run the full suite (scripts/run_test_suite.sh)"
	@echo ""
	@echo "  make build          full dev build   (scripts/build.sh dev)"
	@echo "  make release        full release build (scripts/build.sh release)"
	@echo "  make clean          remove build/ and transient artifacts"
	@echo ""
	@echo "Set VC_ROOT to build/install somewhere other than ./build, e.g.:"
	@echo "  VC_ROOT=/opt/vc make release"

dirs:
	@echo "VC_ROOT: $(VC_ROOT)"
	mkdir -p $(VC_ROOT)/toolchain $(VC_ROOT)/stdlib

# Symlinks runtime/*.lua into $(VC_ROOT)/toolchain so edits to runtime.lua /
# debugger.lua / common.lua are picked up immediately, without a rebuild.
# Mirrors scripts/build.sh's dev mode (cp -as).
runtime:
	cp -as "$(REPO)/runtime/." $(VC_ROOT)/toolchain/

# Bootstraps the compiler from bootstrap/stage1.vc using the minimal Lua
# stage0 compiler, then re-compiles stage1.vc with itself
bootstrap:
	@set -e; \
	trap 'rm -f stage1.vcx' EXIT; \
	$(NVIM) --headless -u $(RUNTIME_LUA) -S bootstrap/stage0.lua -- bootstrap/stage1.vc stage1.vcx; \
	$(NVIM) --headless -u $(RUNTIME_LUA) -S stage1.vcx -- bootstrap/stage1.vc $(VC_ROOT)/toolchain/compiler.vcx 2>/dev/null

# Compiles the linker using the bootstrap compiler.
linker:
	$(NVIM) --headless -u $(RUNTIME_LUA) -S $(VC_ROOT)/toolchain/compiler.vcx -- linker/linker.vc $(VC_ROOT)/toolchain/linker.vcx 2>/dev/null

cli:
	go -C cli build -o $(VC_ROOT)/vc ./cmd/VimComplete

# Rebuilds the *full* compiler (compiler/main.vc and friends) using the
# bootstrap compiler + linker, via the vc CLI itself. This overwrites the
# minimal bootstrap compiler.vcx with the real one.
compiler:
	VC_ROOT=$(VC_ROOT) $(VC_ROOT)/vc build compiler -m compiler/main.vc -o $(VC_ROOT)/toolchain/compiler.vcx -f

stdlib:
	VC_ROOT=$(VC_ROOT) $(VC_ROOT)/vc build stdlib -o $(VC_ROOT)/stdlib -f

test:
	VC_ROOT=$(VC_ROOT) ./scripts/run_test_suite.sh

build:
	VC_ROOT=$(VC_ROOT) ./scripts/build.sh dev

release:
	VC_ROOT=$(VC_ROOT) ./scripts/build.sh release

clean:
	rm -rf $(VC_ROOT) test-scripts vc-build
