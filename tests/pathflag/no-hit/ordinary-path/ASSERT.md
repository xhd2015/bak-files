## Expected

- `Response.Err` is empty.
- `Rule`, `Reason`, `Flags`, and `Owner` are all empty strings.

## Side Effects

- None.

## Errors

- Any non-empty Result field or Classify error fails.

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
