#!/bin/bash
# Bash wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

to() {
    local control_file
    local control_output
    local exit_code

    # Capture only the machine control channel (fd 3). Stdout/stderr are
    # passed through directly so interactive output keeps native TTY behavior.
    control_file="$(mktemp)"
    "to-backend" "$@" 3>"${control_file}"
    exit_code=$?
    control_output="$(cat "${control_file}")"
    rm -f "${control_file}"

    # Apply navigation response on a successful command when a NAV frame exists.
    if [[ "${exit_code}" -eq 0 && "${control_output}" =~ ^NAV\ (.+)$ ]]; then
        # Extract the path from the response and change directory
        local target_dir="${BASH_REMATCH[1]}"
        cd "${target_dir}" || return 1
        return 0
    fi

    # A non-empty control payload with no valid frame indicates protocol drift.
    if [[ "${exit_code}" -eq 0 && -n "${control_output}" ]]; then
        echo "error: invalid navigation control frame: ${control_output}" >&2
        return 1
    fi

    return "${exit_code}"
}

# Export the function so it's available in subshells
export -f to
