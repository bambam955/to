#!/bin/bash
# Integration test runner script

set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <shell>" >&2
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

readonly SHELL=$1

if [[ ${SHELL} == "all" ]]; then
    for s in bash zsh fish; do
        command -v "${s}" >/dev/null 2>&1 && "${s}" "${SCRIPT_DIR}/${s}_test.${s}" || echo "SKIP: ${s} not found"
    done
else
    "${SHELL}" "${SCRIPT_DIR}/${SHELL}_test.${SHELL}"
fi
