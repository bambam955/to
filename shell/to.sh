#!/bin/bash
# Bash wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

to() {
    local output
    local exit_code
    
    # Call the Go backend and capture both stdout and exit code
    output=$("to-backend" "$@")
    exit_code=$?
    
    # Check if this is a navigation response (starts with "[to] ")
    if [[ "${output}" =~ ^\[to\]\ (.+)$ ]]; then
        # Extract the path from the response and change directory
        local -r target_dir="${BASH_REMATCH[1]}"
        cd "${target_dir}" || {
            unset target_dir
            exit 1
        }
        unset target_dir
        return 0
    fi
    
    # For non-navigation commands, output as-is and return the exit code
    echo "${output}"
    return "${exit_code}"
}

# Export the function so it's available in subshells
export -f to
