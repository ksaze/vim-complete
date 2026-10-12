# Architecture

A compiled program is a small Lua script, and running it means sourcing that script
into a headless Neovim instance, which replays a validated sequence of Normal-mode commands
against real buffers via `vim.cmd("normal! " .. cmd)`. Neovim *is* the
runtime — `vc`, the compiler, and the linker all exist to produce and
orchestrate input for it.

## Compiled object format

A `.vco`/`.vcx` file is a Lua script, structured as:

```lua
lua << EOF
VC.run({
  "<line 1's raw keystrokes>",
  "<line 2's raw keystrokes>",
  ...
})
VC.unresolved("some_module")   -- present only if a \r reference is unlinked
VC.run({
  ...
})
EOF
```

Each source line becomes one string element in a `VC.run({...})` table;
`VC.run` (defined in `runtime/common.lua`) iterates the table and executes
each element as its own `normal!` call — which is the mechanism behind
the "one line, one `normal!` call" semantics described in
[language.md](language.md#execution-model).

A `\r <module>` reference becomes a `VC.unresolved("<module>")` call, 
standalone between `VC.run` blocks. This is what `vc link` resolves:
for every `VC.unresolved` marker, the linker splices in the referenced
module's own compiled content. Run an object with unresolved markers still
in it directly (skip linking), and `VC.unresolved` produces the "VC Link Error" message.

## Compiler pipeline

`compiler/main.vc` is itself a `.vc` program — written in the language it
compiles, using `\r <stage>` to pull in each stage of the pipeline in
order:

```
\r preprocess
\r segment
\r validation
\r optimise
\r codegen
```

- **`preprocess`** strips comments (`\\ ...`) and empty lines, and trims
  leading whitespace, before anything else runs.
- **`segment`** is intended to normalize each line into exactly one
  logical motion unit (splitting compound lines, keeping an insert-mode
  operation and its content together as one unit). **This stage is
  currently an unimplemented stub** — worth knowing if you're debugging
  something that looks like a segmentation issue; there isn't real logic
  there yet.
- **`validation`** enforces the no-Ex-commands rule and checks `\r`
  directive syntax, reporting `Invalid Compiler Directive`,
  `Invalid Module Name`, or `Moral Failure: Tried to use an Ex command.`
  for the respective violations.
- **`optimise`** currently strips redundant/consecutive `<Esc>` bytes. A
  planned pass to merge consecutive insert-mode append chains (`a..` /
  `a..` → `a....`) is noted in the source as blocked on `segment` actually
  being implemented.
- **`codegen`** does the real transformation: marking unresolved `\r`
  references as `VC.unresolved(...)`, escaping the remaining content for
  embedding as Lua string literals, and wrapping the result in the
  `VC.run({...})` structure described above.

The compiler's own module system is `\r` itself — there's no separate
bootstrapping mechanism for "including files inside the compiler
specifically." It uses the same directive any `.vc` program does.

## Self-hosting: the bootstrap chain

1. **`bootstrap/stage0.lua`** — a minimal compiler, written directly in
   Lua (not the esolang), just capable enough to compile
   `bootstrap/stage1.vc`.
2. **`bootstrap/stage1.vc`** — a minimal *subset* compiler, written in the
   esolang itself, compiled by `stage0.lua` into `stage1.vcx`.
2. That bootstrap `compiler.vcx` is then used to compile the linker
   (`linker/linker.vc`), and finally to compile the *real* compiler
   (`compiler/main.vc` and its stages), overwriting the minimal bootstrap
   compiler with the full one.

## Runtime

Three Lua files, all under `toolchain/` in an install:

- **`common.lua`** — shared setup: disables swapfile/autoindent/wrapscan,
  defines `VC.init` (resets the `c`/`x`/`y`/`z` registers — see
  [language.md](language.md#registers)), `VC.exec` (runs one line via
  `normal!`, tolerating a small set of harmless Vim errors like a failed
  search), `VC.unresolved`, and the `VimLeave` autocmd that serializes the
  final register state to the path in `VC_RESULT` as JSON
  (`exit_code`/`stdout`/`stderr`) for `vc` to read back.
- **`runtime.lua`** — sources `common.lua`, then defines `VC.run` to run
  straight through a program's instructions with no interaction. Used by
  `vc run`, `vc compile`, `vc build`, `vc link`.
- **`debugger.lua`** — sources `common.lua`, then defines `VC.run` to
  execute one instruction at a time, echoing each command before running
  it and waiting for a keypress (continue / run-to-end / abort) between
  steps. Used by `vc debug`, with Neovim's real UI attached to the
  terminal instead of running headless.

## Install layout

```
<install>/
├── vc                  # the CLI binary
├── toolchain/
│   ├── runtime.lua
│   ├── debugger.lua
│   ├── common.lua
│   ├── compiler.vcx    # built artifact, not source
│   └── linker.vcx      # built artifact, not source
└── stdlib/
    └── *.vco           # built artifacts
```

`vc` resolves this layout relative to its own binary location (or
`VC_ROOT`, if set) — see [cli.md](cli.md#toolchain-resolution). A dev
build (`scripts/build.sh dev` / `make build`) symlinks `toolchain/`'s Lua
files back to `runtime/` in the repo, so edits there are picked up without
rebuilding; a release build (`scripts/build.sh release` / `make release`)
copies real files instead, so the install has no dependency on the source
checkout still existing.
