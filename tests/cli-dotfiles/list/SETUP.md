# Scenario

**Feature**: `list` uses the same home-dot discovery as backup

```
# list discovers mapping paths for auto-dots
operator -> bak-files list --config …
  -> stdout: mapping paths for explicit + discovered jobs
```

## Preconditions

- Leaves supply WorkDir, empty-files config, HOME with top-level dots.

## Steps

1. Leaves write home trees and set Args to `list --config …`.

## Context

- List does not write targetDir; only stdout mapping paths matter.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"list"}
	}
	return nil
}
```
