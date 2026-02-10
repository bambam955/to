---
status: completed
author: Bennett Moore
creation_date: 2026-01-25
approved_by: Bennett Moore
approval_date: 2026-02-03
---

# Database Design and Schema

## Overview

This spec defines the JSON database schema for storing directory aliases, along with Go package strategy for database operations. The design focuses on extensibility, maintainability, and XDG compliance.

## What is Being Proposed

- Design a JSON schema for storing aliases with metadata
- Define versioning strategy for future schema evolution
- Specify Go package usage (`encoding/json`)
- Design database file management and initialization
- Define XDG-compliant configuration paths

## JSON Schema Design

### Database Structure

**Root Object**: Contains version information and alias collection

- `version`: Database schema version
- `created_at`: Database creation timestamp
- `updated_at`: Last modification timestamp
- `aliases`: Array of alias objects

**Alias Object**:

- `name`: Unique identifier for the alias
- `directory`: Absolute path to the directory
- `created_at`: Alias creation timestamp
- `last_visited`: Timestamp of last navigation (optional)

**Version Migration Strategy**:

- Major version bump for breaking changes requiring migration
- Minor version bump for backward-compatible additions

### Field Validation Rules

**Alias Name**: Must start with alphanumeric, can contain alphanumeric, hyphens, underscores, maximum 50 characters
**Directory Path**: Must be absolute path, maximum 4096 characters
**Timestamps**: ISO 8601 format for all time fields

## Go Package Strategy

### Package Selection

**Primary Choice**: Standard `encoding/json` package

**Rationale**:

- Built-in, no external dependencies
- Sufficient performance for expected database size
- Well-maintained and stable
- Good integration with Go ecosystem and Cobra framework

### Data Structure Design

**Core Types**: Define Go structs that map to JSON schema

- Database struct with version, timestamps, and aliases array
- Alias struct with name, directory, creation, and optional last visited

**Database Interface**: Define CRUD operations for alias management

- Load/save database operations
- Add/remove/update alias operations
- Search and filter operations

## File Management

### Configuration Paths

**XDG Compliance**: Follow XDG Base Directory Specification

- Default path: `~/.config/to/database.json`
- Override via environment variable `TO_DB`
- Create parent directories if they don't exist

### File Operations Strategy

**Atomic Writes**: Use temporary file pattern to prevent corruption

- Write to temporary file first
- Atomically rename to final location
- Handle cleanup on failure

**Initialization**: Create empty database with version metadata on first run

- Ensure config directory exists with proper permissions
- Initialize with current timestamp and empty aliases array

## Validation Logic

### Alias Name Validation

**Pattern**: Must start with alphanumeric, can contain alphanumeric, hyphens, underscores

- Enforce length limits and character restrictions
- Ensure uniqueness within database
- Provide clear error messages for invalid formats

### Directory Path Validation

**Requirements**:

- Must be absolute path
- Directory must exist and be accessible
- User must have permissions to navigate to directory
- Normalize paths for consistent storage

## Error Handling

### Database Errors

**File System Errors**: Permission denied, disk full, corrupted JSON

- Provide clear error messages distinguishing permission issues from corruption
- Suggest corrective actions where possible

**Validation Errors**: Invalid alias names, non-existent directories

- Provide specific guidance on fixing validation failures
- Include examples of valid formats

### Error Categories

**Not Found**: Alias or directory doesn't exist
**Exists**: Alias name already in use
**Invalid**: Format or validation failures
**Corrupted**: Database file is malformed or unreadable

## Design Decisions

### JSON Schema: Simple vs Rich Metadata

**Decision**: Include basic metadata from the start

**Options Considered:**

- **Minimal Schema**: Just name and directory
  - Pros: Simple, fast, small files
  - Cons: Hard to add features later, no usage tracking
- **Rich Metadata (chosen)**: Include timestamps and future extensibility
  - Pros: Enables future features, versionable, extensible
  - Cons: Slightly larger files, more complex validation

**Rationale**: Including basic metadata from the start provides foundation for future features without significant overhead.

### Go Package: Standard vs Third-Party

**Decision**: Use standard `encoding/json` package

**Options Considered:**

- **Standard Library**: `encoding/json`
  - Pros: Built-in, no dependencies, stable
  - Cons: Slower than some alternatives
- **Third-Party**: `json-iterator/go` or similar
  - Pros: Faster performance, additional features
  - Cons: External dependency, maintenance overhead

**Rationale**: For expected database size (hundreds of aliases at most), standard library performance is more than sufficient and avoids dependency management complexity.

### Atomic Operations: In-Memory vs File-Based

**Decision**: Load entire database into memory, use atomic file writes

**Options Considered:**

- **File-Based Operations**: Direct file manipulation for each operation
  - Pros: Low memory usage, simple for individual operations
  - Cons: Complex race conditions, difficult to maintain consistency
- **In-Memory with Atomic Saves (chosen)**: Load into RAM, atomic writes
  - Pros: Simple logic, consistent state, easier testing
  - Cons: Higher memory usage, potential performance for very large datasets

**Rationale**: For typical usage patterns (small number of aliases), in-memory operations provide simpler, more reliable code with negligible performance impact.

## Task List

### Schema Design

- [x] Define JSON schema for database structure
- [x] Design versioning strategy for schema evolution
- [x] Create validation rules for all fields
- [x] Design metadata fields for future extensibility

### Database Layer Implementation

- [x] Implement Go data structures matching JSON schema
- [x] Create database interface with CRUD operations
- [x] Implement database loading and atomic file save operations
- [x] Add error categorization and validation error types

### File Management

- [x] Implement XDG-compliant path resolution
- [x] Create database initialization logic
- [x] Add directory creation with proper permissions
- [x] Implement environment variable override support

### Validation System

- [x] Implement alias name validation with regex
- [x] Add directory path validation and existence checking
- [x] Create duplicate detection logic
- [x] Design error categorization and messaging

### Error Handling

- [x] Implement structured error types
- [x] Add file system error handling
- [x] Create database corruption detection
- [x] Design user-friendly error messages
