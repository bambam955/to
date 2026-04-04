#!/usr/bin/env bash

set -euo pipefail

# Resolve the repository root so the script behaves the same from any cwd.
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${repo_root}"

usage() {
    echo "Usage: $0 <version>" >&2
    exit 1
}

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "error: required command not found: $1" >&2
        exit 1
    fi
}

if [[ "$#" -ne 1 ]]; then
    usage
fi

version="$1"
release_branch="release/${version}"
release_regex='^[0-9]+\.[0-9]+\.[0-9]+$'
release_date="$(date +%Y-%m-%d)"

if [[ ! "${version}" =~ ${release_regex} ]]; then
    echo "error: version must be bare semantic versioning (for example 1.2.3)" >&2
    exit 1
fi

require_command git
require_command git-cliff
require_command gh

# Release prep should always start from a clean main checkout so the generated
# changelog is based on the exact release candidate state.
current_branch="$(git branch --show-current)"
if [[ "${current_branch}" != "main" ]]; then
    echo "error: prep-release must be run from main (current: ${current_branch})" >&2
    exit 1
fi

worktree_status="$(git status --porcelain)"
if [[ -n "${worktree_status}" ]]; then
    echo "error: working tree must be clean before preparing a release" >&2
    exit 1
fi

if git show-ref --verify --quiet "refs/heads/${release_branch}"; then
    echo "error: local branch already exists: ${release_branch}" >&2
    exit 1
fi

if git ls-remote --exit-code --heads origin "${release_branch}" >/dev/null 2>&1; then
    echo "error: remote branch already exists: ${release_branch}" >&2
    exit 1
fi

# Release branches must be cut from the exact remote mainline tip. Fetch first
# and fail fast if local main is ahead, behind, or diverged.
git fetch --quiet origin main

local_main_sha="$(git rev-parse HEAD)"
remote_main_sha="$(git rev-parse refs/remotes/origin/main)"
if [[ "${local_main_sha}" != "${remote_main_sha}" ]]; then
    echo "error: local main must match origin/main before preparing a release" >&2
    echo "hint: update main first (for example: git pull --ff-only)" >&2
    echo "local main:  ${local_main_sha}" >&2
    echo "origin/main: ${remote_main_sha}" >&2
    exit 1
fi

git switch -c "${release_branch}"

# Generate the full changelog for the requested release version so CI can
# reproduce it exactly from branch name + repo state. The release date is
# chosen once here and then persisted in the prep commit trailer so later
# reruns stay deterministic until the branch is merged and tagged.
bash ci/generate-changelog.sh --branch "${release_branch}" --tag "${version}" --release-date "${release_date}" --output CHANGELOG.md

git add CHANGELOG.md
git commit \
    -m "chore(changelog): prepare release ${version}" \
    -m "Release-Date: ${release_date}"
git push -u origin "${release_branch}"

gh pr create \
    --base main \
    --head "${release_branch}" \
    --title "release: ${version}" \
    --body "Prepare release ${version}.

- regenerate CHANGELOG.md for ${version}
- open the release PR before pushing the matching tag"

echo "Release branch created: ${release_branch}"
echo "Next steps:"
echo "  1. Review and merge the release PR"
echo "  2. Push the matching tag: git tag ${version} && git push origin ${version}"
