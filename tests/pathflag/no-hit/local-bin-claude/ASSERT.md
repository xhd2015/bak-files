## Expected

- No Classify error.
- `Rule`, `Reason`, `Flags`, and `Owner` are all empty.

## Side Effects

- None.

## Errors

- Any hit from `.local/share/claude` (cache|binary) fails.

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
		t.Fatalf("Classify(%q): want zero Result (bin is not share/claude), got Rule=%q Reason=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Reason, resp.Flags, resp.Owner)
	}
}
```
