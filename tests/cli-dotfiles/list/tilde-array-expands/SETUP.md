# Scenario

**Feature**: `list` expands `"~": [names]` to per-name mapping paths, not bare home

```
# operator
includeDotFiles: false
files: { "~": [".bashrc"] }
HOME has .bashrc (and optional Library — list must not imply full-home sole job)
operator -> bak-files list --config …
  -> stdout contains exact line HOME/alice/.bashrc
  -> stdout must NOT contain exact line HOME/alice  (bare ~ full-home mapping)
  -> exit 0
```

## Preconditions

- `tildeArrayConfigDotsOff(['.bashrc'])` so list reflects only expanded array jobs.
- Same discovery/expand rules as backup ResolveJobs.

## Steps

1. setupDotsWorld list + tildeArrayConfigDotsOff.
2. Write home/.bashrc.
3. ExpectedMappingPath = `HOME/alice/.bashrc`.

## Context

- Bug signature today: sole mapping path `HOME/alice` from key `"~"`.
- Fixed behavior: `HOME/alice/.bashrc` from expanded `~/`.bashrc`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, tildeArrayConfigDotsOff(`[".bashrc"]`), "list", false)
	writeFile(t, filepath.Join(home, ".bashrc"), "export X=1\n")
	writeFile(t, filepath.Join(home, "Library", "x"), "ignore\n")
	req.ExpectedMappingPath = "HOME/alice/.bashrc"
	return nil
}
```
