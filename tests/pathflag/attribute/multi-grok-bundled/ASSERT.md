## Expected

- No error.
- `Rule` is `.grok/bundled`.
- `Flags` is `cache|binary|vendor`.
- `Owner` is `grok`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong flag order or missing bits fails.

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
	if resp.Rule != ".grok/bundled" {
		t.Fatalf("Rule: got %q, want .grok/bundled", resp.Rule)
	}
	if resp.Flags != "cache|binary|vendor" {
		t.Fatalf("Flags: got %q, want cache|binary|vendor", resp.Flags)
	}
	if resp.Owner != "grok" {
		t.Fatalf("Owner: got %q, want grok", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
