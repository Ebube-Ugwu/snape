# Snape

Snape is a local-first snippet manager written in Go. It stores snippets in SQLite, supports template variables like `{{YEAR}}`, exposes a CLI, and serves an embedded React web UI from the same Go binary.

## Build and Run

```sh
go run ./cmd list
make build
./dist/snape serve
```

The web UI runs at `http://localhost:7777` by default. Use `snape serve --addr :8080` to choose another port.

## Installation

The easiest install path is the release installer. It downloads the latest GitHub Release binary for your OS and installs it to `~/.local/bin`:

```sh
curl -fsSL https://raw.githubusercontent.com/Ebube-Ugwu/snape/master/scripts/install.sh | sh
```

Make sure `~/.local/bin` is on your `PATH`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

You can also download a binary manually from the GitHub Releases page, rename it to `snape`, make it executable, and place it anywhere on your `PATH`.

To build from source instead:

```sh
make build
install -m 755 dist/snape ~/.local/bin/snape
```

## CLI Commands

### Add a Snippet

```sh
snape add license --description "MIT header" --language text --type license --tag legal --content "Copyright {{YEAR}} {{AUTHOR}}"
snape add helper --file ./helper.go
cat ./snippet.txt | snape add pasted
```

If `--content`, `--file`, and stdin are not provided, Snape opens `$EDITOR`.

### Get a Snippet

```sh
snape get license --var YEAR=2026 --var AUTHOR=Ebube
```

`get` prints the snippet to stdout and renders `{{VARIABLE}}` placeholders from `--var KEY=VALUE` arguments.

### List and Search

```sh
snape list
snape list --verbose
snape search license
```

`list` shows a bordered, colorized table with names, languages, and tags. Use `-v` or `--verbose` to include type and description. `search` matches name, description, content, and language.

### Delete and Tag

```sh
snape delete license
snape tag helper go
```

Tags are stored in SQLite and linked to snippets.

### Insert Into a File

```sh
snape insert license README.md
snape insert license README.md --line 1 --var YEAR=2026 --var AUTHOR=Ebube
```

Without `--line`, Snape appends to the file. With `--line`, it inserts before the 1-based line number.

### Import and Export JSON

```sh
snape export --output snippets.json
snape export snippets.json
snape import snippets.json
```

`export` writes all snippets, including tags, to JSON. `import` accepts Snape's exported JSON shape or a raw JSON array of snippets. Existing snippets with the same name are updated.

### Serve the Web UI

```sh
snape serve
snape serve --addr :8080
```

The React assets are embedded in the Go binary and served directly by `net/http`.

## Development

```sh
go test ./...
gofmt -w cmd internal
go mod tidy
make build
```

Core code lives in `internal/data`, the CLI entry point is `cmd/main.go`, and embedded web assets live in `internal/server/static`.

## License

Snape is licensed under the MIT License. See `LICENSE`.
