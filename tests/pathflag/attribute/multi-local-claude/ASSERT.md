## Expected

- No error.
- `Rule` is `.local/share/claude`.
- `Flags` is `cache|binary` (not `binary|cache`).
- `Owner` is empty string (no OwnerClaude required).
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong order or missing bits fails.

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
	if resp.Rule != ".local/share/claude" {
		t.Fatalf("Rule: got %q, want .local/share/claude", resp.Rule)
	}
	if resp.Flags != "cache|binary" {
		t.Fatalf("Flags: got %q, want cache|binary", resp.Flags)
	}
	// Prefer empty owner unless implementer adds OwnerClaude; empty is locked default here.
	if resp.Owner != "" && resp.Owner != "claude" {
		t.Fatalf("Owner: got %q, want empty or claude", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
