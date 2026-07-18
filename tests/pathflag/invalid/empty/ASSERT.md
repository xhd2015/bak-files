## Expected

- `Run` harness error is nil (error is in the response).
- `Response.Err` is non-empty (Classify failed).
- Rule / Reason / Flags / Owner remain empty.

## Side Effects

- None (pure function).

## Errors

- Missing error (successful zero Result) fails the test.

## Exit Code

- N/A (library call)

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Err == "" {
		t.Fatalf("Classify(%q): want error, got Rule=%q Flags=%q Owner=%q",
			req.RelPath, resp.Rule, resp.Flags, resp.Owner)
	}
}
```
