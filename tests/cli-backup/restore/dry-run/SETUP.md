# Scenario

**Feature**: `restore --dry-run` plans writes without modifying destinations

```
# dry-run restore
operator -> bak-files restore --config … --dry-run
  -> no dest writes; prefer dry-run/would logs; exit 0
```

## Preconditions

- Backup store seeded; dest controlled by leaf.

## Steps

1. Leaves set `--dry-run` and DestBefore fingerprint/content.

## Context

- Sibling of `restore/real/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	hasDry := false
	for _, a := range req.Args {
		if a == "--dry-run" {
			hasDry = true
			break
		}
	}
	if !hasDry && len(req.Args) > 0 {
		req.Args = append(append([]string{}, req.Args...), "--dry-run")
	}
	return nil
}
```

