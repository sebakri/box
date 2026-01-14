# ETC - Ephemeral Tool Configuration

[![Deploy Documentation](https://github.com/sebakri/etc/actions/workflows/docs.yml/badge.svg)](https://sebakri.github.io/etc/)
[![Release](https://github.com/sebakri/etc/actions/workflows/release.yml/badge.svg)](https://github.com/sebakri/etc/releases)

ETC is a platform-independent tools manager written in Go. It allows project-specific tool installation without root permissions, keeping your system clean and your project reproducible.

## Documentation

Full documentation is available at [https://sebakri.github.io/etc/](https://sebakri.github.io/etc/)

## Quick Start

1.  **Configure**: Create an `etc.yml` in your project root:
    ```yaml
    tools:
      - name: task
        type: go
        source: github.com/go-task/task/v3/cmd/task@latest
    env:
      APP_DEBUG: "true"
    ```
2.  **Install**: Run `etc install`.
3.  **Run**: Run `etc run <tool>` or use `direnv`.

## Features

- **Project-Local Tools**: Installs tools into a local `.etc/bin` directory.
- **Environment Variables**: Define project-specific environment variables in `etc.yml`.
- **No Root Required**: Leverages user-space package managers (Go, npm, Cargo, uv).
- **Declarative Configuration**: Defined in `etc.yml`.
- **Automatic PATH & Env**: Integrated with `direnv` via `.envrc`.

## Installation

Download the binary for your platform from the [latest releases](https://github.com/sebakri/etc/releases).

## Development

Build with Task:
```bash
task build
```

Run tests:
```bash
task test
```