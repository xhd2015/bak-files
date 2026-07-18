# Scenario

**Feature**: multi-bit Flag.String uses `|` in ascending bit order

```
(FlagTmp|FlagCache|FlagLogs).String() -> "tmp|cache|logs"
```

## Preconditions

- Bits set via pathflag constants (not hard-coded numeric literals alone).

## Steps

1. Combine Tmp|Cache|Logs into FlagBits.
2. Expect exact string `tmp|cache|logs`.

## Context

- Order is bit order, not declaration convenience of the caller.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = OpFlagString
	req.FlagBits = uint64(pathflag.FlagTmp | pathflag.FlagCache | pathflag.FlagLogs)
	return nil
}
```
