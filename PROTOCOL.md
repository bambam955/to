# Communication Protocol: Bash Wrapper ↔ Go Backend

## Overview

This document defines the communication protocol between the bash wrapper function `to()` and the Go backend binary `to-backend`. The protocol enables synchronous, request-response communication with clear separation between navigation commands and regular output.

## Architecture

```
┌──────────────────┐
│  User Terminal   │
│   Shell Session  │
└────────┬─────────┘
         │
         ▼
┌──────────────────────────┐
│  Bash Wrapper to()       │
│  - Parses responses      │
│  - Executes cd           │
│  - Forwards errors       │
└──────────┬───────────────┘
           │
      Process Call
      (stdin, stdout, stderr)
           │
           ▼
┌──────────────────────────┐
│   Go Backend to-backend  │
│   - Database operations  │
│   - Alias management     │
│   - Validation logic     │
└──────────────────────────┘
```

## Response Formats

### 1. Navigation Success Response

**Format**: `[to] <absolute_path>`

**Description**: Used exclusively for successful navigation commands to indicate the directory to change to.

**Examples**:
```
[to] /home/user/projects
[to] /home/user/work/project-x
[to] /
```

**Parsing**:
The bash wrapper uses regex matching to detect this pattern:
```bash
if [[ "$output" =~ ^\[to\]\ (.+)$ ]]; then
    target_dir="${BASH_REMATCH[1]}"
    cd "$target_dir"
fi
```

**Key Properties**:
- Exactly 5 characters: `[to] ` (note the space)
- Followed by the absolute path to navigate to
- Printed to stdout
- Always indicates successful navigation
- Cannot be mixed with other output

### 2. Error Response

**Format**: `error: <message>`

**Description**: All error conditions produce messages with the "error: " prefix.

**Examples**:
```
error: alias not found
error: directory /nonexistent does not exist
error: invalid alias name: contains invalid characters
error: permission denied: cannot access /root
error: database corrupted: invalid JSON format
```

**Delivery**: Printed to stderr (not parsed by wrapper, just forwarded)

**Key Properties**:
- Exactly 7 characters: `error: ` (note the space)
- Followed by a human-readable error message
- Printed to stderr
- Forwarded unchanged by wrapper
- Multiple errors may be output (one per line)

### 3. Non-Navigation Output (Info/Status Messages)

**Format**: Regular text output without the `[to]` prefix or `error:` prefix

**Description**: Used for `list`, `clean`, and other non-navigation commands that provide information.

**Examples**:
```
alias1     /home/user/projects/alias1
alias2     /home/user/projects/alias2
/home/user/projects
Removed 3 invalid aliases.
```

**Delivery**: Printed to stdout

**Key Properties**:
- No special parsing by wrapper
- Displayed as-is to the user
- Used for informational commands
- May span multiple lines

## Exit Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 0 | Success | All successful commands return 0 |
| 1 | Error | Any error condition returns 1 |

**Notes**:
- Exit code is not used to distinguish between navigation success and other outputs; use the response format instead
- Bash wrapper preserves the exit code from the backend
- Exit code semantics are simple: 0 = success, 1 = failure

## Command Examples

### Navigation Command

**Bash Wrapper Call**:
```bash
to myalias
```

**Backend Execution**:
```bash
to-backend myalias
```

**Successful Response**:
```
stdout: [to] /home/user/projects
stderr: (empty)
exit code: 0
```

**Wrapper Behavior**:
```bash
cd /home/user/projects
```

**Failed Response**:
```
stdout: (empty)
stderr: error: alias not found
exit code: 1
```

**Wrapper Behavior**:
```bash
echo "error: alias not found"
return 1
```

### Registration Command

**Bash Wrapper Call**:
```bash
to reg myalias /home/user/projects
```

**Backend Execution**:
```bash
to-backend reg myalias /home/user/projects
```

**Successful Response**:
```
stdout: Registered alias 'myalias' → /home/user/projects
stderr: (empty)
exit code: 0
```

**Failed Response** (duplicate):
```
stdout: (empty)
stderr: error: alias 'myalias' already exists
exit code: 1
```

### List Command

**Bash Wrapper Call**:
```bash
to list
```

**Backend Execution**:
```bash
to-backend list
```

**Response**:
```
stdout: 
alias1     /home/user/projects
alias2     /home/user/documents
myalias    /home/user/work
stderr: (empty)
exit code: 0
```

### Expand Command

**Bash Wrapper Call**:
```bash
to exp myalias
```

**Backend Execution**:
```bash
to-backend exp myalias
```

**Successful Response**:
```
stdout: /home/user/projects
stderr: (empty)
exit code: 0
```

**Failed Response**:
```
stdout: (empty)
stderr: error: alias not found
exit code: 1
```

## Process Lifecycle

1. **User Input**: User types command in shell
2. **Wrapper Invocation**: Bash wrapper function `to()` is called
3. **Backend Spawn**: Wrapper spawns `to-backend` process with arguments
4. **Processing**: Backend performs operation, writes to stdout/stderr
5. **Process Exit**: Backend exits with appropriate code
6. **Response Parsing**: Wrapper checks output format and acts accordingly
7. **User Feedback**: Wrapper either changes directory or displays output

**Important**: All operations are synchronous and blocking. The wrapper waits for the backend to complete before returning.

## Error Handling Strategy

### Backend Responsibilities

1. **Validate Input**: Check arguments, formats, database state
2. **Categorize Errors**: Determine error type (validation, database, operation)
3. **Output Errors**: Write clear error messages to stderr with "error: " prefix
4. **Exit with 1**: Return exit code 1 for any error

### Wrapper Responsibilities

1. **Parse Navigation**: Check for `[to] ` prefix in stdout
2. **Execute cd**: Change directory if navigation response detected
3. **Forward Errors**: Print stderr unchanged (already has "error: " prefix)
4. **Preserve Exit Code**: Return the same exit code from backend

## Special Considerations

### Shell Expansion
Paths may contain:
- Spaces: `/home/user/my documents` → preserved in `[to]` response
- Special characters: `$`, `!`, etc. → passed through stdout
- Quotes: handled by shell before wrapper receives them

### Database Consistency
- Each command is atomic: reads or modifies database in a single operation
- No partial updates: transaction either completes or fails entirely
- Stateless processes: each command is independent

### Multi-Line Outputs
- Navigation: Always single line with `[to] <path>` format
- Errors: May be single or multiple lines, each with `error: ` prefix
- Other output: May span multiple lines without special prefix

## Implementation Notes

### Bash Wrapper Implementation
```bash
to() {
    local output
    local exit_code
    
    # Capture stdout and exit code
    output=$("to-backend" "$@")
    exit_code=$?
    
    # Check for navigation response format
    if [[ "$output" =~ ^\[to\]\ (.+)$ ]]; then
        local target_dir="${BASH_REMATCH[1]}"
        cd "$target_dir"
        return 0
    fi
    
    # Forward all other output
    echo "$output"
    return $exit_code
}
```

### Go Backend Error Formatting
```go
import "fmt"

// Navigation success
fmt.Println(protocol.NavigationResponse("/path/to/dir"))

// Error output
fmt.Fprintf(os.Stderr, "%s\n", protocol.ErrorResponse("error message"))
os.Exit(1)

// Other output
fmt.Println("Some information")
```

## Backwards Compatibility

This protocol is designed to be stable and extensible:

- The `[to]` prefix is unlikely to conflict with normal directory paths
- Error prefix `error: ` is standard and extensible
- Exit codes follow Unix conventions
- No binary protocols or special encoding required
- Human-readable format aids debugging

If future enhancements are needed:
- New response types can use different prefixes (e.g., `[to-info]`)
- Additional exit codes could be added without breaking compatibility
- The protocol does not require backend version matching
