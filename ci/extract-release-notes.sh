#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${repo_root}"

usage() {
    echo "Usage: $0 <version> <output-path>" >&2
    exit 1
}

if [[ "$#" -ne 2 ]]; then
    usage
fi

version="$1"
output_path="$2"

if [[ ! "${version}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "error: version must be bare semantic versioning (for example 1.2.3)" >&2
    exit 1
fi

# GitHub Releases should use the exact checked-in changelog section for the
# matching tag so release notes stay aligned with the release PR process.
awk -v version="${version}" '
BEGIN {
    capture = 0
}
$0 ~ ("^## \\[" version "\\]") {
    capture = 1
}
capture && /^## \[/ && $0 !~ ("^## \\[" version "\\]") {
    exit
}
capture && /^<!--/ {
    exit
}
capture {
    print
}
' CHANGELOG.md >"${output_path}"

if [[ ! -s "${output_path}" ]]; then
    echo "error: could not find changelog notes for ${version} in CHANGELOG.md" >&2
    exit 1
fi
