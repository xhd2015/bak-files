## Expected

- No harness error.
- `Response.Flags` is exactly `history`.

## Side Effects

- None.

## Errors

- Missing constant, wrong name, or empty string fails.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Flags != "history" {
		t.Fatalf("FlagHistory.String(): got %q, want history", resp.Flags)
	}
}
```
