# PT - Project Tools

PT is a lightweight, platform-independent tools manager for your projects. It allows you to define and install project-specific tools locally without requiring root permissions, keeping your global environment clean and your development setup reproducible.

## Key Features

- **Project-Local Tools**: Tools are installed in `.pt/bin`, isolated from your system.
- **Environment Variables**: Define project-specific variables that are automatically exported.
- **Declarative Setup**: Define all required tools and env vars in a simple `pt.yml` file.
- **Multi-Runtime Support**: Works seamlessly with Go, npm, Cargo, and uv.
- **direnv Integration**: Automatically manages your `PATH` and `ENV` using `.envrc`.
- **Cross-Platform**: Built in Go, supporting Linux, macOS, and Windows.

## Quick Start

### 1. Installation

Download the latest binary for your platform from the [Releases](https://github.com/sebakri/etc/releases) page and place it in your system PATH.

### 2. Configure Your Project

Create a `pt.yml` in your project root:

```yaml
tools:
  - name: task
    type: go
    source: github.com/go-task/task/v3/cmd/task@latest
  - name: ruff
    type: uv
    source: ruff
  - name: jj
    type: cargo
    source: jj-cli
    args:
      - --strategies
      - crate-meta-data
env:
  DEBUG: "true"
  API_URL: "http://localhost:8080"
```

### 3. Install Tools

Run the install command to fetch and install all defined tools:

```bash
pt install
```

### 4. Use Your Tools

You can run tools directly using the `run` command:

```bash
pt run task --version
```

Or, if you use `direnv`, simply `allow` the `.envrc` and use them directly:

```bash
task --version
```

## Commands

- `pt install`: Installs tools defined in `pt.yml` and sets up `.envrc`.
- `pt run <command>`: Executes a binary from the local `.pt/bin` directory.
- `pt doctor`: Checks if the required host runtimes (Go, npm, Cargo, uv) are installed.
- `pt help`: Displays usage information.

## Contributing

We welcome contributions! Please check the [repository](https://github.com/sebakri/etc) for issues and pull requests.
