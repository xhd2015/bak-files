# Scenario

**Feature**: real restore writes destinations from `targetDir`

```
# real restore
operator -> bak-files restore --config …
  -> SourcePath content from backup store
```

## Preconditions

- Args without `--dry-run`.

## Steps

1. Leaves seed BackupPath; dest may be absent or stale.

## Context

- Sibling of `restore/dry-run/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Real restore: strip accidental --dry-run so leaf Args define write mode.
	out := make([]string, 0, len(req.Args))
	for _, a := range req.Args {
		if a == "--dry-run" {
			continue
		}
		out = append(out, a)
	}
	req.Args = out
	return nil
}
```

