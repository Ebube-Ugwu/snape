Make db layer (migrations) and Models

Do basic subcommands

Do web frontend



# Snape — Snippet Manager: Full Project Plan

## The Idea in One Line
A fast, local-first CLI snippet manager with an embedded web UI, templating, and a future path to community snippet sharing.

---

## Phase 1 — Core Foundation

### Data Model (SQLite)

**snippets**
- `id`, `name`, `description`, `content`, `language`, `type` (file | inline | license | readme), `created_at`, `updated_at`

**tags**
- `id`, `name`

**snippet_tags**
- `snippet_id`, `tag_id`

**snippet_versions**
- `id`, `snippet_id`, `content`, `saved_at`

**variables**
- `id`, `snippet_id`, `key`, `default_value`

---

### CLI Command Structure

```
snape add <name>          # add a new snippet (opens $EDITOR)
snape get <name>          # print snippet to stdout
snape insert <name> <file> [--line N]  # insert into a file
snape list                # list all snippets
snape search <query>      # full-text search
snape edit <name>         # edit a snippet
snape delete <name>       # delete a snippet
snape tag <name> <tag>    # tag a snippet
snape history <name>      # view version history
snape export              # export library to JSON/ZIP
snape import <file>       # import a library
snape serve               # spin up the web UI
```

---

### Templating System

Snippets support `{{VARIABLE}}` placeholders. At insert time Snape prompts for values, or they can be passed as flags:

```
snape get my-license --var YEAR=2026 --var AUTHOR=Ebube
```

If no value is passed and a default exists in the DB, it uses that. Otherwise it prompts interactively.

---

## Phase 2 — Web UI

Spin up with `snape serve` on a configurable port (default `:7777`).

**Pages:**
- **Home/Dashboard** — search bar, recent snippets, tag cloud
- **Browse** — filterable grid of all snippets with language badges
- **Snippet Detail** — full content with syntax highlighting (highlight.js), copy button, variable list, version history
- **Add/Edit** — form with a code editor (CodeMirror)

**Aesthetic direction:** Dark, minimal, utilitarian — monospace-heavy, sharp edges, muted background with a single accent color. Think a terminal that learned design. Fits right alongside Eunotes.

---

## Phase 3 — Polish & Power Features

- **Shell completions** — `snape completion bash/zsh/fish` so snippet names autocomplete in the terminal
- **Watch mode** — `snape insert` detects if the target file is open and inserts cleanly
- **Config file** — `~/.snape/config.toml` for defaults (editor, port, author name, default year for licenses) yaml or toml
- **Syntax validation** — warn if a snippet's declared language doesn't match its content
- **Aliases** — short names for frequently used snippets

---

## Phase 4 — Community & Online Sharing

**Architecture:**
- A separate hosted Snape server (same Go codebase, different binary target)
- Users register and get a personal namespace: `snape.dev/ebube/my-license`
- Snippets can be **public**, **private**, or **unlisted**

**New CLI commands:**
```
snape login               # authenticate with snape.dev
snape publish <name>      # push local snippet online
snape pull <user/name>    # download someone's snippet locally (json)
snape search --remote <query>  # search community snippets
```

**Online features:**
- Browse and star community snippets
- Collections (curated snippet packs)
- Install a pack: `snape pack install go-starters`

---

## Tech Stack Summary

| Layer | Choice |
|---|---|
| Language | Go |
| CLI framework | Cobra |
| Database | SQLite (`modernc.org/sqlite` — pure Go, no CGo) |
| Web server | Chi router + `net/http` |
| Web UI | Server-rendered HTML + HTMX + highlight.js + CodeMirror |
| Config | TOML (`BurntSushi/toml`) |
| Testing | Go's built-in `testing` package |
| Distribution | Single binary via `goreleaser` |

---

## File Structure

```
snape/
├── cmd/              # Cobra commands (one file per command)
├── internal/
│   ├── db/           # SQLite queries and migrations
│   ├── snippet/      # Core snippet logic
│   ├── template/     # Variable substitution engine
│   ├── server/       # Web server and handlers
│   └── config/       # Config loading
├── web/
│   ├── templates/    # HTML templates
│   └── static/       # CSS, JS assets
├── main.go
└── snape.db          # Lives in ~/.snape/
```

---

## Build Order

1. DB schema + migrations
2. `add`, `get`, `list`, `delete` commands
3. Templating engine
4. `insert` command
5. `search` (SQLite FTS5)
6. Version history
7. `export` / `import`
8. Web UI (`serve`)
9. Shell completions
10. Community server

