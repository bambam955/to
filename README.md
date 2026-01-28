# to

Go TO any directory instantly. A modern directory navigation tool with a Go backend and bash wrapper.

## Quick Start

```bash
just install
source ~/.local/bin/to.sh
to reg work /path/to/work
to work  # Changes directory instantly
```

See [CLI Reference](docs/CLI.md) for full command documentation.

## Features

- **Instant navigation**: Jump to registered directories with a single command
- **Multiple aliases**: Register multiple names for the same directory
- **Smart management**: List, clean, and expand aliases
- **Simple protocol**: Clear error messages and straightforward communication
- **XDG-compliant**: Configuration stored in `~/.config/to/`

## Installation

```bash
just install
```

Then source the wrapper in your shell configuration (`.bashrc`, `.zshrc`, etc.):
```bash
source ~/.local/bin/to.sh
```

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Development

```bash
just build      # Build the binary
just test       # Run tests
just clean      # Clean build artifacts
just help       # Show all commands
```

## Architecture

- **Backend**: Go binary with database operations and validation
- **Wrapper**: Bash function for shell integration
- **Database**: JSON format at `~/.config/to/database.json`
- **Protocol**: Text-based with `[to] <path>` for navigation, `error: ` for errors

See [PROTOCOL.md](PROTOCOL.md) for protocol details.
