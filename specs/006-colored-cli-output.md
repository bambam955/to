---
number: 6
status: draft
author: Bennett Moore
creation_date: 2026-03-19
approved_by: Bennett Moore
approval_date: 2026-03-19
---

# Colored CLI Output

## Overview

Add optional, consistent terminal color to `to` command output so the CLI feels polished in interactive use.

The new behavior should cover:

- Cobra help output (root and subcommands)
- Error messages
- List output (and any user-facing table/list style output)

This must not break machine-facing behavior (especially protocol output consumed by the shell wrapper for `to` path switching).

## Why It's Needed

- Current output is entirely plain text, which is harder to scan in common terminal workflows.
- Error paths and list output are difficult to distinguish at a glance.
- Users expect modern CLI affordances for readability, especially for command completion and cleanup tasks.

## Design Decisions

### In-scope output

- In scope: command help text, command errors, and list-style output from interactive commands.
- Out of scope: protocol output consumed by the shell wrapper (`[to] <path>` responses and raw path export behavior), since that must remain machine-parseable and unstyled.
- Out of scope for v1: spinner/progress UI and editor/interactive prompts.

### Color mode policy

- Chosen: Use auto-detection only. Emit ANSI colors when the active output stream is a TTY; otherwise emit plain output.
- Considered: Explicit user color controls
  - increases complexity for a small ergonomics gain in v1
- Considered: Always-on color
  - simpler implementation but produces noisy output in logs and pipes.

### Theme

- Chosen: Use a small, fixed palette:
  - command headers/help labels
  - list names and values
  - success messages
  - warning messages
  - error messages
- Considered: Third-party theme system
  - unnecessary complexity for v1

### Tooling dependency

- Chosen: Use `github.com/fatih/color` for ANSI rendering.
- Considered: New internal styling package
  - unnecessary for v1 now that a small helper exists in the command layer.

## Task List

### Core format layer

- [ ] Add a small command-layer style helper (no new `pkg/ui` package) that wraps `github.com/fatih/color`:
  - shared theme constants (error/success/warn/help/list styles)
  - helper methods for `Sprintf`-style styling that become no-ops when color is disabled
  - color enablement detection via TTY checks
- [ ] Add unit tests for theme resolution and TTY-based enablement in `cmd/to-backend` style helper tests.

### Root/help output

- [ ] Update `cmd/to-backend/root.go` to pass the shared formatter into Cobra.
- [ ] Configure Cobra usage/help template to apply style to:
  - usage lines
  - short/long descriptions
  - example blocks
  - flag sections
- [ ] Add tests in `cmd/to-backend/*_test.go` to verify generated help is styled when output is a TTY and unstyled otherwise.

### Errors

- [ ] Route command-level errors through the formatter helper so they get consistent color treatment.
- [ ] Keep protocol-style error strings compatible; error outputs intended for the wrapper must remain parseable.
- [ ] Add tests in `cmd/to-backend/*_test.go` for:
  - colored error output when output is a TTY
  - uncolored errors when output is non-TTY.

### List output

- [ ] Update `cmd/to-backend/list.go` to colorize:
  - headers or list context
  - alias names
  - directory path output
- [ ] Ensure list formatting remains stable so existing tests can assert token positions where expected.
- [ ] Add or update list tests to cover TTY and non-TTY snapshots.

### Protocol safety

- [ ] Audit `pkg/protocol` and commands that emit path-only output (`navigate`, `exp`, etc.) and ensure those outputs stay unstyled.
- [ ] Add regression tests for protocol path outputs to guarantee no ANSI escapes appear in wrapper-consumed output.

### Docs and defaults

- [ ] Update README with color behavior and examples:
  - colors appear on interactive terminals
  - no colors for non-TTY output.
