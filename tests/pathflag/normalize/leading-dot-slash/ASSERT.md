## Expected

- No error.
- `Rule` is `.bun`.
- `Flags` is `cache`.
- `Owner` is `bun`.

## Side Effects

- None.

## Errors

- Error or miss as if path were unknown fails.

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
	if resp.Rule != ".bun" {
		t.Fatalf("Rule: got %q, want .bun", resp.Rule)
	}
	if resp.Flags != "cache" {
		t.Fatalf("Flags: got %q, want cache", resp.Flags)
	}
	if resp.Owner != "bun" {
		t.Fatalf("Owner: got %q, want bun", resp.Owner)
	}
}
```
