# Scenario

**Feature**: backup `--dry-run` never executes hooks or `cmd` and never writes

```
# dry-run must not spawn shell for hooks/cmd
operator -> bak-files backup --config … --dry-run
  -> log would… (prefer dry-run / would)
  -> marker absent; cmd-generated file absent
```

## Preconditions

- Args include `--dry-run`.

## Steps

1. Leaf configures both `cmd` and `beforeCopy` that would leave fingerprints;
   Assert proves none ran.

## Context

- Core P4 dry-run exit criterion for hooks/cmd.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Dry-run branch: ensure --dry-run will be present once leaf sets Args.
	// Leaves append the flag; record branch intent on empty Args placeholder.
	if len(req.Args) == 0 {
		req.Args = []string{"backup", "--dry-run"}
	}
	hasDry := false
	for _, a := range req.Args {
		if a == "--dry-run" {
			hasDry = true
			break
		}
	}
	if !hasDry && len(req.Args) > 0 && req.Args[0] == "backup" {
		req.Args = append(req.Args, "--dry-run")
	}
	return nil
}
```
