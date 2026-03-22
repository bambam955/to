#!/bin/zsh
# Zsh wrapper for the 'to' directory navigation tool
# This function intercepts the 'to' command and routes it to the Go backend

to() {
    local control_file
    local control_output
    local exit_code

    # Capture only fd 3 control frames; stdout/stderr flow through directly
    # so terminal formatting behavior remains accurate.
    control_file="$(mktemp)"
    "to-backend" "$@" 3>"${control_file}"
    exit_code=$?
    control_output="$(cat "${control_file}")"
    rm -f "${control_file}"

    # Apply navigation only when backend succeeded and emitted NAV payload.
    if [[ "${exit_code}" -eq 0 && "${control_output}" =~ '^NAV (.+)$' ]]; then
        # Extract the path from the response and change directory
        local target_dir="${match[1]}"
        cd "${target_dir}" || return 1
        return 0
    fi

    # Guard against malformed control payloads to fail loudly on protocol drift.
    if [[ "${exit_code}" -eq 0 && -n "${control_output}" ]]; then
        echo "error: invalid navigation control frame: ${control_output}" >&2
        return 1
    fi

    return "${exit_code}"
}
