# Scenario

**Feature**: .gitignore covers required safety patterns

```
# reviewer checks .gitignore for P1 safety patterns
reviewer -> read module .gitignore
  -> includes files/, bak.stats, sum.index, .env, .DS_Store, and binary/bin ignore
```

## Preconditions

- `Mode=gitignore`.

## Steps

1. Rely on Run to load `.gitignore` text.

## Context

- Assert by content scan (not `git check-ignore`), so the leaf works without
  requiring a full index of ignored paths.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "gitignore"
	return nil
}
```
