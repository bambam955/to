#!/usr/bin/env fish
# Integration tests for the fish wrapper (wrappers/to.fish)

set SCRIPT_DIR (status dirname)
set REPO_ROOT "$SCRIPT_DIR/../.."

# Use the real built backend
set -x PATH "$REPO_ROOT/bin" $PATH

# Point at canned test database (use a copy so --clean doesn't mutate it)
set -x TO_DB (mktemp)
cp "$SCRIPT_DIR/testdata/database.json" "$TO_DB"

# Source the wrapper
source "$REPO_ROOT/wrappers/to.fish"

set pass 0
set fail 0

function assert
    set -l name $argv[1]
    set -l expected $argv[2]
    set -l actual $argv[3]
    if test "$expected" = "$actual"
        echo "  PASS: $name"
        set pass (math $pass + 1)
    else
        echo "  FAIL: $name"
        echo "    expected: $expected"
        echo "    actual:   $actual"
        set fail (math $fail + 1)
    end
end

echo "=== Fish wrapper integration tests ==="

# Test 1: Navigation changes directory
echo "Test 1: navigation changes PWD"
set start_dir $PWD
to tmp
assert "PWD changed to /tmp" "/tmp" "$PWD"
cd $start_dir; or exit 1

# Test 2: Non-navigation output is passed through
echo "Test 2: list output passthrough"
set output (to --list)
set has_tmp (string match -r 'tmp' "$output" | head -1)
assert "list contains tmp alias" "tmp" "$has_tmp"

# Test 3: Error preserves non-zero exit code
echo "Test 3: error exit code preserved"
to nonexistent 2>/dev/null
if test $status -ne 0
    assert "exit code non-zero" "non-zero" "non-zero"
else
    assert "exit code non-zero" "1" "0"
end

# Test 4: Expand shows path without triggering cd
echo "Test 4: expand shows path without cd"
set start_dir $PWD
set output (to --exp tmp)
assert "exp output is /tmp" "/tmp" "$output"
assert "PWD unchanged" "$start_dir" "$PWD"

# Cleanup
rm -f "$TO_DB"

# Summary
echo ""
echo "Results: $pass passed, $fail failed"
test $fail -eq 0
