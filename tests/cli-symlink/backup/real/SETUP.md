# Scenario

**Feature**: real backup (no `--dry-run`) may write under targetDir

```
# real backup
operator -> bak-files backup --config …
  -> targetDir receives regular files + OS symlinks
```

## Preconditions

- Args must not include `--dry-run`.
- Leaves build Scripts fixtures and path fields on Request.

## Steps

1. Strip accidental `--dry-run` from Args at this branch.
2. Leaves overwrite Args with full backup invocation and fixtures.

## Context

- Real mode creates `symlinked …` logs; dry-run branch is sibling.

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Real branch: ensure Args will not carry --dry-run from a mistaken parent.
	filtered := req.Args[:0]
	for _, a := range req.Args {
		if a == "--dry-run" {
			continue
		}
		filtered = append(filtered, a)
	}
	req.Args = filtered
	// Leaves overwrite Args; this documents the real (non-dry-run) branch.
	_ = strings.Join(req.Args, " ")
	return nil
}
```
