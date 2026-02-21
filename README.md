# Box - Minimalist project-local toolbox

[![Deploy Documentation](https://github.com/sebakri/box/actions/workflows/docs.yml/badge.svg)](https://sebakri.github.io/box/)
[![Release](https://github.com/sebakri/box/actions/workflows/release.yml/badge.svg)](https://github.com/sebakri/box/releases)

Box is a minimalist, project-local toolbox that keeps your development tools, binaries, and environment variables neatly packed and isolated within your project. It ensures a reproducible environment without cluttering your global system.

## Why use Box?

While there are many tool managers (like `asdf`, `mise`, or `aqua`), Box is designed for simplicity and project-level isolation:

- **No Magic/No Registry**: Box doesn't need a central database of "supported" tools. If it can be installed via `go install`, `npm`, `cargo`, `uv`, or a shell script, Box can manage it.
- **True Project Isolation**: Everything—binaries, caches, and metadata—lives inside your project's `.box` folder. Deleting the folder completely removes the tools.
- **Zero Dependencies**: Box is a single Go binary. You don't need Nix, a plugin system, or a complex runtime to get started.
- **Transparent Wrapper**: It doesn't replace your package managers; it coordinates them to keep your workspace clean.

## Documentation

Full documentation is available at [https://sebakri.github.io/box/](https://sebakri.github.io/box/)

## Quick Start

1.  **Configure**: Create a `box.yml` in your project root:
    ```yaml
    tools:
      - type: go
        source: github.com/go-task/task/v3/cmd/task@latest
      - type: cargo
        source: jj-cli
        args:
          - --strategies
          - crate-meta-data
    env:
      APP_DEBUG: "true"
    ```
2.  **Install**: Run `box install`.
3.  **Setup Shell (Optional)**: Run `box generate direnv` if you use `direnv`.
4.  **Run**: Run `box run <tool>` or use `direnv`.

## Features

- **Project-Local Tools**: Installs tools into a local `.box/bin` directory.
- **Environment Variables**: Define project-specific environment variables in `box.yml`. `BOX_DIR` and `BOX_BIN_DIR` are automatically provided.
- **No Root Required**: Leverages user-space package managers (Go, npm, Cargo, uv, gem) or custom shell scripts.
- **Declarative Configuration**: Defined in `box.yml`.
- **Manual or Automatic PATH**: Use `box run` or generate a `.envrc` for `direnv`.

## Installation

Install using curl:

```bash
curl -fsSL https://raw.githubusercontent.com/sebakri/box/main/scripts/install.sh | sh
```

Or download the binary for your platform from the [latest releases](https://github.com/sebakri/box/releases).

## Commands

- `box install`: Installs tools defined in `box.yml`.
- `box list`: Lists installed tools and their binaries.
- `box run <command>`: Executes a binary from the local `.box/bin` directory.
- `box env`: Displays the merged list of environment variables.
- `box generate direnv`: Generates a `.envrc` file for `direnv` integration.
- `box doctor`: Checks if the host runtimes are installed.

## Development

Build with Task:

```bash
task build
```

Run tests:

```bash
task test
```
