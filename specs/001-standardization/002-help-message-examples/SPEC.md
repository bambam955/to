---
number: 2
status: draft
author: Bennett Moore
creation_date: 2026-04-13
---

# Help Message Examples

The root help output should include a short examples section so new and returning users can see the common workflows directly from `to --help`.

The current help text lists flags, but it does not show how those flags combine with aliases and paths. Adding examples makes the command easier to learn without requiring users to leave the terminal or inspect documentation.

The examples should follow the broad shape of `mise help`: a distinct `Examples:` section after the flags, with command lines on the left and brief descriptions on the right. The section should cover the basic navigation and alias-management flows without turning help output into full documentation.

## Goals

- Show practical examples when running `to --help`.
- Include the most common root command flows: navigate, register, list, expand, unregister, and clean.
- Keep examples concise enough that the help output remains scannable.
- Preserve the existing root command usage, flags, version output, and error behavior.
- Keep example formatting stable for tests and readable in plain terminals.

## Design Decisions

### Examples section scope

- Chosen: add examples to the root command help output only
  - Directly addresses the user-facing `to --help` gap.
  - Keeps the change focused on the single command surface that currently exposes the CLI contract.
- Considered: add separate long-form documentation instead
  - Useful for deeper onboarding.
  - Does not help users who are already in the terminal and asking the CLI for help.

### Example set

- Chosen: include examples for navigation, registration, listing, path expansion, unregistering, and stale-alias cleanup
  - Covers the main daily workflow from creating an alias through using and maintaining it.
  - Keeps destructive or advanced install-management operations out of a basic examples section.
- Considered: show every available flag
  - Maximizes completeness.
  - Makes the section longer than needed for quick help.
- Considered: show only navigation and registration
  - Very compact.
  - Leaves users without examples for the other common maintenance operations.

### Formatting style

- Chosen: use a command-plus-description layout modeled after `mise help`
  - Makes examples easy to scan vertically.
  - Allows the examples to explain intent without adding paragraphs.
- Considered: list commands without descriptions
  - Easier to copy.
  - Less useful for users learning what each command is for.
- Considered: add longer prose under each example
  - Can explain edge cases.
  - Turns help output into documentation instead of a quick reference.

### Prompt prefix

- Chosen: prefix examples with `$`
  - Matches the familiar shell-command style used by `mise help`.
  - Makes clear that each example is a command to run.
- Considered: omit the prompt for easier copy and paste
  - Avoids copying an extra character.
  - Is less consistent with the requested comparison point.
