---
number: 7
status: draft
author: Bennett Moore
creation_date: 2026-03-22
---

# Changelog Generation with git-cliff

## Overview

Adopt `git-cliff` to generate and maintain `CHANGELOG.md` from Conventional Commit history.

This spec defines a repeatable workflow for:

- local changelog generation before release
- CI validation to prevent manual drift
- predictable changelog sections based on commit types

## Why It's Needed

- The project currently has no standardized, automated changelog generation flow.
- Manual changelog editing is error-prone and can miss important changes.
- A `git-cliff` based workflow creates consistent release notes from commit history.

## Design Decisions

### Tool Choice: git-cliff

- Chosen: use `git-cliff` as the single source of truth for changelog generation.
- Considered: handwritten `CHANGELOG.md`
  - simple for small projects
  - does not scale and is easy to make inconsistent
- Considered: custom scripts
  - flexible but adds maintenance burden versus a mature tool

### Commit Parsing Strategy

- Chosen: parse Conventional Commit prefixes (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`) into changelog groups.
- Chosen: include breaking changes in a dedicated section when marked in commit messages.
- Considered: flat chronological list
  - easier to implement
  - harder for users to scan by change type

### Generation Workflow

- Chosen: generate changelog via explicit project commands (`just`) and checked-in config.
- Chosen: CI checks that generated output is up to date.
- Considered: generate only in release CI
  - less local setup
  - hides drift until release time

### Versioning and Tag Behavior

- Chosen: generate per-release sections based on Git tags, with an unreleased section for changes since last tag.
- Considered: full regenerate without release boundaries
  - simpler mental model
  - less useful for release-oriented reading

## Task List

### Tooling and Configuration

- [ ] Add `git-cliff` configuration file at repository root (for example `cliff.toml`) with:
  - commit parser/grouping rules
  - release header/body/footer templates
  - tag and unreleased configuration
- [ ] Add or update project commands in `justfile` for local changelog generation (for example `just changelog`).
- [ ] Document prerequisite installation guidance for `git-cliff` in developer docs.

### Changelog File Management

- [ ] Add `CHANGELOG.md` if it does not exist, with generated content seed.
- [ ] Ensure generated content is stable and deterministic for repeat runs.
- [ ] Preserve any required project-specific intro/header text above generated sections.

### CI Enforcement

- [ ] Add CI step to verify changelog generation is up to date.
- [ ] Fail CI when `CHANGELOG.md` differs from generated output.
- [ ] Document the remediation command in CI failure output.

### Documentation

- [ ] Update `README.md` with changelog workflow:
  - how to generate/update changelog locally
  - when to run it in the release process
- [ ] Add a short contribution note reinforcing Conventional Commit usage for clean changelog output.
