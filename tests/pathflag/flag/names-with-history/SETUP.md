# Scenario

**Feature**: history appears last in Names ascending bit order

```
(FlagVendor|FlagHistory|FlagTmp).Names()
  -> []string{"tmp", "vendor", "history"}
```

## Preconditions

- Bit order: tmp … vendor, then history last.

## Steps

1. Combine Tmp|Vendor|History.
2. Expect Names with history after vendor.

## Context

- Extends names-order to the new FlagHistory bit.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = OpFlagNames
	req.FlagBits = uint64(pathflag.FlagTmp | pathflag.FlagVendor | pathflag.FlagHistory)
	return nil
}
```
