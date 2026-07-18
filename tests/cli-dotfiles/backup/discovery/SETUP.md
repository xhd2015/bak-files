# Scenario

**Feature**: auto-discover top-level home dots when includeDotFiles is on

```
# discovery branch
HOME has top-level .names; files may be empty or sparse
operator -> bak-files backup …
  -> jobs for uncovered ~/name  |  no jobs when dots disabled
```

## Preconditions

- Prefer empty `files` so only auto-discovery creates jobs (except dedupe leaf).

## Steps

1. Leaves create home dots and set enablement (default / config / flag).

## Context

- Sibling of `filters/` and `walk-skip/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Document discovery branch; leaves fully set Args/fixtures.
	if req.HomeDir == "" && req.WorkDir == "" {
		// Leaf will set world; no-op marker for grouping.
		_ = req.Args
	}
	return nil
}
```
