# Changelog

All notable changes to this project are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). This project is
pre-1.0; the CLI surface, module format, and object file layout may still
change between minor versions.

## [Unreleased]

## [0.1.0] - beta

Initial public beta.

### Added
- Self-hosted compiler (`compiler/`), bootstrapped from a minimal Lua
  seed compiler (`bootstrap/stage0.lua` → `bootstrap/stage1.vc`).
- Linker (`linker/`) with local/stdlib module resolution.
- `vc` CLI (`compile`, `build`, `link`, `run`, `debug`, `modules`), with
  staleness-aware recompilation and a `-force` override.
- `vc build` library mode: compiles and links a whole directory of
  modules against each other, for building a standalone stdlib-style
  package rather than a single executable.
- Interactive step-through debugger (`vc debug`) with a real Neovim UI
  attached to the terminal.
- Standard library covering register/pattern/text manipulation
  primitives (`add`, `subtract`, `trim`, `search_and_replace`,
  `search_and_transform`, `copy_register`, `assert_present`,
  `assert_absent`, and `file_diff_assert`).
- `examples/game_of_life`: a complete, documented example program.
- `tests/`: stdlib and compiler test suites, run via
  `scripts/run_test_suite.sh`.
- `tools/nvim/`: syntax highlighting and filetype support for editing
  `.vc` source in Neovim.

[Unreleased]: https://github.com/ksaze/vim-complete/compare/v0.1.0b...HEAD
[0.1.0b]: https://github.com/ksaze/vim-complete/releases/tag/v0.1.0b
