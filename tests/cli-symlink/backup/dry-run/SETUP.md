# Scenario

**Feature**: backup `--dry-run` logs symlink intent without writing targetDir

```
# dry-run
operator -> bak-files backup --config … --dry-run
  -> would symlink…; targetDir unchanged
```

## Preconditions

- Leaves set dry-run Args and TargetDirBefore fingerprint.

## Steps

1. Ensure `--dry-run` is present on Args when partially set by parents.
2. Leaves replace Args with full backup + `--dry-run` and fixtures.

## Context

- Zero file writes under targetDir is mandatory.
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
