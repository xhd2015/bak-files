# Scenario

**Feature**: `backup --dry-run` plans copies without writing `targetDir`

```
# dry-run backup
operator -> bak-files backup --config … --dry-run
  -> log intent (prefer "dry-run" or "would")
  -> zero writes under targetDir
  -> exit 0
```

## Preconditions

- Source exists; targetDir starts absent or empty.

## Steps

1. Leaves set `--dry-run` and capture TargetDirBefore fingerprint.

## Context

- Sibling of `backup/real/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Dry-run branch marker: leaves append --dry-run; ensure flag token is known.
	hasDry := false
	for _, a := range req.Args {
		if a == "--dry-run" {
			hasDry = true
			break
		}
	}
	if !hasDry && len(req.Args) > 0 {
		// Leaf Setup runs after this and sets full Args including --dry-run.
		// Record intent so this Setup is not a no-op stub.
		req.Args = append(append([]string{}, req.Args...), "--dry-run")
	}
	return nil
}
```

