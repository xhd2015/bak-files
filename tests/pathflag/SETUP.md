# Scenario

**Feature**: pure pathflag.Classify and Flag helpers for home-relative paths

```
# caller supplies a home-relative path; library attributes without I/O
caller -> pathflag.Classify(relPath)
  -> Result{Rule, Reason, Flags, Owner} | error

# Flag helpers (assert formatting)
pathflag.Flag bits -> String() "tmp|cache|..." | Names() []string
```

## Preconditions

- Module root: `DOCTEST_ROOT/../..` (repo with `go.mod`).
- Production package (implementer): `github.com/xhd2015/bak-files/pathflag`.
- No filesystem fixtures required — Classify is pure string logic.
- Classic TDD: package may be missing; compile/runtime RED is expected until implementer lands code.
- Default `Request.Op` is `classify`. Flag leaves set `flag_string` / `flag_names`.

## Steps

1. Root Setup defaults `Op` to classify when unset.
2. Grouping / leaf Setup narrows `RelPath` or Flag helper fields.
3. Root `Run` calls real `pathflag.Classify` or `Flag.String` / `Flag.Names`.
4. Leaf Assert checks Rule, Reason, Flags string, Owner, error, or Names.

## Context

- Flag.String names (ascending bit order): tmp, cache, logs, binary, trash, meta, vendor.
- Owner strings: codex, opencode, grok, cursor, bun, npm, cargo, chromium (or empty).
- Longest catalog prefix wins among path rules; owner table is always resolved.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if req.Op == "" {
		req.Op = OpClassify
	}
	return nil
}
```
