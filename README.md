# VimComplete
An esolang whose instruction set is Neovim's motions including search commands, but not any Ex commands.

Each line of the vc source file is run using a normal! command on the runtime. Thus, all lines come with an implicit escape back to normal mode. 
Additionally, because the execution model of normal! command is similar to macros, if a motion fails, then everything afterwards isn't executed.

A consequence of this on the language semantics is that these two versions aren't equivalent:
```
\\ no-op, no delete occurs because j fails
Gjdd

\\ deletes the last line
Gj
dd
```

Refer [language docs](./docs/language.md) for the complete VimComplete language specification.

## Building
```bash
git clone https://github.com/ksaze/vim-complete
cd vim-complete
chmod +x scripts/build.sh
./scripts/build.sh release
```

### Dependencies (tested with)
- nvim (0.12.6)
- go (1.26.4)

If you're planning on contributing, then prefer `build.sh dev` which creates symlinks for runtime files instead of copying them to build directory.

### Rebuilding
`build.sh` script should only be used for a clean install. For dev work, affecting only one component, make commands like make cli, make stdlib are provided--run make help for a list of all possible commands.

## CLI

| Command | What it does |
|---|---|
| `vc compile <file> [-o path] [-f]` | Compile one source file to an object file (default: `vc-build/<name>.vco`). |
| `vc build <dir> [-m main.vc] [-o out] [-f]` | Compile every source file in `dir`. With `-m`, link into a single executable at `out`. Without `-m`, link each file individually against its siblings and install the objects into `out` — for building a library/stdlib-style package. |
| `vc link <main-file> [-o out]` | Link a main object against every discovered local + stdlib object. |
| `vc run <target> [args...]` | Run a compiled `.vcx`/`.vco` on the runtime, headless. |
| `vc debug <target> [args...]` | Same, but with an interactive Neovim UI and step-through control. |
| `vc modules` | List the local and stdlib modules currently visible to the linker. |

All build outputs are skipped if already up to date (compared by mtime);
pass `-f`/`--force` to rebuild anyway. Set `VC_ROOT` to point the CLI at a
toolchain install other than the one next to the `vc` binary itself.

## Testing
The test suite includes self-hosted unit tests for stdlib and compiler and bash scripts for bootstrap and fixed point tests.

```bash
make test
# or directly:
./scripts/build.sh dev && ./scripts/run_test_suite.sh
```

The toolchain currently doesn't provide any mechanism for writing unit tests easily. The same goes on the CLI side--the script currently has to compile and run each test manually. 
