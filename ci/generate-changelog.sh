#!/usr/bin/env bash

set -euo pipefail

# Resolve the repository root once so branch detection, fetches, and output
# paths behave the same regardless of the caller's current directory.
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${repo_root}"

usage() {
    cat >&2 <<'EOF'
Usage: ci/generate-changelog.sh [--branch <name>] [--tag <version>] [--output <path>]
EOF
    exit 1
}

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "error: required command not found: $1" >&2
        exit 1
    fi
}

branch_name="${GITHUB_HEAD_REF:-${GITHUB_REF_NAME:-$(git branch --show-current)}}"
tag_value=""
output_path="CHANGELOG.md"

while [[ "$#" -gt 0 ]]; do
    case "$1" in
    --branch)
        [[ "$#" -ge 2 ]] || usage
        branch_name="$2"
        shift 2
        ;;
    --tag)
        [[ "$#" -ge 2 ]] || usage
        tag_value="$2"
        shift 2
        ;;
    --output)
        [[ "$#" -ge 2 ]] || usage
        output_path="$2"
        shift 2
        ;;
    *)
        usage
        ;;
    esac
done

require_command git
require_command git-cliff

config_path="${repo_root}/cliff.toml"

if [[ ! -f "${config_path}" ]]; then
    echo "error: cliff.toml not found at ${config_path}" >&2
    exit 1
fi

# Normalize the output path up front so git-cliff can write to the same target
# whether it runs in the current checkout or a temporary mainline worktree.
if [[ "${output_path}" = /* ]]; then
    output_abs="${output_path}"
else
    output_abs="${repo_root}/${output_path}"
fi

run_git_cliff() {
    local workdir="$1"
    shift

    (
        cd "${workdir}"
        git-cliff --config "${config_path}" "$@" --output "${output_abs}"
    )
}

temp_worktree=""
cleanup() {
    if [[ -n "${temp_worktree}" ]]; then
        git -C "${repo_root}" worktree remove --force "${temp_worktree}" >/dev/null 2>&1 || true
    fi
}
trap cleanup EXIT

git_cliff_args=()
if [[ -n "${tag_value}" ]]; then
    git_cliff_args+=(--tag "${tag_value}")
fi

case "${branch_name}" in
main | release/*)
    # Main and release branches are allowed to generate from their current
    # checked-out state because their commits are part of the intended
    # changelog surface.
    run_git_cliff "${repo_root}" "${git_cliff_args[@]}"
    ;;
*)
    # Feature branches should render the changelog against the latest
    # mainline history so local implementation commits never leak into the
    # checked-in file.
    git fetch --quiet origin main

    temp_worktree="$(mktemp -d)"
    rmdir "${temp_worktree}"
    git worktree add --detach --quiet "${temp_worktree}" refs/remotes/origin/main

    run_git_cliff "${temp_worktree}" "${git_cliff_args[@]}"
    ;;
esac
