---
number: 8
status: draft
author: Bennett Moore
creation_date: 2026-03-22
---

# Semantic Versioning and Version Flag

## Overview

Adopt Semantic Versioning 2.0.0 for `to` releases and add a version flag so users can inspect the installed binary version.

This spec defines a repeatable workflow for:

- release tag naming and version numbering
- build-time version metadata embedded into binaries
- `to -v` / `to --version` output for installed binaries
- release-process alignment between tags, binaries, and release artifacts

## Why It's Needed

- Users need a simple way to confirm which version of `to` is installed.
- Release automation needs a single version source of truth.
- Changelog generation and published artifacts should align with the same semantic version.

## Design Decisions

### Semantic Version Format

- Chosen: use Semantic Versioning 2.0.0 with `MAJOR.MINOR.PATCH` releases.
- Chosen: publish bare release tags like `1.2.3` without a `v` prefix.
- Considered: `v1.2.3` tags
  - common in Go projects
  - adds an extra normalization step for release artifacts and changelog headings

### Version Source of Truth

- Chosen: embed version metadata at build time from the tagged release ref.
- Chosen: include supplementary build metadata, such as commit SHA, only when needed for development builds.
- Considered: runtime Git inspection
  - not available in release tarballs or installed binaries
  - makes version reporting unreliable after distribution

### Version Flag Behavior

- Chosen: support `to -v` and `to --version` at the CLI root.
- Chosen: print a short, machine-readable version string and exit 0.
- Chosen: untagged development builds fall back to a clear placeholder such as `dev`.
- Considered: a verbose multi-line banner
  - harder to parse
  - unnecessary for a simple CLI version check

### Release Integration

- Chosen: the version flag output must match the Git tag used to build the release artifact exactly.
- Chosen: tagged release builds should stamp the binary and release artifacts with the same semantic version.
- Considered: deriving version from package state at runtime
  - mismatched with distributed artifacts
  - complicates reproducible release builds

## Task List

### CLI Version Support

- [ ] Add build-time version metadata to the Go backend.
- [ ] Add `-v` / `--version` support at the CLI root to print the installed version.
- [ ] Define development-build fallback behavior for binaries built outside a tagged release.
- [ ] Add unit tests for version output in release and development build modes.

### Build and Release Wiring

- [ ] Update build and release commands so tagged builds inject the release version into the binary.
- [ ] Ensure release artifacts report the same semantic version as the tagged binary.
- [ ] Document the versioning contract for tagged releases in developer docs.

### Documentation

- [ ] Update README with version flag usage and release tag format.
