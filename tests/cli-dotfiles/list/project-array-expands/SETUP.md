# Scenario

**Feature**: `list` expands `"$W0/proj": [".vscode"]` to child mapping, not bare PREFIX

```
# operator
includeDotFiles: false
files: { "$W0/proj": [".vscode"] }
mapping: { "$W0": "W" }
W0 tree has proj/.vscode/settings.json and proj/vendor/poison.go
operator -> bak-files list --config …
  -> stdout contains exact line W/proj/.vscode
  -> stdout must NOT contain exact line W/proj  (bare PREFIX full-tree job)
  -> exit 0
```

## Preconditions

- Same ResolveJobs rules as backup; list uses `engine.MappingPaths` → jobs.
- Dots off so only the expanded array child job is scheduled from this PREFIX.

## Steps

1. setupProjectArrayWorld list + projectArrayConfig.
2. Plant keep + poison under W0/proj (list does not write store; tree optional but realistic).
3. ExpectedMappingPath = `W/proj/.vscode`.

## Context

- Bug signature today: sole mapping path `W/proj` from unexpanded array value.
- Fixed: `W/proj/.vscode` from expanded `$W0/proj/.vscode` child key.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, w0, _ := setupProjectArrayWorld(t, req, projectArrayConfig(`[".vscode"]`), "list", false)
	writeFile(t, filepath.Join(w0, "proj", ".vscode", "settings.json"), "{}\n")
	writeFile(t, filepath.Join(w0, "proj", "vendor", "poison.go"), "vendor-poison\n")
	req.ExpectedMappingPath = "W/proj/.vscode"
	return nil
}
```
