# Scenario

**Feature**: Flag.Names returns names in ascending bit order

```
(FlagVendor|FlagTmp|FlagBinary).Names()
  -> []string{"tmp", "binary", "vendor"}
```

## Preconditions

- Names slice order matches String split order.

## Steps

1. Combine Tmp|Binary|Vendor (non-adjacent bits).
2. Expect Names: tmp, binary, vendor.

## Context

- Complements String(); Classify Response uses String(), Names is for callers needing a list.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = OpFlagNames
	req.FlagBits = uint64(pathflag.FlagTmp | pathflag.FlagBinary | pathflag.FlagVendor)
	return nil
}
```
