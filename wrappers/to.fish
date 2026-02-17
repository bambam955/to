#!/usr/bin/env fish
# Fish wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

function to
    set -l output (to-backend $argv)
    set -l exit_code $status

    # Check if this is a navigation response (starts with "[to] ")
    if string match -rq '^\[to\] (.+)$' -- $output
        # Extract the path from the response and change directory
        set -l target_dir (string match -r '^\[to\] (.+)$' -- $output)[2]
        cd $target_dir; or return 1
        return 0
    end

    # For non-navigation commands, output as-is and return the exit code
    if test -n "$output"
        echo $output
    end
    return $exit_code
end
