#!/bin/bash
# Integration test runner script

set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <shell>" >&2
    exit 1
fi

readonly SHELL=$1

if [[ $SHELL == "all" ]]; then
    for s in bash zsh fish; do
        command -v $s >/dev/null 2>&1 && $s tests/integration/${s}_test.${s} || echo "SKIP: $s not found"
    done
else
    $SHELL "tests/integration/${SHELL}_test.${SHELL}"
fi
