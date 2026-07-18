## Expected

- No error.
- `Rule` is `.bun`.
- `Flags` is `cache`.
- `Owner` is `bun`.
- `Reason` mentions Bun / install / cache (non-empty; contains `cache` or `Bun`).

## Side Effects

- None.

## Errors

- Wrong rule/flags/owner fails.

## Exit Code

- N/A

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Err != "" {
		t.Fatalf("Classify(%q): unexpected error: %s", req.RelPath, resp.Err)
	}
	if resp.Rule != ".bun" {
		t.Fatalf("Rule: got %q, want .bun", resp.Rule)
	}
	if resp.Flags != "cache" {
		t.Fatalf("Flags: got %q, want cache", resp.Flags)
	}
	if resp.Owner != "bun" {
		t.Fatalf("Owner: got %q, want bun", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
	low := strings.ToLower(resp.Reason)
	if !strings.Contains(low, "bun") && !strings.Contains(low, "cache") {
		t.Fatalf("Reason %q: want mention of bun or cache", resp.Reason)
	}
}
```
