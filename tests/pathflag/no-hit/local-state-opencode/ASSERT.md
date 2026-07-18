## Expected

- No Classify error.
- `Rule`, `Reason`, `Flags` empty (no attribute skip).
- `Owner` empty (no owner prefix for this path).

## Side Effects

- None.

## Errors

- Any catalog skip flags or whole-`.local` match fails.

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
	if resp.Rule != "" || resp.Reason != "" || resp.Flags != "" || resp.Owner != "" {
		t.Fatalf("Classify(%q): want zero Result, got Rule=%q Reason=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Reason, resp.Flags, resp.Owner)
	}
}
```
