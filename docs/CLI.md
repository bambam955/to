# CLI Reference

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | sh
source ~/.local/bin/to.bash
```

Or build from source:

```bash
just install
source ~/.local/bin/to.bash
```

Both install paths copy `to-backend` plus all three wrappers: `to.bash`, `to.zsh`, and `to.fish`.
Set `TO_INSTALL_DIR` to override the install location for `just install`; the resolved install path is recorded in `~/.config/to/config.toml` so `to -U` / `to -P` remove the same location later.

## Commands

### Navigate to an alias

```bash
to <alias>
```

Changes directory to the registered alias. Navigation control is emitted on fd 3 (`NAV <path>`) for shell wrappers; normal stdout is reserved for user-facing output.

### Register an alias

```bash
to --reg <alias> <directory>
to -r <alias> <directory>
```

Register a new alias pointing to a directory. Relative paths are resolved to absolute. Warns if another alias already points to the same directory.

### Unregister an alias

```bash
to --unreg <alias>
to -u <alias>
```

Remove an alias from the database.

### List all aliases

```bash
to --list
to -l
```

Display all registered aliases and their paths in aligned columns, sorted alphabetically.

### Clean invalid aliases

```bash
to --clean
to -c
```

Remove aliases pointing to directories that no longer exist. Shows each removed alias and a summary of how many were cleaned.

### Expand an alias

```bash
to --exp <alias>
to -e <alias>
```

Output the absolute path for an alias (useful for scripts). Outputs only the path with no extra formatting.

### Uninstall the tool

```bash
to --uninstall
to -U
```

Remove the backend and all installed wrappers. Prompts for confirmation before deleting files.

### Purge install data

```bash
to --purge
to -P
```

Remove the backend, all installed wrappers, and the default configuration directory at `~/.config/to/` along with `config.toml`.

### Inspect the installed version

```bash
to --version
to -v
```

Print the installed semantic version as `TO <version>`. Development builds fall back to `TO dev`, and release builds are expected to match the semantic tag used to compile the binary.

Release tags use bare Semantic Versioning identifiers such as `1.2.3`, and tagged builds should pass the same value through `TO_VERSION` so the binary and any packaged release artifacts stay aligned.

## Error Handling

All errors print to stderr with `error: ` prefix and exit with code 1.

Success returns exit code 0.

## Build & Development

```bash
TO_VERSION=1.2.3 just build   # Build a tagged release binary
just build                    # Build a local development binary
just test       # Run tests
just clean      # Clean build artifacts
just dev        # Full dev build (test + build)
just            # Show all commands
```
