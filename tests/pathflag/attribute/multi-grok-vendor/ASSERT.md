## Expected

- No error.
- `Rule` is `.grok/vendor`.
- `Flags` is `cache|vendor`.
- `Owner` is `grok`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong flags/owner fails.

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
	if resp.Rule != ".grok/vendor" {
		t.Fatalf("Rule: got %q, want .grok/vendor", resp.Rule)
	}
	if resp.Flags != "cache|vendor" {
		t.Fatalf("Flags: got %q, want cache|vendor", resp.Flags)
	}
	if resp.Owner != "grok" {
		t.Fatalf("Owner: got %q, want grok", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
