---
number: 7
status: in-progress
author: Bennett Moore
creation_date: 2026-03-22
approved_by: Bennett Moore
approval_date: 2026-03-23
---

# Changelog Generation with git-cliff

## Overview

Adopt `git-cliff` to generate and maintain `CHANGELOG.md` from Conventional Commit history and the project's semantic release tags defined in spec 008.

The initial changelog backfill should anchor on the existing `0.1.0` tag, so running `git-cliff` automatically generates history from that release forward.

This spec defines a repeatable workflow for:

- release-PR changelog preparation via `just prep-release <version>`
- CI validation to prevent manual drift on release branches
- predictable changelog sections based on commit types
- GitHub Release notes sourced from the matching `CHANGELOG.md` section

## Why It's Needed

- The project currently has no standardized, automated changelog generation flow.
- Manual changelog editing is error-prone and can miss important changes.
- A `git-cliff` based workflow creates consistent release notes from commit history.
- Release preparation currently requires too much manual branching, changelog editing, and PR setup.

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
- Chosen: skip `spec:` commits and merge-noise commits from user-facing release notes.
- Considered: flat chronological list
  - easier to implement
  - harder for users to scan by change type

### Generation Workflow

- Chosen: generate changelog via explicit project commands (`just`) and checked-in config.
- Chosen: `just prep-release <version>` calls `ci/prepare-release.sh <version>` to automate release-branch creation and PR setup.
- Chosen: release preparation starts from a clean `main` checkout and fails hard on unsafe state.
- Chosen: CI checks that generated output is up to date.
- Considered: generate only in release CI
  - less local setup
  - hides drift until release time

### Changelog as Release Notes Source of Truth

- Chosen: `CHANGELOG.md` is the single source of truth for both the checked-in changelog and GitHub Release notes.
- Chosen: the tag-triggered publishing workflow defined in spec 005 must use the matching generated `CHANGELOG.md` section as the GitHub Release body.
- Considered: generating GitHub Release notes separately from the checked-in changelog
  - creates two release-note sources that can drift
  - increases release-time ambiguity

### CI Scope

- Chosen: changelog freshness checks run only on `release/*` branch pushes and pull requests whose head branch matches `release/*`.
- Considered: enforcing changelog freshness on every branch
  - catches drift earlier
  - adds noise to ordinary feature development when release notes are not being prepared

### Release Tag Integration

- Chosen: generate per-release sections based on the repository's release tags, with an unreleased section for changes since last tag.
- Chosen: keep tag names and changelog headings aligned so the release version in `CHANGELOG.md` matches the release tag exactly.
- Chosen: treat the semantic versioning rules in spec 008 as the source of truth for release tag names.
- Chosen: after the release PR is merged, the release is cut by manually pushing the matching bare semantic version tag.
- Chosen: the publishing workflow may require GitHub deployment approval before the release is published.
- Considered: full regenerate without release boundaries
  - simpler mental model
  - less useful for release-oriented reading

## Task List

### Tooling and Configuration

- [ ] Add `git-cliff` configuration file at repository root (for example `cliff.toml`) with:
  - commit parser/grouping rules
  - skip rules for `spec:` and merge-noise commits
  - release header/body/footer templates
  - tag and unreleased configuration
- [ ] Add or update project commands in `justfile` for local changelog generation, including `just prep-release <version>`.
- [ ] Add `ci/prepare-release.sh` to validate release-prep prerequisites, create `release/<version>` from clean `main`, regenerate `CHANGELOG.md`, commit the result, push the branch, and open a pull request with GitHub CLI.
- [ ] Document prerequisite installation guidance for `git-cliff` and `gh` in `README.md`.

### Changelog File Management

- [ ] Add `CHANGELOG.md` if it does not exist, with generated content seed.
- [ ] Seed the initial changelog generation from the existing `0.1.0` tag so historical sections are backfilled automatically.
- [ ] Ensure generated content is stable and deterministic for repeat runs.
- [ ] Preserve any required project-specific intro/header text above generated sections.
- [ ] Ensure generated release sections use exact release tag identifiers from Git tags.
- [ ] Ensure `just prep-release <version>` generates the requested release section before the tag exists.

### CI Enforcement

- [ ] Add CI step to verify changelog generation is up to date only for `release/*` branch pushes and pull requests from `release/*` branches.
- [ ] Fail CI when `CHANGELOG.md` differs from generated output.
- [ ] Document the remediation command in CI failure output using `just prep-release <version>`.

### Documentation and Release Process

- [ ] Update `README.md` with changelog workflow:
  - how to run `just prep-release <version>` from clean `main`
  - how the command creates `release/<version>` and opens the release PR
  - that the release is cut by merging the release PR and then pushing the matching tag
- [ ] Document that spec 005 publishes the GitHub Release from the matching `CHANGELOG.md` section, with deployment approval if configured.
- [ ] Add a short contribution note reinforcing Conventional Commit usage for clean changelog output and noting that `spec:` commits are excluded from user-facing release notes.
