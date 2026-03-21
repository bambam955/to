#!/bin/bash
# Integration tests for the bash wrapper (src/wrappers/to.bash)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Use the real built backend
export PATH="${REPO_ROOT}/bin:${PATH}"

# Point at canned test database (use a copy so --clean doesn't mutate it)
TEST_DB="$(mktemp)"
cp "${SCRIPT_DIR}/testdata/database.json" "${TEST_DB}"
export TO_DB="${TEST_DB}"
trap 'rm -f "${TEST_DB}"' EXIT

# Source the wrapper
# shellcheck disable=SC1091
source "${REPO_ROOT}/src/wrappers/to.bash"

pass=0
fail=0

assert() {
    local name="$1"
    local expected="$2"
    local actual="$3"
    if [[ "${expected}" == "${actual}" ]]; then
        echo "  PASS: ${name}"
        ((pass++)) || true
    else
        echo "  FAIL: ${name}"
        echo "    expected: ${expected}"
        echo "    actual:   ${actual}"
        ((fail++)) || true
    fi
}

echo "=== Bash wrapper integration tests ==="

# Test 1: Navigation changes directory
echo "Test 1: navigation changes PWD"
start_dir="${PWD}"
to tmp
assert "PWD changed to /tmp" "/tmp" "${PWD}"
cd "${start_dir}" || exit 1

# Test 2: Non-navigation output is passed through
echo "Test 2: list output passthrough"
output=$(to --list)
# shellcheck disable=SC2312
assert "list contains tmp alias" "yes" "$(echo "${output}" | grep -q 'tmp' && echo yes || echo no)"

# Test 3: Error preserves non-zero exit code
echo "Test 3: error exit code preserved"
# shellcheck disable=SC2310
if to nonexistent 2>/dev/null; then
    assert "exit code non-zero" "1" "0"
else
    assert "exit code non-zero" "non-zero" "non-zero"
fi

# Test 4: Expand shows path without triggering cd
echo "Test 4: expand shows path without cd"
start_dir="${PWD}"
output=$(to --exp tmp)
assert "exp output is /tmp" "/tmp" "${output}"
assert "PWD unchanged" "${start_dir}" "${PWD}"

# Test 5: to function is exported
echo "Test 5: function is exported"
if declare -F to >/dev/null 2>&1; then
    assert "to function exists" "yes" "yes"
else
    assert "to function exists" "yes" "no"
fi

# Summary
echo ""
echo "Results: ${pass} passed, ${fail} failed"
[[ ${fail} -eq 0 ]]
