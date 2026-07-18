## Expected

- No harness error.
- `DefaultSkipMask.String()` equals `tmp|cache|logs|binary|trash|meta|vendor` (no `history`).
- `BackupSkipMask.String()` equals that string with `|history` appended (history last).
- `BackupSkipMask` has FlagHistory; `DefaultSkipMask` does not.

## Side Effects

- None.

## Errors

- history in DefaultSkipMask, missing from BackupSkipMask, or wrong order fails.

## Exit Code

- N/A

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	wantDefault := "tmp|cache|logs|binary|trash|meta|vendor"
	if resp.Flags != wantDefault {
		t.Fatalf("DefaultSkipMask.String(): got %q, want %q", resp.Flags, wantDefault)
	}
	if strings.Contains(resp.Flags, "history") {
		t.Fatalf("DefaultSkipMask must not include history: %q", resp.Flags)
	}
	if pathflag.DefaultSkipMask.Has(pathflag.FlagHistory) {
		t.Fatalf("DefaultSkipMask.Has(FlagHistory) want false")
	}

	bak := pathflag.BackupSkipMask
	gotBak := bak.String()
	wantBak := wantDefault + "|history"
	if gotBak != wantBak {
		t.Fatalf("BackupSkipMask.String(): got %q, want %q", gotBak, wantBak)
	}
	if !bak.Has(pathflag.FlagHistory) {
		t.Fatalf("BackupSkipMask.Has(FlagHistory) want true")
	}
	if bak != pathflag.DefaultSkipMask|pathflag.FlagHistory {
		t.Fatalf("BackupSkipMask want DefaultSkipMask|FlagHistory")
	}
}
```
