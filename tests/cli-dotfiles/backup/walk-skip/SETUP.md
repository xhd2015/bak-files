# Scenario

**Feature**: walk under a home dir applies pathflag and basename excludes

```
# nested walk
operator -> bak-files backup … (dots on; dir job under $HOME)
  -> pathflag DefaultSkipMask paths skipped
  -> global.excludes basename/glob still applied
```

## Preconditions

- Discovered or walkable top-level home directory with nested keep/skip files.

## Steps

1. Leaves plant nested trees under home dots / app dirs.

## Context

- pathflag runs on any path under $HOME during any job walk.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req.Args
	return nil
}
```
