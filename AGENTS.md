# AGENTS.md

## Build & Test
- Build: `just build` (or `go build -o bin/to-backend ./cmd/to-backend`)
- Test all: `just test` (or `go test ./...`)
- Test single package: `go test ./pkg/database/` — single test: `go test ./pkg/database/ -run TestName`
- Lint: `just lint` (`go vet ./...` + `shellcheck --enable=all --shell=bash shell/*`)
- Format: `just fmt` (`gofmt -w .` + `shfmt -i 4 -w .`)
- Full dev cycle: `just dev` (clean → fmt → lint → test → build)

## Architecture
Go CLI tool (`to`) for instant directory navigation. Two-part design:
- **`shell/to.sh`**: Bash wrapper that calls the backend and performs `cd` (shell can't change its own cwd from a subprocess).
- **`cmd/to-backend/`**: Go binary (cobra CLI) — the backend that resolves aliases.
- **`pkg/`**: Core packages — `config/` (XDG paths), `database/` (JSON alias store), `errors/` (typed errors with `ErrorType`), `protocol/` (stdout `[to] <path>` format parsed by shell wrapper), `install/`.
- Config/data stored in `~/.config/to/` (XDG-compliant).

## Code Style
- Go module name: `to`. Imports use `to/pkg/...`.
- Use `pkg/errors` structured errors (`errors.New`, `errors.Wrap`, `errors.NotFound`, etc.) — not bare `fmt.Errorf`.
- Format with `gofmt`; shell scripts with `shfmt -i 4`. Lint shell with `shellcheck --enable=all`.
- All packages have `_test.go` files; keep tests alongside source.
- Pre-commit hooks enforce formatting, vetting, and shellcheck. Specs live in `specs/`.
