# Scenario

**Feature**: list shows discovered mapping paths when dots are on

```
# operator
HOME=…/home  WORKING_ROLE=alice
write home/.bashrc = "export X=1\n"
config: files={}, mapping ~ → HOME/$WORKING_ROLE
operator -> bak-files list --config …
  -> stdout contains HOME/alice/.bashrc (exit 0)
```

## Preconditions

- Empty `files`; default includeDotFiles; home has only `.bashrc`.

## Steps

1. setupDotsWorld with emptyFilesConfig, command list.
2. Write home/.bashrc.
3. ExpectedMappingPath = `HOME/alice/.bashrc`.

## Context

- Same discovery as backup; list is the observability surface.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "list", false)
	writeFile(t, filepath.Join(home, ".bashrc"), "export X=1\n")
	req.ExpectedMappingPath = "HOME/alice/.bashrc"
	return nil
}
```
