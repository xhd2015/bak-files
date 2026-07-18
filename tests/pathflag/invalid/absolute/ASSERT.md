## Expected

- `Response.Err` is non-empty.
- No successful attribute Result.

## Side Effects

- None.

## Errors

- Successful Classify on absolute path fails the test.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Err == "" {
		t.Fatalf("Classify(%q): want error for absolute path, got Rule=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Flags, resp.Owner)
	}
}
```
