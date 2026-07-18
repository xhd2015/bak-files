# Scenario

**Feature**: files entry with `gitPatch: true` writes a patch from git diff

```
# gitPatch
config files["~/repo"] = { "gitPatch": true }
operator -> bak-files backup …
  -> (dirty, real) .patch under mapping from git diff HEAD
  -> (--dry-run) no patch file
```

## Preconditions

- Entry mode is **gitPatch** (not gitTree).
- Local-only repo; MVP base is **HEAD** (no origin required).

## Steps

1. Descendants set real vs dry-run with dirty worktree.

## Context

- Requirement: prefer single `.patch` file content from `git diff HEAD`.
- Implementers may write `files/HOME/repo/worktree.patch` or any `*.patch`
  under `files/HOME/repo/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
