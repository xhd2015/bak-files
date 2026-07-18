# Scenario

**Feature**: list help paths do not require config or env

```
# operator asks for list help
operator -> bak-files list --help
  -> Usage on stdout (exit 0); no bak.config needed
```

## Preconditions

- No config file, no HOME/WORKING_ROLE required for help.

## Steps

1. Grouping marks this branch as help-only.

## Context

- Regression: help must stay green while list body gains config load.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Help must not depend on HOME/WORKING_ROLE; children may refine Args.
	req.Env = minimalEnv()
	if len(req.Args) == 0 {
		req.Args = []string{"list", "--help"}
	}
	return nil
}
```
