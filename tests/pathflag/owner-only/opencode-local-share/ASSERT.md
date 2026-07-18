## Expected

- No error.
- `Owner` is `opencode`.
- `Rule`, `Reason`, and `Flags` are empty.

## Side Effects

- None.

## Errors

- Wrong owner or unexpected Rule/Flags fails.

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
	if resp.Owner != "opencode" {
		t.Fatalf("Owner: got %q, want opencode", resp.Owner)
	}
	if resp.Rule != "" || resp.Reason != "" || resp.Flags != "" {
		t.Fatalf("want owner-only, got Rule=%q Reason=%q Flags=%q",
			resp.Rule, resp.Reason, resp.Flags)
	}
}
```
