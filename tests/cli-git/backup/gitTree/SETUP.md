# Scenario

**Feature**: files entry with `gitTree: true` (working-tree-aware directory backup)

```
# gitTree
config files["~/repo"] = { "gitTree": true }
operator -> bak-files backup …
  -> dirty: copy + bak.stats
  -> clean: skip INFO
  -> dry-run: no writes
```

## Preconditions

- Entry mode is **gitTree** (not gitPatch, not plain copy).
- Source `~/repo` is a local git directory.

## Steps

1. Descendants set dirty/clean and real/dry-run.

## Context

- TS parity: skip when `!hasChange`; write stats then copy when dirty.
- No remotes required.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
