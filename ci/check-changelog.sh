#!/usr/bin/env bash

set -euo pipefail

# Run from the repository root so local and CI behavior stay identical.
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${repo_root}"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "error: required command not found: $1" >&2
        exit 1
    fi
}

require_command git-cliff

branch_name="${1:-${GITHUB_HEAD_REF:-${GITHUB_REF_NAME:-$(git branch --show-current)}}}"

case "${branch_name}" in
release/*)
    version="${branch_name#release/}"
    ;;
*)
    echo "error: changelog checks only support release/* branches (got: ${branch_name})" >&2
    exit 1
    ;;
esac

if [[ ! "${version}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "error: release branch must end with a bare semantic version (got: ${version})" >&2
    exit 1
fi

if [[ ! -f CHANGELOG.md ]]; then
    echo "error: CHANGELOG.md not found" >&2
    exit 1
fi

tmp_changelog="$(mktemp)"
trap 'rm -f "${tmp_changelog}"' EXIT

# Regenerate the changelog exactly as the release workflow expects and compare
# the result to the checked-in file.
bash ci/generate-changelog.sh --branch "${branch_name}" --tag "${version}" --output "${tmp_changelog}"

if ! diff -u CHANGELOG.md "${tmp_changelog}"; then
    echo "error: CHANGELOG.md is out of date for release ${version}" >&2
    # Existing release branches only need the changelog regenerated in place.
    # `prep-release` creates a new branch and is intentionally not reusable here.
    echo "run: just gen-changelog" >&2
    exit 1
fi

echo "CHANGELOG.md is up to date for release ${version}"
