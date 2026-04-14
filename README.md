# to

**Go TO any directory instantly.** A lightweight CLI tool for bookmark-style directory navigation, built with Go and bash.

## Quick Start

```bash
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | sh
source ~/.local/bin/to.bash

to --reg work ~/projects/work   # Register an alias
to work                         # Jump there instantly
```

## Installation

### Curl installer

```bash
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | sh
```

The installer downloads the latest GitHub Release for Linux (`amd64` and `arm64`), verifies its SHA-256 checksum, and installs `to-backend` plus all shell wrappers to `~/.local/bin` by default.

Override the version or target directory when needed:

```bash
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | TO_INSTALL_VERSION=0.1.0 sh
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | TO_INSTALL_DIR=./bin sh
```

`TO_INSTALL_DIR` accepts relative or absolute paths. The installer resolves the final path before copying files.

### Build from source

Prerequisites:

- **Go** 1.25+
- **just** (command runner)
- **bash** (or zsh / fish)

```bash
git clone https://github.com/bambam955/to.git
cd to
just install
```

This builds the `to-backend` binary and copies it along with the selected shell wrapper to `~/.local/bin/` by default. Set `TO_INSTALL_DIR` to override the target directory for `just install` and `just uninstall`.

### Make it persistent

Add the wrapper that matches your shell to your startup config:

```bash
source ~/.local/bin/to.bash
source ~/.local/bin/to.zsh
source ~/.local/bin/to.fish
```

## Usage

| Action | Command | Short form |
|---|---|---|
| Navigate | `to <alias>` | — |
| Register | `to --reg <alias> <path>` | `to -r <alias> <path>` |
| Unregister | `to --unreg <alias>` | `to -u <alias>` |
| List all | `to --list` | `to -l` |
| Clean stale | `to --clean` | `to -c` |
| Expand path | `to --exp <alias>` | `to -e <alias>` |
| Version | `to --version` | `to -v` |

### Examples

```bash
# Register aliases
to -r proj ~/projects/myapp
to -r docs ~/Documents

# Navigate
to proj                  # cd ~/projects/myapp
to docs                  # cd ~/Documents

# List all registered aliases
to -l

# Get the raw path (useful in scripts)
cp file.txt "$(to -e proj)/src/"

# Remove an alias
to -u docs

# Clean aliases pointing to deleted directories
to -c
```

See [docs/CLI.md](docs/CLI.md) for the full CLI reference.

### Color Output

`to` applies ANSI colors automatically when output is connected to an interactive terminal (TTY), including:

- help/usage text
- list output (`to --list`)
- warnings, success messages, and errors

When output is redirected or piped, color is disabled automatically so output remains plain text.

```bash
# Interactive terminal: colored output
to --help
to --list

# Non-interactive: plain output (no ANSI escapes)
to --list | cat
to --help > help.txt
```

Version output is always plain text so release checks and scripts can parse it reliably.

```bash
to --version
# TO 1.2.3
```

Release builds should be tagged with bare Semantic Versioning identifiers like `1.2.3`, and the same value should be passed to `TO_VERSION` when building the binary.

## Configuration

Alias data is stored in `~/.config/to/database.json` (XDG-compliant).

Override the database location with the `TO_DB` environment variable:

```bash
export TO_DB=~/custom/path/database.json
```

## Architecture

`to` uses a hybrid Go + bash design:

- **`to-backend`** — A Go binary (using cobra) that resolves aliases, manages the database, and writes results to stdout.
- **`to.bash`** — A bash function wrapper that calls `to-backend`, intercepts navigation responses, and performs the actual `cd`.

This split exists because a subprocess cannot change the parent shell's working directory. On successful navigation, the backend emits a control frame on file descriptor 3 (`NAV <path>`), and the shell wrapper parses that control channel to run `cd`.

## Development

```bash
TO_VERSION=1.2.3 just build   # Build a tagged release binary
just build                    # Build a local development binary
just gen-changelog            # Regenerate CHANGELOG.md from mainline or refresh the active release branch
just prep-release 1.2.3       # Create release/1.2.3, regenerate CHANGELOG.md, and open the release PR
just test       # Run all tests
just fmt        # Format Go and shell code
just lint       # go vet + shellcheck
just dev        # Full cycle: clean → fmt → lint → test → build
just upgrade    # Rebuild and reinstall
```

Release workflow details live in [docs/release-workflow.md](docs/release-workflow.md).

## Uninstalling

```bash
just uninstall  # Remove binary and shell wrapper
just purge      # Also remove ~/.config/to/ data
```

`just uninstall` also respects `TO_INSTALL_DIR` when you installed outside `~/.local/bin`.

Remember to remove the `source ~/.local/bin/to.<shell>` line from your shell config.

## License

[MIT](LICENSE) — Bennett Moore
