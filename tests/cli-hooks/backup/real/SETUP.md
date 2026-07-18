# Scenario

**Feature**: real backup may execute hooks/`cmd` and write under `targetDir`

```
# real backup executes shell and may create files
operator -> bak-files backup --config …
  -> beforeCopy / cmd / afterCopy run; targetDir receives output
```

## Preconditions

- Args must not include `--dry-run`.

## Steps

1. Leaves build configs with `cmd` or `beforeCopy`; Assert checks side effects.

## Context

- Sibling of `backup/dry-run/`.

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Real branch: strip accidental --dry-run if present (leaves overwrite Args).
	filtered := req.Args[:0]
	for _, a := range req.Args {
		if a == "--dry-run" {
			continue
		}
		filtered = append(filtered, a)
	}
	req.Args = filtered
	_ = strings.Join(req.Args, " ")
	return nil
}
```
