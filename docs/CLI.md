# CLI Reference

## Installation

```bash
just install
source ~/.local/bin/to.sh
```

## Commands

### Navigate to an alias
```bash
to <alias>
```

Changes directory to the registered alias. Outputs `[to] <path>` to stdout on success.

### Register an alias
```bash
to reg <alias> <directory>
```

Register a new alias pointing to a directory.

### Unregister an alias
```bash
to unreg <alias>
```

Remove an alias from the database.

### List all aliases
```bash
to list
```

Display all registered aliases and their paths.

### Clean invalid aliases
```bash
to clean
```

Remove aliases pointing to directories that no longer exist.

### Expand an alias
```bash
to exp <alias>
```

Output the absolute path for an alias (useful for scripts).

## Error Handling

All errors print to stderr with `error: ` prefix and exit with code 1.

Success returns exit code 0.

## Build & Development

```bash
just build      # Build the binary
just test       # Run tests
just clean      # Clean build artifacts
just dev        # Full dev build (test + build)
just help       # Show all commands
```
