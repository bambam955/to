# just recipe file for the 'to' directory navigation tool

# Display help information
default:
    @just --list

# --------------- USER COMMANDS --------------- #

# Install the backend and wrapper to the local bin dir
install shell="bash": build
    @mkdir -p ~/.local/bin
    cp bin/to-backend ~/.local/bin/to-backend
    cp shell/to.{{ if shell == "fish" { "fish" } else if shell == "zsh" { "zsh" } else { "sh" } }} ~/.local/bin/
    @chmod +x ~/.local/bin/to-backend
    @echo "Add the following to your shell configuration:"
    @echo "  {{ if shell == "fish" { "source ~/.local/bin/to.fish" } else if shell == "zsh" { "source ~/.local/bin/to.zsh" } else { "source ~/.local/bin/to.sh" } }}"

# Build the Go backend binary
build:
    go build -o bin/to-backend ./cmd/to-backend

# Uninstall both components
uninstall shell="bash":
    @rm -f ~/.local/bin/to-backend
    @rm -f ~/.local/bin/to.{{ if shell == "fish" { "fish" } else if shell == "zsh" { "zsh" } else { "sh" } }}
    @echo "Remember to remove the source line from your shell configuration."

# Rebuild and reinstall both components
upgrade shell="bash": (uninstall shell) (install shell)

# Uninstall and remove all configuration and data
purge shell="bash": (uninstall shell)
    @echo "Warning: removing all configuration and data from ~/.config/to/"
    @rm -rf ~/.config/to/

# --------------- DEV COMMANDS --------------- #

# Development build (clean, fmt, lint, test, build)
dev: clean fmt lint test build

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
