# Repository Guidelines

## Project Structure

- `cmd/lametric/`: CLI entrypoint
- `internal/`: implementation
  - `cmd/`: command routing (kong CLI framework)
  - `api/`: LaMetric device API client
  - `auth/`: keyring-based credential storage
  - `config/`: YAML configuration
  - `discovery/`: mDNS/SSDP device discovery
  - `output/`: terminal rendering (lipgloss)
  - `errfmt/`: error formatting
- `bin/`: build outputs

## Build, Test, and Development Commands

- `make build`: compile to `bin/lametric`
- `make lametric -- <args>`: build + run (e.g., `make lametric -- notify "Hello"`)
- `make fmt` / `make lint` / `make test` / `make ci`: format, lint, test, full local gate
- `make tools`: install pinned dev tools into `.tools/`
- `make clean`: remove bin/ and .tools/

## Coding Style & Naming Conventions

- Formatting: `make fmt` (goimports local prefix `github.com/dedene/lametric-cli` + gofumpt)
- Output: keep stdout parseable (`--json`); send human hints/progress to stderr
- Linting: golangci-lint v2.8.0 with project config

## Testing Guidelines

- Unit tests: stdlib `testing` (files: `*_test.go` next to code)
- 11 test files; comprehensive coverage
- CI gate: fmt-check, lint, test

## Config & Secrets

- **Keyring**: 99designs/keyring for API key storage (macOS Keychain, Linux SecretService, file backend)
- **Env var fallback**: `LAMETRIC_API_KEY`
- **Keyring backend**: `LAMETRIC_KEYRING_BACKEND` (auto/keychain/file)
- **Keyring password**: `LAMETRIC_KEYRING_PASSWORD` (for file backend)
- **Device discovery**: automatic via mDNS/SSDP on local network

## Key Features

- Send notifications to LaMetric TIME displays
- Device discovery on local network
- Multiple device support

## Commit & Pull Request Guidelines

- Conventional Commits: `feat|fix|refactor|build|ci|chore|docs|style|perf|test`
- Group related changes; avoid bundling unrelated refactors
- PR review: use `gh pr view` / `gh pr diff`; don't switch branches

## Security Tips

- Never commit API keys
- Prefer OS keychain; file backend only for headless environments
- Device API keys should be unique per device
