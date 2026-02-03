---
status: draft
author: Bennett Moore
creation_date: 2026-01-25
approved_by: Bennett Moore
approval_date: 2026-02-03
---

# Utility Commands

## Overview

This spec defines the utility and management commands for the "to" tool: listing aliases, cleaning invalid directories, and expanding alias paths. These commands provide database management and information retrieval capabilities that complement the core navigation functionality.

## What is Being Proposed

- Implement `list` command for displaying all aliases
- Implement `clean` command for removing invalid directory entries
- Implement `exp` command for resolving alias paths
- Define output formatting standards for display commands
- Design bulk operation handling and user feedback

## Command Specifications

### List Command

**Usage**: `to list`
**Purpose**: Display all registered aliases in a readable format

**Output Requirements**:

- Show alias name and corresponding directory path
- Use column formatting for readability
- Display aliases in a consistent order (alphabetical by name)
- Handle long paths gracefully without breaking formatting

**Display Options**:

- Show alias names and directory paths in aligned columns
- Use color coding or formatting to distinguish aliases from paths
- Handle edge cases like very long paths or alias names

**Success Scenario**: Clear, readable list of all aliases
**Error Scenarios**: Database access errors, corruption, permission issues

### Clean Command

**Usage**: `to clean`
**Purpose**: Remove aliases pointing to directories that no longer exist

**Behavior Flow**:

1. Iterate through all aliases in database
2. Check if each directory path still exists and is accessible
3. Remove invalid entries from database
4. Show user what was cleaned up
5. Save the updated database

**Output Requirements**:

- Show each invalid entry being removed
- Provide summary of cleanup actions taken
- Handle case where no cleanup is needed

**Success Scenario**: Invalid aliases removed with confirmation
**Error Scenarios**: Database access errors, permission issues during validation

### Expand Command

**Usage**: `to exp <alias>`
**Purpose**: Show the full directory path for a given alias

**Output Requirements**:

- Display just the absolute path (no prefix or formatting)
- Return exit code 0 for valid aliases, 1 for invalid
- No additional formatting or explanatory text

**Use Cases**:

- Script integration where just the path is needed
- Checking where an alias points before navigation
- Automated tools that need to resolve alias paths

**Success Scenario**: Output to absolute path to stdout
**Error Scenarios**: Alias doesn't exist, database access errors

## Output Formatting Standards

### Display Principles

**Readability**: Prioritize human-readable output over machine parsing
**Consistency**: Use similar formatting across display commands
**Clarity**: Distinguish between different types of information

### Error Messages

**Error Prefix**: All error messages use "error:" prefix
**Validation Errors**: Clear guidance on how to fix issues
**Database Errors**: Clear database-related issues
**System Errors**: File system and permission issues

### Clean Command Reporting

**Action Reporting**: Show each invalid alias being removed
**Summary**: Provide count of items removed and items remaining
**No-Op Case**: Clear message when no cleanup needed

## User Experience Design

### Command Discovery

**Helpful Output**: Commands should be self-documenting
**Clear Usage**: Show usage patterns when invoked incorrectly
**Consistent Behavior**: Similar error handling across all utility commands

### Progress Feedback

**Long Operations**: Provide feedback for potentially slow operations (like clean with many entries)
**Success Confirmation**: Clear confirmation when operations complete
**Error Context**: Provide enough context for users to understand and fix issues

### Error Handling

**Graceful Degradation**: Handle partial failures appropriately
**Clear Messages**: Explain what went wrong and why
**Recovery Guidance**: Suggest corrective actions where possible

## Design Decisions

### List Formatting: Simple vs Complex

**Decision**: Simple column-based formatting

**Options Considered:**

- **Complex Formatting**: Colors, icons, tree structures
  - Pros: Visually appealing, information-rich
  - Cons: Complex maintenance, dependency on terminal capabilities
- **Simple Columns (chosen)**: Basic aligned columns
  - Pros: Simple, reliable, works everywhere
  - Cons: Less visually impressive

**Rationale**: Simple formatting ensures reliability across different terminal environments and focuses on function over appearance.

### Clean Strategy: Interactive vs Automatic

**Decision**: Automatic cleanup with reporting

**Options Considered:**

- **Interactive**: Prompt for each invalid entry
  - Pros: User control, safe
  - Cons: Tedious for many invalid entries
- **Automatic (chosen)**: Clean all invalid entries automatically
  - Pros: Fast, efficient
  - Cons: Less user control

**Rationale**: Invalid directory entries are clearly errors that should be removed automatically, while user can always re-register if needed.

### Expand Output: Path Only vs Formatted

**Decision**: Output just the absolute path

**Options Considered:**

- **Formatted Output**: Include labels, prefixes, or explanations
  - Pros: Clear, self-documenting
  - Cons: Requires parsing for scripts, verbose
- **Path Only (chosen)**: Raw absolute path to stdout
  - Pros: Perfect for script integration, simple
  - Cons: Less context for humans

**Rationale**: The expand command is primarily for script integration and automation, where raw path output is most useful.

## Task List

### List Command

- [ ] Define column-based output formatting
- [ ] Design handling for long paths and alias names
- [ ] Specify sorting and ordering behavior
- [ ] Design error handling and display

### Clean Command

- [ ] Define directory validation logic for each alias
- [ ] Design batch removal behavior
- [ ] Specify progress reporting and user feedback
- [ ] Design summary output for cleanup results

### Expand Command

- [ ] Define path-only output format
- [ ] Design behavior for non-existent aliases
- [ ] Specify exit code conventions
- [ ] Design error handling and messaging

### Output Standards

- [ ] Define consistent formatting principles across commands
- [ ] Design error message standards using "error:" prefix
- [ ] Specify clean formatting guidelines

### User Experience

- [ ] Define clear success and error feedback
- [ ] Specify error recovery guidance
- [ ] Design consistent interaction patterns
