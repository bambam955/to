---
status: completed
author: Bennett Moore
creation_date: 2026-01-25
approved_by: Bennett Moore
completed_date: 2026-02-03
---

# Bash-Go Integration Architecture

## Overview

This spec defines the integration architecture between the bash wrapper function `to()` and the Go backend binary `to-backend`. The focus is on clear communication protocols, error handling, and seamless user experience.

## What is Being Proposed

- Define process communication protocol between bash and Go components
- Specify error handling strategy using stderr and Unix exit codes
- Design installation and deployment approach for both components
- Outline shell completion integration patterns

## Communication Protocol

### Navigation Success Response

**Format**: `[to] <absolute_path>`

**Usage**: When navigation succeeds, `to-backend` outputs this exact format to stdout. The bash wrapper parses this prefix to distinguish navigation success from other output and performs the `cd` operation.

### Error Communication

**Strategy**: All errors go to stderr, no parsing required by bash

**Exit Codes**: Standard Unix conventions

- `0`: Success (only for non-navigation commands)
- `1`: Any error condition

**Error Messages**: Clear, human-readable messages prefixed with "error:" that the bash wrapper simply forwards to the user without modification.

## Installation Architecture

### Component Locations

**Binary**: `~/.local/bin/to-backend`
**Wrapper**: `~/.local/bin/to.sh` (sourced by shell as `to`)

### Installation Process

1. Build Go binary and install to `~/.local/bin/to-backend`
2. Install bash wrapper to `~/.local/bin/to.sh` and source in shell configuration
3. Verify user's PATH includes `~/.local/bin`

## Error Handling Strategy

### Go Backend Error Categories

**Validation Errors**: Invalid input, missing arguments, format issues
**Database Errors**: File operations, permission issues, corrupted data
**Operation Errors**: Alias not found, directory doesn't exist

### Bash Wrapper Error Propagation

**Exit Code Preservation**: Bash wrapper returns the same exit code as Go backend
**Error Forwarding**: stderr messages pass through unchanged without modification

## Process Lifecycle

### Startup Behavior

**First Run**: Automatically create database directory and file
**Database Location**: XDG-compliant `~/.config/to/database.json`

### Process Communication

**Synchronous**: All operations are synchronous and blocking
**Single-shot**: Each command spawns new Go process, no daemon
**Stateless**: No persistent process or shared memory

## Design Decisions

### Communication Protocol: Prefix Format vs JSON

**Decision**: `[to] <path>` prefix format for navigation

**Options Considered:**

- **JSON Response**: `{ "status": "success", "path": "/path" }`
  - Pros: Structured, extensible
  - Cons: Parsing overhead, verbose
- **Exit Codes Only**: Different codes for different outcomes
  - Pros: Standard Unix practice
  - Cons: Limited information, cannot convey paths
- **Prefix Format (chosen)**: `[to] <path>`
  - Pros: Simple parsing, clear separation, human-readable
  - Cons: Custom protocol, limited to navigation success

**Rationale**: Prefix format provides the simplest parsing while clearly distinguishing navigation success from other output.

### Error Handling: Parsing vs Forwarding

**Decision**: Forward stderr messages without parsing

**Options Considered:**

- **Parse and Reformat**: Parse errors and provide custom messages
  - Pros: Consistent user experience
  - Cons: Complex maintenance, loses Go backend context
- **Structured Error Format**: JSON errors with error codes
  - Pros: Programmatic handling
  - Cons: Overkill for CLI tool, verbose output
- **Direct Forwarding (chosen)**: Pass through stderr unchanged
  - Pros: Simple, preserves context, no maintenance burden
  - Cons: Less control over user-facing messages

**Rationale**: Direct forwarding maintains simplicity and preserves the Go backend's error handling context.

## Task List

### Communication Protocol

- [x] Define `[to] <path>` navigation response format
- [x] Implement Go backend navigation output formatting
- [x] Create bash wrapper parsing logic for navigation
- [x] Design error message format standards

### Installation System

- [x] Create Go build process for `to-backend`
- [x] Design bash wrapper installation to `~/.local/bin/`
- [x] Implement shell sourcing mechanism
- [x] Create PATH verification logic

### Error Handling

- [x] Implement Go backend error categorization
- [x] Design stderr message formatting standards
- [x] Create bash wrapper error forwarding logic
- [x] Implement exit code preservation

### Process Management

- [x] Design database initialization behavior
- [x] Implement XDG-compliant configuration paths
- [x] Create error handling for missing database
- [x] Design graceful startup behavior
