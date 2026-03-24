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
3. regenerates `CHANGELOG.md` for that exact release version
4. commits the changelog update
5. pushes the release branch
6. opens the release PR against `main`

Example:

```bash
just prep-release 1.2.3
```

## Generate the Changelog Locally

`just gen-changelog` behaves differently depending on the current branch:

- on `main`, it uses the current `main` checkout
- on `release/*`, it uses the current release branch
- on any other branch, it fetches `origin/main` and renders from mainline history

That last rule keeps local implementation-branch commits from leaking into the
checked-in changelog.

## Publish the Release

After the release PR is merged, cut the release by pushing the matching bare
tag:

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
