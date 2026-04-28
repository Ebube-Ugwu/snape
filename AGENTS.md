# Repository Guidelines

## Project Structure

Snape is a Go module (`github.com/ebube-ugwu/snape`) for a local-first snippet manager. `PRD.md` describes the product direction; `README.md` documents current usage.

- `cmd/main.go` contains the CLI entry point and subcommand wiring.
- `internal/data/` contains SQLite setup, migrations, store methods, templating helpers, and tests.
- `internal/data/migrations/` contains embedded SQL migrations. Use paired `NNNN_description.up.sql` and `NNNN_description.down.sql` files.
- `internal/server/` serves the web UI and JSON API.
- `internal/server/static/` contains embedded React-style browser assets.
- `scripts/install.sh` installs the latest GitHub Release binary into `~/.local/bin`.

Treat `snape.db` as local runtime data, not source.

## Commands

- `go run ./cmd list` runs Snape from source.
- `go run ./cmd serve --addr :7777` starts the embedded web UI.
- `make build` builds a local binary at `dist/snape`.
- `go test ./...` runs all tests.
- `gofmt -w cmd internal` formats Go code.
- `go mod tidy` cleans dependency metadata after dependency changes.

When running in restricted environments, set `GOCACHE=/tmp/snape-go-build`.

## Coding Style

Use standard Go formatting and small lowercase package names. Keep new functionality close to the existing boundaries: CLI behavior in `cmd`, persistence and template logic in `internal/data`, and HTTP/UI behavior in `internal/server`. Continue using `database/sql` with `modernc.org/sqlite`; do not introduce CGo-only SQLite dependencies.

The web UI is embedded in the Go binary. Avoid adding a Node build step unless the project explicitly moves to one.

## Testing

Use Go's built-in `testing` package. Place `*_test.go` beside the code under test. Prefer temporary SQLite files or in-memory databases so tests never mutate `snape.db`. Add focused tests for store behavior, migrations, and template rendering when those areas change.

## Commits and PRs

Use concise Conventional Commit-style messages, matching the existing history: `feat:`, `fix:`, `test:`, `docs:`, and `chore:`. PRs should summarize the change, list tests run, and call out migration or CLI behavior changes.

## CI and Releases

GitHub Actions lives in `.github/workflows/ci.yml`. Keep it validating `go test ./...`, cross-building Linux, macOS, Windows, and Android artifacts with `CGO_ENABLED=0`, and publishing release assets for version tags.
