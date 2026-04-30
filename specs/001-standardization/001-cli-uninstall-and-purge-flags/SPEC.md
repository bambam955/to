---
number: 1
status: in-progress
author: Bennett Moore
creation_date: 2026-04-07
approved_by: Bennett Moore
approval_date: 2026-04-07
---

# CLI Uninstall and Purge Flags

This spec adds root-level CLI flags for removing an installed `to` setup without relying on `just` recipes. The new flags make uninstall and purge available anywhere the backend binary is installed, while preserving the existing shell-wrapper navigation protocol and error formatting.

The design builds on the installation layout defined in [Bash-Go Integration Architecture](specs/000-mvp/001-bash-golang-integration/SPEC.md), but formalizes destructive install-management operations as top-level behavior outside the MVP spec tree.

## Goals

- Add `-U` / `--uninstall` to remove the installed backend and known shell wrappers.
- Add `-P` / `--purge` to perform uninstall and remove default configuration data.
- Keep uninstall and purge as explicit operation modes alongside register, unregister, list, clean, and expand.
- Prevent ambiguous flag combinations by rejecting multiple operation-mode flags in one invocation.
- Keep missing install targets and missing config directories as successful no-op cases.

## Design Decisions

### Root flag surface

- Chosen: add both short and long flags: `-U` / `--uninstall` and `-P` / `--purge`
  - Matches the existing root command pattern where management operations expose both forms.
  - Keeps the short destructive forms visually distinct from the existing lowercase aliases.
- Considered: short-only flags
  - Smaller surface area.
  - Less discoverable in help and documentation.
- Considered: long-only flags
  - Clearer for infrequent operations.
  - Inconsistent with the rest of the CLI.

### Uninstall scope

- Chosen: remove `to-backend` plus all known wrappers: `to.bash`, `to.zsh`, and `to.fish`
  - Cleans up installed assets regardless of which shell is currently active.
  - Avoids leaving stale wrappers behind after switching shells.
- Considered: remove only the current shell wrapper
  - Matches the `just install shell=...` flow more closely.
  - Leaves behind other installed wrappers that still belong to `to`.
- Considered: remove only a legacy `to.sh`
  - Preserves an older naming assumption.
  - Does not match the current repository layout or shipped wrapper filenames.

### Purge scope

- Chosen: purge mirrors `just purge` by removing the full default config directory `~/.config/to/`
  - Removes the database and any future tool-managed config in one step.
  - Keeps the destructive scope simple and consistent with the documented recipe.
- Considered: remove only the database file
  - Narrower blast radius.
  - Leaves behind now-empty tool state that users asked to purge.
- Considered: resolve the purge target through `TO_DB`
  - More closely follows runtime database lookup.
  - Makes purge behavior depend on the current shell environment instead of the standard installation footprint.

### Flag conflict handling

- Chosen: reject any invocation that sets more than one operation-mode flag
  - Prevents silent precedence bugs in the current first-match dispatch model.
  - Makes destructive operations predictable and explicit.
- Considered: keep first-match dispatch
  - No extra validation logic.
  - Unsafe once uninstall and purge share the same root command surface.
