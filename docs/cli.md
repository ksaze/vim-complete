# `vc` CLI reference

`vc` is the wrapper around the compiler, linker, and runtime — the
language itself has no way to read the file system, take command-line
arguments, or produce output on its own, so `vc` provides all of that from
the outside.

## Toolchain resolution

`vc` locates its own compiler, linker, and runtime files relative to
wherever the `vc` binary itself lives — normally `<install>/toolchain/`
and `<install>/stdlib/`, sitting next to the `vc` binary (see
[architecture.md](architecture.md#install-layout)). Set `VC_ROOT` to point
at a different install instead:

```bash
VC_ROOT=/some/other/build vc compile foo.vc
```

## `vc compile <file>`

Compiles a single source file to an object file.

| Flag | Default | Meaning |
|---|---|---|
| `-o`, `--out` | `vc-build/<name>.vco` | Output object file path. |
| `-f`, `--force` | off | Recompile even if the object is already up to date. |

Staleness is checked by mtime: if the object already exists and is newer
than the source, `compile` does nothing and reports `up to date`. This is
what lets `vc build` re-run cheaply over a large directory.

```bash
vc compile stdlib/add.vc                 # -> vc-build/add.vco
vc compile stdlib/add.vc -o out/add.vco  # explicit output path
vc compile stdlib/add.vc -f              # force recompile
```

## `vc build <dir>`

Compiles every source file in `dir` (non-recursive), then — depending on
whether `-m` is given — either links a single executable or builds a
library.

| Flag | Default | Meaning |
|---|---|---|
| `-m`, `--main` | *(none)* | Entry file, relative to `dir`. Presence of this flag selects executable mode. |
| `-o`, `--out` | `a.vcx` | Executable path (`-m` given) or install directory (library mode). |
| `-f`, `--force` | off | Recompile/relink everything even if outputs are up to date. |

**Executable mode** (`-m` given): every source file in `dir` is compiled
into `vc-build/` as an intermediate, then the object for `-m`'s file is
linked — pulling in every other object in `vc-build/` plus the stdlib —
into a single executable at `-o`.

```bash
vc build examples/game_of_life/src -m src/main.vc -o gol.vcx
```

**Library mode** (`-m` omitted): every source file is still compiled into
`vc-build/` first (so sibling modules can depend on each other), but
instead of one link step there's one *per file* — each module is linked
against its siblings and the stdlib individually, and each linked result
is installed into `-o` as its own object. This is how the stdlib itself
gets built.

```bash
vc build stdlib -o build/stdlib
```

## `vc link <main-file>`

Links a single already-compiled main object against every object
currently discoverable — everything in `vc-build/`, plus everything in the
stdlib directory. This is the lower-level primitive `vc build` uses
internally; call it directly if you've compiled things by hand with
`vc compile` and just want the link step without recompiling.

| Flag | Default | Meaning |
|---|---|---|
| `-o`, `--out` | `a.vcx` | Output executable path. |

```bash
vc compile src/main.vc
vc link vc-build/main.vco -o app.vcx
```

## `vc run <target> [args...]`

Runs a compiled object or executable on the runtime, headless. Any
trailing arguments are passed through to the program's `c` register (see
[language.md](language.md#registers)).

```bash
vc run gol.vcx input.txt 3
```

An object that still has unresolved `\r` references (i.e. you ran a
`.vco` that was never linked) will fail with a link error at the point the
unresolved reference is reached — see
[language.md](language.md#modules-the-r-directive).

## `vc debug <target> [args...]`

Same as `run`, but with a real, interactive Neovim UI attached to your
terminal instead of running headless. Steps through the program one
instruction at a time, showing each command before it executes and the
buffer state after. Takes the same trailing arguments as `run`.

```bash
vc debug gol.vcx input.txt 1
```

Requires the install's `debugger.lua` to be present — this is part of the
standard toolchain layout, not something you need to set up separately.

## `vc modules`

Lists every object currently visible to the linker: local objects in
`vc-build/`, and stdlib objects, each labeled accordingly. Useful for
sanity-checking what a `vc link`/`vc build -m` invocation will actually
pull in, without needing to inspect the (temporary, auto-deleted) module
list the linker itself is handed internally.

```bash
vc modules
# local: /path/to/vc-build/main.vco
# local: /path/to/vc-build/utils.vco
# std:   /path/to/build/stdlib/add.vco
# ...
```

## Exit codes and output

Every command forwards the compiled program's own stdout/stderr (its `y`
and `z` registers) to the real stdout/stderr, and its exit code (`x`
register) becomes `vc`'s own exit code. A nonzero exit from `vc compile`,
`vc build`, or `vc link` means the *toolchain* failed (compiler/linker
error) — a nonzero exit from `vc run`/`vc debug` means the *program you
ran* set a nonzero exit code, which is a normal, expected outcome for a
program signaling failure, not a `vc` bug.
