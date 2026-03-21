#!/bin/zsh
# Zsh wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

to() {
    local output
    local exit_code

    # Call the Go backend and capture both stdout and exit code
    output=$("to-backend" "$@")
    exit_code=$?

    # Check if this is a navigation response (starts with "[to] ")
    if [[ "${output}" =~ '^\[to\] (.+)$' ]]; then
        # Extract the path from the response and change directory
        local target_dir="${match[1]}"
        cd "${target_dir}" || return 1
        return 0
    fi

    # For non-navigation commands, output as-is and return the exit code
    if [[ -n "${output}" ]]; then
        echo "${output}"
    fi
    return "${exit_code}"
}
