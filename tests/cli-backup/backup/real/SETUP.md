# Scenario

**Feature**: real backup writes under `targetDir` (no `--dry-run`)

```
# real backup may create/overwrite under targetDir
operator -> bak-files backup --config …
  -> files under WorkDir/files/…
```

## Preconditions

- Args must not include `--dry-run`.

## Steps

1. Leaves create sources and config; Assert checks BackupPath content.

## Context

- Sibling of `backup/dry-run/`.

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

