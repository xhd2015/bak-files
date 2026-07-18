## Expected

- No error.
- `Rule` is `.cache`.
- `Flags` is `cache`.
- `Owner` is empty.

## Side Effects

- None.

## Errors

- Error treating spaces as invalid path fails if trim is missing.

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
	if resp.Rule != ".cache" {
		t.Fatalf("Rule: got %q, want .cache", resp.Rule)
	}
	if resp.Flags != "cache" {
		t.Fatalf("Flags: got %q, want cache", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
}
```
