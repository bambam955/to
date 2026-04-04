# Release Workflow

The release workflow for `to` is designed to keep manual steps small and keep
`CHANGELOG.md` as the source of truth for both the checked-in changelog and the
GitHub Release notes.

## Requirements

- `git-cliff` on your `PATH`
- GitHub CLI (`gh`) on your `PATH`
- an authenticated GitHub CLI session (`gh auth login`)

## Prepare a Release PR

Run `just prep-release <version>` from a clean `main` checkout.

The command:

1. fetches `origin/main` and verifies local `main` is in sync
2. creates `release/<version>`
3. chooses the release date for that branch
4. regenerates `CHANGELOG.md` for that exact release version
5. commits the changelog update with a `Release-Date: YYYY-MM-DD` trailer
6. pushes the release branch
7. opens the release PR against `main`

Example:

```bash
just prep-release 1.2.3
```

## Generate the Changelog Locally

`just gen-changelog` behaves differently depending on the current branch:

- on `main`, it uses the current `main` checkout
- on `release/*`, it uses the current release branch, infers the matching
  release version from the branch name, and reuses the most recent
  `Release-Date: YYYY-MM-DD` trailer reachable on that branch
- on any other branch, it fetches `origin/main` and renders from mainline history

That last rule keeps local implementation-branch commits from leaking into the
checked-in changelog.

The `Release-Date` trailer makes pending release branches deterministic. The
checked-in `CHANGELOG.md` no longer depends on the day `git-cliff` was run; it
depends on the release date committed to the branch.

If changelog validation fails on an existing `release/*` branch, run:

```bash
just gen-changelog
```

## Correct a Pending Release Date

If the intended release day changes after the release branch already exists, add
a follow-up commit on `release/<version>` with a new `Release-Date:` trailer and
regenerate the changelog:

```bash
git commit --allow-empty -m "chore(changelog): update release date for 1.2.3" -m "Release-Date: 2026-04-05"
just gen-changelog
git add CHANGELOG.md
git commit -m "chore(changelog): regenerate changelog for 1.2.3"
```

The newest valid `Release-Date` trailer on the release branch wins.

Use `just prep-release <version>` only to create a new release branch from a
clean `main` checkout.

## Publish the Release

After the release PR is merged, cut the release by pushing the matching bare
tag on the same calendar day as the committed `Release-Date`:

```bash
git tag 1.2.3
git push origin 1.2.3
```

The tag-triggered publishing workflow defined in spec `005` is expected to use
the matching `CHANGELOG.md` section as the GitHub Release notes source of
truth, with deployment approval if configured.

## Commit Expectations

Only user-facing Conventional Commit types are included in generated release
notes. `spec:` commits are intentionally excluded, so use `feat`, `fix`,
`docs`, `refactor`, `test`, and `chore` for releasable work.
