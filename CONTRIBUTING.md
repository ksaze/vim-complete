# Contributing to VimComplete

Thank you for your interest in contributing to VimComplete! Contributions, bug reports, feature requests, and documentation improvements are welcome.

## Getting Started

A [techniques doc](docs/techniques.md) is provided for contributors to familiarise themselves with the techniques/structures used in the source code.

Please review GitHub issues for more concrete problems and [TODO](TODO.md) for more general directions.

### Prerequisites

You will need to install and configure the following to build and work on VimComplete as a contributor:
- Git
- Neovim
- Go
- Make
- golangci-lint 
- Stylua
- luacheck
- pre-commit

## Local Development

1. Fork the repository.

2. Clone your fork:
```
git clone https://github.com/<your_github_username>/vim-complete.git
cd vim-complete
```

3. Add the upstream remote so you can pull in future releases:
```
git remote add upstream https://github.com/ksaze/vim-complete.git
```

4. Create a new branch:
```
git checkout -b feature/your-feature-name
```

## Building

Use the build script, or `make build`, to build the entire project:
```
./scripts/build.sh dev
```

For rebuilding specific components of the project, use the provided make commands:
```
make help
```

## Running Tests

The test suite assumes a build already exists. Run the script directly or with 
the make command.

```
make test
```

```
./scripts/run_test_suite.sh
```

This repo also includes a `.pre-commit-config.yaml` with a hook that runs tests before each commit.

## Community Channels

If you get stuck or have questions, open a GitHub Discussions thread.

## Issues

If you find a bug or have a concrete feature suggestion, please create an issue.

- Ensure that a similar issue doesn't exist already
- If the feature suggestion significantly changes the architecture, go through GitHub Discussions before creating an issue
- For bug reports, please include a clear description of the problem along with steps to reproduce it

## Pull Requests

After making your changes, open a pull request against the `main` branch.

- Link your PR to an existing issue, if present
- If you're fixing an issue, make sure someone hasn't already created a PR fixing the same issue
- As with issues, a pull request proposing a significant architectural change should go through GitHub Discussions first
- Prior to submitting your PR, ensure that all existing tests pass using `make test`
