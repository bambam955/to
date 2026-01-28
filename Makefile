.PHONY: build install install-bash uninstall test clean

BINARY_NAME=to-backend
BINARY_PATH=cmd/$(BINARY_NAME)
WRAPPER_NAME=to.sh
WRAPPER_PATH=internal/bash/$(WRAPPER_NAME)

# Build the Go backend binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./$(BINARY_PATH)

# Install the backend binary to ~/.local/bin/
install-backend: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin/..."
	@mkdir -p ~/.local/bin
	@cp bin/$(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@chmod +x ~/.local/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) successfully"

# Install the bash wrapper to ~/.local/bin/
install-bash:
	@echo "Installing $(WRAPPER_NAME) to ~/.local/bin/..."
	@mkdir -p ~/.local/bin
	@cp $(WRAPPER_PATH) ~/.local/bin/$(WRAPPER_NAME)
	@chmod +x ~/.local/bin/$(WRAPPER_NAME)
	@echo "Installed $(WRAPPER_NAME) successfully"
	@echo "Add the following to your shell configuration (.bashrc, .zshrc, etc.):"
	@echo "  source ~/.local/bin/$(WRAPPER_NAME)"

# Install both backend and wrapper
install: install-backend install-bash

# Uninstall both components
uninstall:
	@echo "Uninstalling to tool..."
	@rm -f ~/.local/bin/$(BINARY_NAME)
	@rm -f ~/.local/bin/$(WRAPPER_NAME)
	@echo "Uninstalled to tool"

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	@rm -rf bin/
	go clean

# Development build with verbose output
dev: clean test build
	@echo "Development build complete"
