## Expected

- No harness error.
- `Response.Flags` is exactly `tmp|cache|logs`.

## Side Effects

- None.

## Errors

- Wrong separators, case, or order fails.

## Exit Code

- N/A

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	want := "tmp|cache|logs"
	if resp.Flags != want {
		t.Fatalf("Flag multi.String(): got %q, want %q", resp.Flags, want)
	}
}
```
