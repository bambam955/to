#!/usr/bin/env fish
# Fish wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

function to
    # Capture only fd 3 control frames and let stdout/stderr pass through.
    set -l control_file (mktemp)
    to-backend $argv 3>"$control_file"
    set -l exit_code $status
    set -l control_output (cat "$control_file")
    rm -f "$control_file"

    # Apply navigation on successful backend runs with a NAV control frame.
    if test "$exit_code" -eq 0; and string match -rq '^NAV (.+)$' -- "$control_output"
        set -l target_dir (string match -r '^NAV (.+)$' -- "$control_output")[2]
        cd "$target_dir"; or return 1
        return 0
    end

    # A non-empty unmatched frame indicates protocol mismatch.
    if test "$exit_code" -eq 0; and test -n "$control_output"
        echo "error: invalid navigation control frame: $control_output" >&2
        return 1
    end

    return $exit_code
end
