# just recipe file for the 'to' directory navigation tool

# Display help information
default:
    @just --list --unsorted

# --------------- USER COMMANDS --------------- #

# Install the backend and all shell wrappers to the local bin dir
install shell="bash": build
    @mkdir -p "${TO_INSTALL_DIR:-$HOME/.local/bin}"
    @mkdir -p "$HOME/.config/to"
    cp bin/to-backend "${TO_INSTALL_DIR:-$HOME/.local/bin}/to-backend"
    cp src/wrappers/to.bash "${TO_INSTALL_DIR:-$HOME/.local/bin}/"
    cp src/wrappers/to.zsh "${TO_INSTALL_DIR:-$HOME/.local/bin}/"
    cp src/wrappers/to.fish "${TO_INSTALL_DIR:-$HOME/.local/bin}/"
    @chmod +x "${TO_INSTALL_DIR:-$HOME/.local/bin}/to-backend"
    @chmod +x "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.bash"
    @chmod +x "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.zsh"
    @chmod +x "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.fish"
    @# Record the canonical install dir in the same TOML shape as pkg/config.Config.
    @install_dir="$(cd "${TO_INSTALL_DIR:-$HOME/.local/bin}" && pwd -P)"; install_dir_toml="$(printf '%s' "${install_dir}" | sed 's/\\/\\\\/g; s/"/\\"/g')"; printf 'install_dir = "%s"\n' "${install_dir_toml}" > "$HOME/.config/to/config.toml"
    @echo "Add the following to your shell configuration:"
    @# Resolve relative install dirs before printing the shell rc snippet.
    @echo "  source $(cd "${TO_INSTALL_DIR:-$HOME/.local/bin}" && pwd -P)/to.{{ shell }}"

# Build the Go backend binary
build:
    cd src/backend && go build -ldflags "-X main.buildVersion=${TO_VERSION:-dev}" -o ../../bin/to-backend ./cmd

# Rebuild and reinstall all components
upgrade shell="bash": (install shell)

# Remove installed files and all configuration/data
purge:
    rm -f "${TO_INSTALL_DIR:-$HOME/.local/bin}/to-backend"
    rm -f "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.bash"
    rm -f "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.zsh"
    rm -f "${TO_INSTALL_DIR:-$HOME/.local/bin}/to.fish"
    @echo "Warning: removing all configuration and data from ~/.config/to/"
    rm -rf ~/.config/to/

# Generate the checked-in changelog from the current repository state
gen-changelog:
    bash ci/generate-changelog.sh --output CHANGELOG.md

# Create and open a release PR for the requested semantic version
prep-release version:
    bash ci/prepare-release.sh {{ version }}

# --------------- DEV COMMANDS --------------- #

# Development build (clean, lint, fmt, test, build)
dev: clean lint fmt test build

# Run all tests
test: test-unit (test-integration "all")

# Run unit tests
test-unit:
    cd src/backend && go test -race ./...

# Run integration tests (shell = bash, zsh, fish, or all)
test-integration shell="all": build
    bash tests/integration/runner.sh {{ shell }}

# Clean build artifacts
clean:
    @rm -rf bin/
    go clean

# Lint all source code
lint:
    cd src/backend && go vet ./...
    @command -v shellcheck >/dev/null
    shellcheck --enable=all --shell=bash src/wrappers/*.bash tests/**/*.bash tests/**/*.sh ci/*.sh
    shellcheck --enable=all --shell=sh install.sh

# Format all Go source code
fmt:
    cd src/backend && gofmt -w .
    shfmt -i 4 -w src/wrappers/*.bash tests/**/*.bash tests/**/*.sh ci/*.sh install.sh
