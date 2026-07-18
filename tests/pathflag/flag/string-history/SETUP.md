# Scenario

**Feature**: FlagHistory.String is `history` (last bit name)

```
FlagHistory.String() -> "history"
```

## Preconditions

- `FlagHistory` exists as a named bit after vendor in ascending order.

## Steps

1. Set FlagBits to FlagHistory only.
2. Expect String `history`.

## Context

- Sessions and similar keep-local-but-skip-bak attributes use this name.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = OpFlagString
	req.FlagBits = uint64(pathflag.FlagHistory)
	return nil
}
```
