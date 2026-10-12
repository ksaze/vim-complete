# Language reference

## Instruction set

VimComplete's instruction set is Neovim's normal-mode motions. Searches are allowed
but all Ex commands are rejected by the compiler if detected.

```
Moral Failure: Tried to use an Ex command.
```

## Execution model

Each line of a `.vc` source file is executed by the runtime with a `normal!` call
on the provided file, or an empty buffer. Two semantic consequences follow from
`normal!` being used as the execution engine.

**Every line implicitly returns to Normal mode.** `normal!` always ends in
Normal mode, even if the line left off mid-Insert. You never need a
trailing `<Esc>` to close out a line.

**A failing motion aborts the rest of that line.** `normal!` command in Neovim behaves
like a macro, which means that if a motion fails like trying to `l` at the end of a line,
then all keystrokes after that point in `normal!` command aren't executed.

The implication of this on the language's semantics is that these two programs aren't equivalent.
```
\\ no-op — j fails (already on the last line), so dd never runs
Gjdd

\\ deletes the last line — G and j are one normal! call each,
\\ dd is a separate call that isn't affected by j's failure
Gj
dd
```

The runtime modifies search behaviour for increasing the expressive power of the language.

1. `wrapscan` is disabled for making searches more predictable and offering better isolation
of search space

2. Search failure doesn't interrupt the script like other failures. They are ignored using 
`pcall`. This allows easy implementation of structures like regex conditionals and regex asserts

As mentioned before these changes don't change the computation power of the language. The language 
is still Turing Complete without these modifiers. Search failures can be prevented by adding terms 
that are bound to match while differentiating these from actual matches using cursor position, for example,
trying to match the EOF as the alternate.

This technique is used in stdlib for implementing search and replace and other related macros.

## Comments

A line beginning with `\\` is a comment and produces no instruction:

```
\\ this line does nothing at runtime
```

A line beginning with `\\#` could be treated as documentation in future releases when doc tooling is developed.
The doc line prefix is still a comment so it doesn't affect compilation in any way.

## Modules: the `\r` directive

A line of the form `\r <module>` includes another module's source at that
point. The includes itself doesn't happen at compilation time—only the directive's
syntax is validated there. Resolutoin is deferred to the linker, so that a shared
module is only compiled once rather than re-inlined at every include set.

```
\r delete_pattern      \\ a local module, resolved relative to the current directory
\r std:left_trim       \\ a stdlib module, resolved from the stdlib search path
```

The `std:` prefix selects the standard library; without it, the module is
resolved locally. An explicit `local:` prefix can also be added if required.

## Registers

Four registers carry information between your program and the CLI
around it. They're ordinary Vim registers — read and write them with the
usual register syntax (`"cp`, `"xy`, etc.) — but the runtime treats these
four specially:

| Register | Meaning |
|---|---|
| `c` | Command-line arguments passed to the program |
| `x` | Exit code. Read as a number when the program finishes; `0` unless your program sets it otherwise. |
| `y` | Captured stdout after execution |
| `z` | Captured stderr after execution |

Before program execution starts, all four are reset to their defaults
(`c` = argv, `x` = `"0"`, `y`/`z` = empty).

Registers x, y, and z are polled only after program exits, so all these registers can be freely used
in the program for general usage, provided that they are set to their appropriate values for all exit paths
to prevent garbage values being forwarded to stdout/stderr.
