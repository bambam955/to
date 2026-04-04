#!/usr/bin/env bash

set -euo pipefail

# Resolve the repository root once so branch detection, fetches, and output
# paths behave the same regardless of the caller's current directory.
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${repo_root}"

usage() {
    cat >&2 <<'EOF'
Usage: ci/generate-changelog.sh [--branch <name>] [--tag <version>] [--release-date <YYYY-MM-DD>] [--output <path>]
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
release_date_override=""
output_path="CHANGELOG.md"
release_version_regex='^[0-9]+\.[0-9]+\.[0-9]+$'
release_date_regex='^[0-9]{4}-[0-9]{2}-[0-9]{2}$'

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
    --release-date)
        [[ "$#" -ge 2 ]] || usage
        release_date_override="$2"
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

date_to_epoch() {
    local release_date="$1"

    # Support both GNU date (Linux CI) and BSD date (local macOS usage).
    if date -u -d "${release_date}T00:00:00Z" +%s >/dev/null 2>&1; then
        date -u -d "${release_date}T00:00:00Z" +%s
        return 0
    fi

    if date -j -u -f "%Y-%m-%dT%H:%M:%SZ" "${release_date}T00:00:00Z" +%s >/dev/null 2>&1; then
        date -j -u -f "%Y-%m-%dT%H:%M:%SZ" "${release_date}T00:00:00Z" +%s
        return 0
    fi

    return 1
}

validate_release_date() {
    local value="$1"
    local source_description="$2"
    local epoch_value=""
    local epoch_status=0

    if [[ ! "${value}" =~ ${release_date_regex} ]]; then
        echo "error: ${source_description} must use YYYY-MM-DD format (got: ${value})" >&2
        exit 1
    fi

    set +e
    epoch_value="$(date_to_epoch "${value}" 2>/dev/null)"
    epoch_status=$?
    set -e

    if [[ "${epoch_status}" -ne 0 || -z "${epoch_value}" ]]; then
        echo "error: ${source_description} is not a valid calendar date (got: ${value})" >&2
        exit 1
    fi
}

resolve_branch_ref() {
    local requested_branch="$1"
    local current_branch=""

    # Prefer a local branch ref so trailer lookups match the branch state that
    # release CI and local release-prep commands actually operate on.
    if git show-ref --verify --quiet "refs/heads/${requested_branch}"; then
        echo "refs/heads/${requested_branch}"
        return 0
    fi

    if git show-ref --verify --quiet "refs/remotes/origin/${requested_branch}"; then
        echo "refs/remotes/origin/${requested_branch}"
        return 0
    fi

    current_branch="$(git branch --show-current || true)"
    if [[ "${current_branch}" == "${requested_branch}" ]]; then
        echo "HEAD"
        return 0
    fi

    echo "error: branch ref not found locally for release-date lookup: ${requested_branch}" >&2
    exit 1
}

find_latest_release_date_trailer() {
    local branch_ref="$1"

    # The newest release-date trailer wins, which gives release branches a
    # deterministic correction path when the intended release day changes.
    git log "${branch_ref}" --format='%H %(trailers:key=Release-Date,valueonly)' |
        awk '$2 != "" {print $1 "\t" $2; exit}'
}

# Normalize the output path up front so git-cliff can write to the same target
# whether it runs in the current checkout or a temporary mainline worktree.
if [[ "${output_path}" = /* ]]; then
    output_abs="${output_path}"
else
    output_abs="${repo_root}/${output_path}"
fi

# Release branches need an explicit tag value so git-cliff renders the pending
# release heading instead of falling back to an unreleased section. Allow
# callers like `just gen-changelog` to stay argument-free by deriving the
# version from the branch name when the caller did not override `--tag`.
if [[ -z "${tag_value}" ]]; then
    case "${branch_name}" in
    release/*)
        inferred_version="${branch_name#release/}"
        if [[ ! "${inferred_version}" =~ ${release_version_regex} ]]; then
            echo "error: release branch must end with a bare semantic version (got: ${inferred_version})" >&2
            exit 1
        fi
        tag_value="${inferred_version}"
        ;;
    *) ;;
    esac
fi

if [[ -n "${release_date_override}" && "${branch_name}" != release/* ]]; then
    echo "error: --release-date is only supported for release/* branches" >&2
    exit 1
fi

release_date_value=""
if [[ "${branch_name}" == release/* && -n "${tag_value}" ]]; then
    if ! git rev-parse -q --verify "refs/tags/${tag_value}" >/dev/null 2>&1; then
        if [[ -n "${release_date_override}" ]]; then
            validate_release_date "${release_date_override}" "--release-date"
            release_date_value="${release_date_override}"
        else
            branch_ref="$(resolve_branch_ref "${branch_name}")"
            trailer_entry="$(find_latest_release_date_trailer "${branch_ref}")"

            if [[ -z "${trailer_entry}" ]]; then
                echo "error: pending release ${branch_name} requires a Release-Date trailer or --release-date" >&2
                exit 1
            fi

            IFS=$'\t' read -r trailer_commit release_date_value <<<"${trailer_entry}"
            validate_release_date "${release_date_value}" "Release-Date trailer in ${trailer_commit}"
        fi
    elif [[ -n "${release_date_override}" ]]; then
        echo "error: --release-date cannot override an existing tag timestamp for ${tag_value}" >&2
        exit 1
    fi
fi

run_git_cliff() {
    local workdir="$1"
    shift

    if [[ -n "${release_date_value}" ]]; then
        local context_path
        local rendered_context_path
        local release_epoch

        require_command python3

        context_path="$(mktemp)"
        rendered_context_path="$(mktemp)"
        release_epoch="$(date_to_epoch "${release_date_value}")"

        (
            cd "${workdir}"
            git-cliff --config "${config_path}" "$@" --context >"${context_path}"
        )

        # Re-render the latest release from explicit JSON context so pending
        # release branches use the committed release date instead of wall clock
        # generation time. Tagged releases still use git metadata directly.
        python3 - "${context_path}" "${rendered_context_path}" "${release_epoch}" <<'PY'
import json
import sys

context_path, rendered_context_path, release_epoch = sys.argv[1], sys.argv[2], int(sys.argv[3])

with open(context_path, encoding="utf-8") as context_file:
    context = json.load(context_file)

if not context:
    raise SystemExit("error: git-cliff produced an empty changelog context")

context[0]["timestamp"] = release_epoch

with open(rendered_context_path, "w", encoding="utf-8") as rendered_context_file:
    json.dump(context, rendered_context_file, separators=(",", ":"))
PY

        git-cliff --config "${config_path}" --from-context "${rendered_context_path}" --output "${output_abs}"
        rm -f "${context_path}" "${rendered_context_path}"
        return 0
    fi

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
