---
number: 5
status: completed
author: Bennett Moore
creation_date: 2026-02-17
approved_by: Bennett Moore
approval_date: 2026-03-19
---

# Curl Installer Script

## Overview

Create a one-line curl installer so users can install `to` without cloning the repo or having Go installed. The installer downloads prebuilt binaries from GitHub Releases and places them alongside the shell wrapper, similar to how tools like mise and opencode distribute themselves.

Target experience:

```bash
curl -fsSL https://raw.githubusercontent.com/bambam955/to/main/install.sh | bash
```

This requires two pieces:

1. **GitHub Releases with prebuilt binaries** — a CI workflow that cross-compiles the Go backend for supported platforms and uploads tarballs to GitHub Releases on each tagged version.
2. **Installer script** — a POSIX shell script that detects OS/arch, downloads the correct binary + shell wrapper from the latest release, installs them to `~/.local/bin/`, and prints shell config instructions.

## Why It's Needed

- Currently users must clone the repo and have Go 1.25+ installed to build from source
- A curl installer is the standard distribution method for modern CLI tools
- Removes the Go toolchain as a prerequisite for end users
- Makes installation a single command instead of clone → build → copy → configure

## Design Decisions

### Scope: Linux Only (Initially)

- Chosen: Linux only (amd64 + arm64)
  - Simpler initial scope, covers the primary target
  - macOS support can be added later by extending the release matrix and installer detection
- Considered: Linux + macOS from day one
  - Broader audience immediately
  - More CI matrix entries and testing surface area

### Release Artifact Format: Tarball

- Chosen: `.tar.gz` tarballs containing the binary + shell wrappers
  - Universally supported, `tar` and `gzip` are available on all Linux systems
  - Single archive per platform containing `to-backend`, `to.bash`, `to.zsh`, `to.fish`
- Considered: Raw binaries (no archive)
  - Simpler download but requires multiple fetches (binary + each wrapper)
  - No checksum file for the set
- Considered: `.tar.zst` (zstd compression)
  - Better compression but `zstd` not universally installed

### Release Workflow: GitHub Actions

- Chosen: GitHub Actions workflow triggered on semantic version tags (`*.*.*`, e.g. `1.2.3`)
  - Already using GitHub Actions for CI
  - Native integration with GitHub Releases
  - Uses `GOARCH` cross-compilation (no CGo, so cross-compiling is trivial)
- Considered: GoReleaser
  - Feature-rich but adds a dependency and configuration complexity
  - Overkill for a simple binary with no CGo

### Install Location: `~/.local/bin/`

- Chosen: `~/.local/bin/` (consistent with existing `just install`)
  - XDG-compliant, no sudo required
  - Matches the existing install target so curl-installed and locally-built installs are interchangeable
- Considered: Configurable via environment variable
  - Installed path is overridden via `TO_INSTALL_DIR` (default: `$HOME/.local/bin`)
  - Value is resolved to a directory path and may be relative to current working directory or absolute

### Shell Wrapper Detection

- Chosen: Detect user's current shell via `$SHELL` and print the appropriate `source` line
  - Simple, covers the common case
  - All three wrappers (`to.bash`, `to.zsh`, `to.fish`) are always installed; the instruction just tells the user which one to source
- Considered: Only install the wrapper for the detected shell
  - Saves trivial disk space but breaks if user switches shells

### Checksum Verification

- Chosen: SHA-256 checksum file uploaded alongside release artifacts
  - The installer downloads and verifies the checksum before installing
  - Standard practice for curl-pipe-bash installers
- Considered: No checksum verification
  - Simpler but insecure for a curl-pipe-bash installer

### Version Pinning

- Chosen: Support `TO_INSTALL_VERSION` to pin a specific release.
  - Supports reproducible installs for automation and air-gapped/controlled environments.
  - If unset, installer defaults to fetching the latest release.
- Considered: Hardcoded latest-only behavior
  - Simpler but prevents reproducible installs and rollback workflows.

## Task List

### GitHub Release Workflow

- [x] Create `.github/workflows/release.yml` triggered on `*.*.*` tags
- [x] Cross-compile `to-backend` for `linux/amd64` and `linux/arm64`
- [x] Package each build into a tarball: `to-<version>-linux-<arch>.tar.gz` containing `to-backend`, `to.bash`, `to.zsh`, `to.fish`
- [x] Generate SHA-256 checksums file (`checksums.txt`)
- [x] Upload tarballs and checksums to the GitHub Release

### Installer Script

- [x] Create `install.sh` at the repo root (POSIX shell, `#!/bin/sh`)
- [x] Detect OS (Linux only initially; error on unsupported)
- [x] Detect architecture (`x86_64` → `amd64`, `aarch64`/`arm64` → `arm64`)
- [x] Determine target version:
  - Use `TO_INSTALL_VERSION` if set
  - Otherwise fetch latest release version from GitHub API (`/repos/bambam955/to/releases/latest`)
- [x] Validate `TO_INSTALL_VERSION` if set:
  - Value must match an exact tag present in GitHub releases (e.g. `1.2.3`)
  - Fail with a clear error when unresolved/invalid
- [x] Download the correct tarball and checksums file to a temp directory
- [x] Verify SHA-256 checksum before extracting
- [x] Resolve and validate install directory:
  - `TO_INSTALL_DIR` default is `$HOME/.local/bin`
  - Accept both relative and absolute values, resolved to an absolute path before use
  - Fail with a clear error for non-directory/uncreatable paths
- [x] Extract tarball and install `to-backend` + shell wrappers to resolved install directory
- [x] Print post-install instructions (which `source` line to add based on `$SHELL`)

### Documentation and Testing

- [x] Update README.md with curl install instructions
- [x] Add `TO_INSTALL_DIR` environment variable override for install location
- [x] Test installer script with shellcheck
