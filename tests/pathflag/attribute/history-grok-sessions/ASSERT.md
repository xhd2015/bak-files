## Expected

- No error.
- `Rule` is `.grok/sessions`.
- `Flags` is exactly `history` (not `cache` or `tmp`).
- `Owner` is `grok`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule, cache/tmp flags, or missing history fails.

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
	if resp.Rule != ".grok/sessions" {
		t.Fatalf("Rule: got %q, want .grok/sessions", resp.Rule)
	}
	if resp.Flags != "history" {
		t.Fatalf("Flags: got %q, want history", resp.Flags)
	}
	if strings.Contains(resp.Flags, "cache") || strings.Contains(resp.Flags, "tmp") {
		t.Fatalf("Flags: sessions must not be cache/tmp, got %q", resp.Flags)
	}
	if resp.Owner != "grok" {
		t.Fatalf("Owner: got %q, want grok", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
