## Expected

- No harness error.
- `Response.Flags` is exactly `""`.

## Side Effects

- None.

## Errors

- Non-empty string for zero flag fails.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.Flags != "" {
		t.Fatalf("Flag(0).String(): got %q, want \"\"", resp.Flags)
	}
}
```
