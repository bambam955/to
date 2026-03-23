# to

**Go TO any directory instantly.** A lightweight CLI tool for bookmark-style directory navigation, built with Go and bash.

## Quick Start

```bash
just install
source ~/.local/bin/to.bash

to --reg work ~/projects/work   # Register an alias
to work                      # Jump there instantly
```

## Installation

### Prerequisites

- **Go** 1.25+
- **just** (command runner)
- **bash** (or zsh)

### Install

```bash
git clone https://github.com/bambam955/to.git
cd to
just install
```

This builds the `to-backend` binary and copies it along with `to.bash` to `~/.local/bin/`. Make sure `~/.local/bin` is on your `PATH`.

### Make it persistent

Add this line to your `.bashrc` or `.zshrc`:

```bash
source ~/.local/bin/to.bash
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
just gen-changelog            # Regenerate CHANGELOG.md for the current repo state
just prep-release 1.2.3       # Create release/1.2.3, regenerate CHANGELOG.md, and open the release PR
just test       # Run all tests
just fmt        # Format Go and shell code
just lint       # go vet + shellcheck
just dev        # Full cycle: clean → fmt → lint → test → build
just upgrade    # Rebuild and reinstall
```

### Release Changelog Workflow

Release preparation requires:

- `git-cliff` on your `PATH`
- GitHub CLI (`gh`) on your `PATH`
- an authenticated GitHub CLI session (`gh auth login`)

Run `just prep-release <version>` from a clean `main` checkout. The command:

1. creates `release/<version>`
2. regenerates `CHANGELOG.md` for that exact release version
3. commits the changelog update
4. pushes the release branch
5. opens the release PR against `main`

After the release PR is merged, cut the release by pushing the matching bare tag:

```bash
git tag 1.2.3
git push origin 1.2.3
```

The tag-triggered publishing workflow defined in spec `005` is expected to use the matching `CHANGELOG.md` section as the GitHub Release notes source of truth, with deployment approval if configured.

Only user-facing Conventional Commit types are included in generated release notes. `spec:` commits are intentionally excluded, so use `feat`, `fix`, `docs`, `refactor`, `test`, and `chore` for releasable work.

## Uninstalling

```bash
just uninstall  # Remove binary and shell wrapper
just purge      # Also remove ~/.config/to/ data
```

Remember to remove the `source ~/.local/bin/to.bash` line from your shell config.

## License

[MIT](LICENSE) — Bennett Moore
