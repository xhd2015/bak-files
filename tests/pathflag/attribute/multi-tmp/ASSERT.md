## Expected

- No error.
- `Rule` is `.tmp`.
- `Flags` is `tmp|cache` (ascending bit order).
- `Owner` is empty.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong flag order, missing bits, or non-empty Owner fails.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Err != "" {
		t.Fatalf("Classify(%q): unexpected error: %s", req.RelPath, resp.Err)
	}
	if resp.Rule != ".tmp" {
		t.Fatalf("Rule: got %q, want .tmp", resp.Rule)
	}
	if resp.Flags != "tmp|cache" {
		t.Fatalf("Flags: got %q, want tmp|cache", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
