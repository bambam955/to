# just recipe file for the 'to' directory navigation tool

# Display help information
default:
    @just --list

# --------------- USER COMMANDS --------------- #

# Install the backend and wrapper to the local bin dir
install: build
    @mkdir -p ~/.local/bin
    cp -t ~/.local/bin/ bin/to-backend shell/to.sh
    @chmod +x ~/.local/bin/to-backend ~/.local/bin/to.sh
    @echo "Add the following to your shell configuration (.bashrc, .zshrc, etc.):"
    @echo "  source ~/.local/bin/to.sh"

# Build the Go backend binary
build:
    go build -o bin/to-backend ./cmd/to-backend

# Uninstall both components
uninstall:
    @rm -f ~/.local/bin/to-backend
    @rm -f ~/.local/bin/to.sh
    @echo "Remember to remove the following from your shell configuration (.bashrc, .zshrc, etc.):"
    @echo "  source ~/.local/bin/to.sh"

# Rebuild and reinstall both components
upgrade: uninstall install

# Uninstall and remove all configuration and data
purge: uninstall
    @echo "Warning: removing all configuration and data from ~/.config/to/"
    @rm -rf ~/.config/to/

# --------------- DEV COMMANDS --------------- #

# Development build (clean, fmt, lint, test, build)
dev: clean fmt lint test build

# Run tests
test:
    go test ./...

# Clean build artifacts
clean:
    @rm -rf bin/
    go clean

# Lint all source code
lint:
    go vet ./...
    @command -v shellcheck >/dev/null
    shellcheck --enable=all --shell=bash shell/*

# Format all Go source code
fmt:
    gofmt -w .
    shfmt -i 4 -w .
