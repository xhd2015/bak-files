# Scenario

**Feature**: safe repo-root .gitignore for public bak-files module

```
# reviewer inspects module-root .gitignore
reviewer -> read .gitignore
  -> patterns cover binaries, payload dirs, stats/index, secrets, OS junk
```

## Preconditions

- File path: module root `.gitignore` (implementer creates it).

## Steps

1. Child sets `Mode=gitignore`.
2. `Run` reads file content into `Response.GitignoreContent`.

## Context

- Patterns may be written as path segments (`files/`) or root-anchored (`/files/`).
- Binary ignore may be `/bin/`, `bin/`, or the binary name `bak-files`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "gitignore"
	return nil
}
```
