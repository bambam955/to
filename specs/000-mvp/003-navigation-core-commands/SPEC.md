---
number: 3
status: completed
author: Bennett Moore
creation_date: 2026-01-25
approved_by: Bennett Moore
approval_date: 2026-02-03
---

# Navigation and Core Commands

## Overview

This spec defines the primary user-facing commands for the "to" tool: default navigation, alias registration, and alias unregistration. These commands represent the core CRUD operations that users will interact with most frequently.

## What is Being Proposed

- Implement default navigation command: `to <alias>`
- Implement `reg` command for alias registration
- Implement `unreg` command for alias removal
- Define validation and error handling for all operations
- Design user interaction patterns and feedback
- Use Cobra package for CLI framework

## Command Specifications

### Default Navigation

**Usage**: `to <alias>`
**Purpose**: Navigate to the directory associated with an alias

**Behavior Flow**:

1. Look up alias in database
2. Validate that the directory still exists
3. Update last_visited timestamp in database
4. Output `[to] <absolute_path>` for bash wrapper to process
5. Return appropriate exit code (0 for success, 1 for error)

**Success Scenario**: User navigates to an existing alias with valid directory
**Error Scenarios**: Alias doesn't exist, directory no longer exists, database errors

### Registration Command

**Usage**: `to reg <alias> <directory>`
**Purpose**: Register a new alias pointing to a directory

**Validation Requirements**:

- Alias name: must follow pattern (starts alphanumeric, contains alphanumeric/hyphens/underscores)
- Directory: must exist and be accessible
- Uniqueness: alias name must not already exist in database
- Path handling: resolve relative paths to absolute paths

**Success Scenario**: New alias created and saved to database
**Error Scenarios**: Invalid alias format, alias already exists, directory doesn't exist, permission errors

**Duplicate Directory Handling**: If other aliases point to same directory, show warning but allow registration

### Unregistration Command

**Usage**: `to unreg <alias>`
**Purpose**: Remove an alias from the database

**Validation Requirements**:

- Alias must exist in database
- Must have proper permissions to modify database file

**Success Scenario**: Alias removed from database and file saved
**Error Scenarios**: Alias doesn't exist, database permission errors, file corruption

## User Experience Design

### Success Messages

**Navigation Success**: No output to user (handled via `[to] <path>` protocol)
**Registration Success**: Clear confirmation message
**Unregistration Success**: Clear confirmation message

### Error Messages

**Error Prefix**: All error messages use "error:" prefix
**Validation Errors**: Clear guidance on how to fix the issue
**Database Errors**: Distinguish between permission issues, corruption, and other problems
**Operation Errors**: Specify what went wrong and suggest fixes

### Warnings

**Duplicate Directory**: Inform user of existing aliases pointing to same directory

## Go Package Strategy

### CLI Framework

**Package Choice**: Cobra for CLI command structure and argument parsing

**Rationale**:

- Industry standard for Go CLI applications
- Built-in support for subcommands and validation
- Automatic help generation via --help flag
- Good integration with other Go ecosystem tools

### Command Structure

**Root Command**: `to`

- Default subcommand handles navigation
- Subcommands: `reg`, `unreg`, `list`, `clean`, `exp`
- Flags and arguments handled automatically by Cobra

**Error Handling**: Cobra's built-in error handling combined with custom "error:" prefixed messages
**Help System**: Cobra's automatic help generation via --help flag

## Design Decisions

### Navigation: Database Update vs Read-Only

**Decision**: Update last_visited timestamp on navigation

**Options Considered**:

- **Read-Only Navigation**: Only lookup, no updates
  - Pros: Faster, simpler, less database writes
  - Cons: No usage tracking, missed analytics opportunity
- **Update on Navigation (chosen)**: Update timestamp
  - Pros: Usage data for future features, tracks activity
  - Cons: Extra database write, slightly slower

**Rationale**: Timestamp updates provide valuable usage data with minimal performance impact and enable future features.

### Registration: Strict vs Lenient Validation

**Decision**: Strict validation with clear error messages

**Options Considered**:

- **Lenient**: Accept more formats, fix automatically
  - Pros: More forgiving, easier for users
  - Cons: Ambiguous behavior, potential security issues
- **Strict (chosen)**: Clear validation rules
  - Pros: Predictable behavior, security-conscious
  - Cons: More restrictive, learning curve

**Rationale**: Strict validation prevents ambiguous behavior and security issues while providing clear feedback for users.

### Duplicate Directory Handling: Error vs Warning

**Decision**: Allow duplicates with warning

**Options Considered**:

- **Error on Duplicates**: Prevent multiple aliases to same directory
  - Pros: Prevents confusion, simpler database
  - Cons: Reduces flexibility, different logical groupings needed
- **Warning on Duplicates (chosen)**: Allow but inform user
  - Pros: Maximum flexibility, allows logical organization
  - Cons: More complex, potential confusion

**Rationale**: Users often want different logical aliases for the same directory (e.g., "work" and "projectx"), so warnings provide information without restricting functionality.

## Task List

### Default Navigation

- [x] Define navigation behavior and success/error scenarios
- [x] Design timestamp update mechanism
- [x] Specify directory validation requirements
- [x] Design error messages for missing/invalid aliases
- [x] Create Cobra command structure for navigation

### Registration Command

- [x] Define validation requirements for alias names and directories
- [x] Design path resolution strategy (relative to absolute)
- [x] Specify duplicate detection and warning behavior
- [x] Design success and error message patterns
- [x] Create Cobra subcommand for registration

### Unregistration Command

- [x] Define existence validation requirements
- [x] Design removal behavior and database update process
- [x] Specify error handling for permission and file issues
- [x] Design success confirmation messages
- [x] Create Cobra subcommand for unregistration

### User Experience

- [x] Design consistent message formatting across commands
- [x] Define error categorization and user guidance
- [x] Specify warning display behavior
- [x] Design feedback for all user actions

### Integration Points

- [x] Define how commands interact with database layer
- [x] Specify error propagation from database operations
- [x] Design validation flow between command and database layers
- [x] Define state change requirements for each operation

### Command Implementation

- [x] Implement default navigation command (uses database interface from 002)
- [x] Implement registration command with validation
- [x] Implement unregistration command
- [x] Set up Cobra root command structure with subcommands
- [x] Configure error message formatting with "error:" prefix

### CLI Framework

- [x] Design argument validation using Cobra
- [x] Implement help system and usage messages
- [x] Configure error handling integration
