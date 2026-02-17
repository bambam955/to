# just recipe file for the 'to' directory navigation tool

# Display help information
default:
    @just --list --unsorted

# --------------- USER COMMANDS --------------- #

# Install the backend and wrapper to the local bin dir
install shell="bash": build
    @mkdir -p ~/.local/bin
    cp bin/to-backend ~/.local/bin/to-backend
    cp wrappers/to.{{ shell }} ~/.local/bin/
    @chmod +x ~/.local/bin/to-backend
    @echo "Add the following to your shell configuration:"
    @echo "  source ~/.local/bin/to.{{ shell }}"

# Build the Go backend binary
build:
    go build -o bin/to-backend ./cmd/to-backend

# Uninstall both components
uninstall shell="bash":
    rm -f ~/.local/bin/to-backend
    rm -f ~/.local/bin/to.{{ shell }}

# Rebuild and reinstall both components
upgrade shell="bash": (uninstall shell) (install shell)

# Uninstall and remove all configuration and data
purge shell="bash": (uninstall shell)
    @echo "Warning: removing all configuration and data from ~/.config/to/"
    rm -rf ~/.config/to/

# --------------- DEV COMMANDS --------------- #

# Development build (clean, lint, fmt, test, build)
dev: clean lint fmt test build

# Run all tests
test: test-unit (test-integration "all")

# Run unit tests
test-unit:
    go test ./...

# Run integration tests (shell = bash, zsh, fish, or all)
test-integration shell="all": build
    bash tests/integration/runner.sh {{ shell }}

# Clean build artifacts
clean:
    @rm -rf bin/
    go clean

# Lint all source code
lint:
    go vet ./...
    @command -v shellcheck >/dev/null
    shellcheck --enable=all --shell=bash wrappers/*.bash tests/**/*.bash tests/**/*.sh

# Format all Go source code
fmt:
    gofmt -w .
    shfmt -i 4 -w .
