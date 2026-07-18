## Expected

- No Classify error.
- `Rule`, `Reason`, `Flags`, and `Owner` are all empty (zero Result / no skip attributes).

## Side Effects

- None.

## Errors

- Any attribute hit (especially binary via v2ray-core) fails.

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
		t.Fatalf("Classify(%q): want zero Result (config not under v2ray-core rule), got Rule=%q Reason=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Reason, resp.Flags, resp.Owner)
	}
}
```
