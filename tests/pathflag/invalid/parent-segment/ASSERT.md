## Expected

- `Response.Err` is non-empty when the path contains a `..` segment.

## Side Effects

- None.

## Errors

- Successful Classify fails the test.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Err == "" {
		t.Fatalf("Classify(%q): want error for parent segment, got Rule=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Flags, resp.Owner)
	}
}
```
