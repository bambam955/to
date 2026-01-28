# just recipe file for the 'to' directory navigation tool

set shell := ["bash", "-c"]

# Display help information
help:
    @just --list

# Build the Go backend binary
build:
    @echo "Building to-backend..."
    go build -o bin/to-backend ./cmd/to-backend

# Install the backend binary to ~/.local/bin/
install-backend: build
    @echo "Installing to-backend to ~/.local/bin/..."
    @mkdir -p ~/.local/bin
    @cp bin/to-backend ~/.local/bin/to-backend
    @chmod +x ~/.local/bin/to-backend
    @echo "Installed to-backend successfully"

# Install the bash wrapper to ~/.local/bin/
install-bash:
    @echo "Installing to.sh to ~/.local/bin/..."
    @mkdir -p ~/.local/bin
    @cp internal/bash/to.sh ~/.local/bin/to.sh
    @chmod +x ~/.local/bin/to.sh
    @echo "Installed to.sh successfully"
    @echo "Add the following to your shell configuration (.bashrc, .zshrc, etc.):"
    @echo "  source ~/.local/bin/to.sh"

# Install both backend and wrapper
install: install-backend install-bash

# Uninstall both components
uninstall:
    @echo "Uninstalling to tool..."
    @rm -f ~/.local/bin/to-backend
    @rm -f ~/.local/bin/to.sh
    @echo "Uninstalled to tool"

# Run tests
test:
    go test -v ./...

# Clean build artifacts
clean:
    @rm -rf bin/
    go clean

# Development build (clean, test, build)
dev: clean test build
    @echo "Development build complete"
