# Scenario

**Feature**: `.codex` prefix owns nested non-catalog paths

```
caller -> Classify(".codex/config.toml")
  -> Owner=codex, Flags empty, Rule empty
```

## Preconditions

- Catalog rules under codex include `.codex/.tmp`, `.codex/skills/.system`,
  `.codex/sessions` (history), and `.codex/logs_2.sqlite` (logs) — not
  arbitrary config files like `config.toml`.

## Steps

1. Set path to `.codex/config.toml` (uncatalogued under `.codex`).
2. Expect owner-only Result.

## Context

- Distinguishes owner attribution from history/cache/tmp/logs attribute rows.
- Session trees are **history**, not owner-only (see `attribute/history-codex-sessions`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".codex/config.toml"
	return nil
}
```
