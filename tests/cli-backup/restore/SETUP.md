# Scenario

**Feature**: `bak-files restore` copies from `targetDir` back to the filesystem

```
# restore pipeline
operator -> bak-files restore --config … [--dry-run]
  -> destinations updated (real) | no dest writes (dry-run)
```

## Preconditions

- Leaves seed backup store and control whether dest exists.

## Steps

1. Grouping documents real vs dry-run branches.

## Context

- Same config mapping as backup tests (`~/notes.txt` ↔ `files/HOME/notes.txt`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"restore"}
	} else if req.Args[0] != "restore" && req.Args[0] != "backup" {
		req.Args = append([]string{"restore"}, req.Args...)
	}
	return nil
}
```

