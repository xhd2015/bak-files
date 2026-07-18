# Scenario

**Feature**: DefaultSkipMask excludes history; BackupSkipMask includes it

```
DefaultSkipMask.String()  // reclaim-oriented: no history
BackupSkipMask.String()   // bak skip: DefaultSkipMask | FlagHistory
```

## Preconditions

- `DefaultSkipMask` = tmp|cache|logs|binary|trash|meta|vendor (no history).
- `BackupSkipMask` (or equivalent) = DefaultSkipMask | FlagHistory.

## Steps

1. This leaf exercises both masks via Flag.String (two assertions on constants).
2. Expect history only in BackupSkipMask Names/String.

## Context

- History is skip-for-bak but not reclaim-as-cache/tmp junk.

```go
import (
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Setup(t *testing.T, req *Request) error {
	// Drive Flag.String on DefaultSkipMask; Assert also checks BackupSkipMask.
	req.Op = OpFlagString
	req.FlagBits = uint64(pathflag.DefaultSkipMask)
	return nil
}
```
