# Scenario

**Feature**: zero Flag formats as empty string

```
pathflag.FlagNone.String() -> ""
```

## Preconditions

- Flag value 0 / FlagNone.

## Steps

1. Set Op `flag_string` and FlagBits 0.
2. Expect `Flags` == `""`.

## Context

- Zero Result.Flags.String() is empty, not `"none"`.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = OpFlagString
	req.FlagBits = uint64(pathflag.FlagNone)
	return nil
}
```
