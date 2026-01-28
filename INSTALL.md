# Installation Guide for "to" Tool

## Overview

The "to" tool consists of two components:
- **Go Backend** (`to-backend`): Core binary handling database operations
- **Bash Wrapper** (`to.sh`): Shell integration for directory navigation

## Installation Methods

### Using Make (Recommended)

#### Install Both Components
```bash
make install
```

This will:
1. Build the Go backend binary
2. Install `to-backend` to `~/.local/bin/`
3. Install `to.sh` to `~/.local/bin/`
4. Display instructions for sourcing the wrapper in your shell configuration

#### Install Only Backend
```bash
make install-backend
```

#### Install Only Wrapper
```bash
make install-bash
```

### Manual Installation

1. Build the backend:
```bash
go build -o bin/to-backend ./cmd/to-backend
```

2. Install to `~/.local/bin/`:
```bash
mkdir -p ~/.local/bin
cp bin/to-backend ~/.local/bin/to-backend
cp internal/bash/to.sh ~/.local/bin/to.sh
chmod +x ~/.local/bin/to-backend ~/.local/bin/to.sh
```

3. Source the wrapper in your shell configuration (`.bashrc`, `.zshrc`, `.profile`, etc.):
```bash
source ~/.local/bin/to.sh
```

## Uninstallation

```bash
make uninstall
```

Or manually:
```bash
rm ~/.local/bin/to-backend ~/.local/bin/to.sh
```

Remove the `source ~/.local/bin/to.sh` line from your shell configuration.

## Verification

After installation, verify everything is working:

1. Check that `~/.local/bin` is in your PATH:
```bash
echo $PATH | grep -q "$HOME/.local/bin" && echo "PATH is correct" || echo "Add ~/.local/bin to PATH"
```

2. Reload your shell configuration or start a new terminal
3. Test the tool:
```bash
to --help
```

## PATH Configuration

The installation assumes `~/.local/bin` is in your PATH. If it's not, add this to your shell configuration:

For bash (`.bashrc`):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

For zsh (`.zshrc`):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Architecture Details

### Component Locations
- **Binary**: `~/.local/bin/to-backend`
- **Wrapper**: `~/.local/bin/to.sh` (sourced into shell)

### Database Location
- **Database**: `~/.config/to/database.json` (XDG-compliant, created on first use)
- **Override**: Set `TO_DB` environment variable to use a custom database location

### Communication Protocol

The bash wrapper communicates with the backend using a simple protocol:

#### Navigation Success
When navigation succeeds, the backend outputs:
```
[to] /path/to/directory
```

The wrapper parses this response, extracts the path, and executes `cd` to change directory.

#### Errors
All errors are printed to stderr with the prefix `error:`:
```
error: alias not found: myalias
```

The wrapper forwards these messages unchanged to the user.

#### Exit Codes
- `0`: Command succeeded
- `1`: Command failed (any error condition)

## Troubleshooting

### Command Not Found
If you get "command not found: to":
1. Verify installation: `ls -la ~/.local/bin/to-backend`
2. Check PATH: `echo $PATH | grep -q ~/.local/bin && echo OK || echo FAIL`
3. Reload shell: `source ~/.bashrc` or start a new terminal

### Database Permission Error
If you get permission errors:
1. Check database directory: `ls -la ~/.config/to/`
2. Check database file: `ls -la ~/.config/to/database.json`
3. Fix permissions if needed: `chmod 644 ~/.config/to/database.json`

### Backend Binary Permission Error
If you get permission errors running the backend:
```bash
chmod +x ~/.local/bin/to-backend
```

## Build and Development

### Build Only
```bash
make build
```

### Run Tests
```bash
make test
```

### Clean Build Artifacts
```bash
make clean
```

### Development Build (Clean, Test, Build)
```bash
make dev
```

## Dependencies

The Go backend requires:
- Go 1.21 or later
- Standard library packages
- `github.com/spf13/cobra` for CLI framework

The bash wrapper requires:
- Bash 4.0 or later
- Standard Unix utilities (`mkdir`, etc.)

No additional system dependencies are required for runtime.
