# bak-files

Config-driven backup and restore of selected files under `$HOME` and other whitelisted paths. Copies into a local `targetDir` tree using path mapping (e.g. `HOME/$WORKING_ROLE`), with automatic home-dot discovery and a catalog of paths to skip (caches, sessions, toolchains).

## Install

```bash
go install github.com/xhd2015/bak-files/cmd/bak-files@latest
# or from this repo:
go build -o bak-files ./cmd/bak-files
```

Requires Go (see `go.mod`).

## Quick start

```bash
export HOME WORKING_ROLE   # builtins; add any env listed in config validate
# optional: W0, etc.

bak-files list --config bak.config.json
bak-files backup --config bak.config.json --dry-run
bak-files backup --config bak.config.json
bak-files restore --config bak.config.json --dry-run
```

Default config path: `bak.config.json` in the current working directory.

## Commands

| Command | Purpose |
|---------|---------|
| `backup` | Copy sources → `targetDir` per config |
| `restore` | Copy `targetDir` → sources |
| `list` | Print resolved mapping paths for jobs |

Common flags (backup/restore):

| Flag | Meaning |
|------|---------|
| `--config FILE` | Config JSON (default `bak.config.json`) |
| `--dry-run` | Log actions without writing |
| `--no-dot-files` | Disable auto-discovery of top-level `$HOME` dots |
| `--include PATH` | Force-include home-relative path (repeatable) |
| `--exclude PATH` | Force-exclude home-relative path (repeatable) |
| `-v, --verbose` | More detail |

## Configuration model

Backup is **whitelist-driven**: only paths that match `files` entries (plus optional default home dots) are included.

### Minimal example

```json
{
  "validate": [{ "env": ["HOME", "WORKING_ROLE"] }],
  "files": {
    "~": [".bash_profile", ".bashrc", ".ssh"],
    "~/Scripts": {
      "excludes": ["node_modules", ".git"]
    },
    "~/.config": {
      "excludes": [".DS_Store"]
    }
  },
  "targetDir": "./files",
  "global": {
    "excludes": [".DS_Store", "node_modules", ".git"],
    "dotExcludes": [".local", ".knowledge-hub"],
    "includeDotFiles": true
  },
  "mapping": {
    "~": "HOME/$WORKING_ROLE",
    "~/Scripts": "HOME/Scripts"
  }
}
```

### `files` entry shapes

| Value | Behavior |
|-------|----------|
| `true` | Include the expanded key path (full tree) |
| `{ "excludes": [...], "file": "...", "gitTree": true, hooks… }` | Full tree (or `file` source) with options |
| `[ "a", "b" ]` | **Whitelist basenames under the key prefix** — one job per `prefix/a`, `prefix/b`; never the full prefix alone |

Examples:

```json
"~": [".bashrc", ".ssh"]
"$W0/credit-pricing-center": [".vscode"]
```

→ only `$HOME/.bashrc`, `$HOME/.ssh`, and `$W0/credit-pricing-center/.vscode`.

Bare `"~"` is **not** a full-home recursive backup. Non-array values under bare `"~"` do not schedule a whole-`$HOME` job.

### Home dots (`includeDotFiles`)

When enabled (**default true** if omitted), bak-files treats discovery as a synthetic whitelist of **top-level** `$HOME` names matching `~/.*` (dotfiles and dotdirs). Disable with:

- `"global": { "includeDotFiles": false }`, or  
- `--no-dot-files`

Discovered dots still respect `dotExcludes`, `--exclude` / `--include`, and **pathflag** skips.

### Mapping and layout

- `mapping` maps source path prefixes (after env/`~` expand) into relative paths under `targetDir`.
- Typical layout: `files/HOME/<WORKING_ROLE>/…`, `files/HOME/Scripts/…`.

### Global filters

| Field | Role |
|-------|------|
| `global.excludes` | Basename/glob skip during walks (e.g. `.DS_Store`, `node_modules`) |
| `global.dotExcludes` | Force-exclude home-relative prefixes (CLI `--exclude` family) |
| `global.dotIncludes` | Force-include home-relative prefixes |
| `global.gitExcludes` | Patterns for git-aware modes (when used) |

Walk skip order under `$HOME` (simplified): force-exclude → pathflag backup-skip catalog → basename excludes. Force-include can keep pathflag-skipped paths when both apply (exclude still wins if both match).

## Path catalog (`pathflag`)

Package `pathflag` classifies home-relative paths (no I/O). Flags in **`BackupSkipMask`** are skipped on backup/restore walks (caches, tmp, logs, binaries, trash, meta, vendor, **history**).

Examples:

| Path | Typical flags |
|------|----------------|
| `.cache`, `.npm`, `.tmp`, `.sandbox`, `.wrk` | cache / tmp |
| `.grok/sessions`, `.codex/sessions` | **history** (user-local; skip bak, not “safe delete”) |
| `.grok/projects`, `.claude/projects` | cache |
| `.nvm` | cache\|binary |
| `.local/share/containers` | cache |
| `.Trash` | trash |

`DefaultSkipMask` is reclaim-oriented (**without** history). Bak uses **`BackupSkipMask`** = default + history.

Owners (`grok`, `codex`, `claude`, `android`, …) are metadata for the catalog; skip is driven by flags.

## Environment

Always required (even if omitted from config): `HOME`, `WORKING_ROLE`.

Config `validate[].env` can require more (e.g. `W0`).

## Development

```bash
go build -o bak-files ./cmd/bak-files
go test ./pathflag/... ./internal/...

# doctests (requires doctest CLI)
doctest vet ./tests/pathflag
doctest test ./tests/pathflag
doctest test ./tests/cli-dotfiles
doctest test ./tests/cli-backup
```

### Layout

```text
cmd/bak-files/     CLI entrypoint
internal/config/   bak.config load / mapping / env
internal/engine/   resolve jobs, copy/restore, skip policy
pathflag/          pure home-relative path classifier
tests/             doctest trees (cli-*, pathflag)
```

## Notes

- Prefer **explicit** `files` keys and array whitelists; rely on default dots only when you want all top-level home dots.
- After changing config or pathflag rules, rebuild the binary and clear a stale `targetDir` if old full-tree copies remain.
- Private/local configs (outside this repo) are not required to run the tool.

## License

See repository license if present; otherwise all rights reserved by the author unless stated.
