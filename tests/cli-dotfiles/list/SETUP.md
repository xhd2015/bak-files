# Scenario

**Feature**: `list` uses the same ResolveJobs as backup (dots + `"~"` / PREFIX array expand)

```
# list discovers mapping paths for auto-dots and expanded array names
operator -> bak-files list --config …
  -> stdout: mapping paths for explicit + discovered + "~":[name] + PREFIX:[name] jobs
  -> bare "~" must not appear as sole HOME/$ROLE line without name expand
  -> bare PREFIX must not appear alone when value is a string array whitelist
```

## Preconditions

- Leaves supply WorkDir, config (empty files, `"~"` array, or `"$W0/…"` array), fixtures.

## Steps

1. Leaves write trees and set Args to `list --config …`.

## Context

- List does not write targetDir; only stdout mapping paths matter.
- `list/project-array-expands` is RED until ResolveJobs generalizes array expand.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"list"}
	}
	return nil
}
```
