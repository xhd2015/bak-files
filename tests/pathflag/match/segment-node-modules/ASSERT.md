## Expected

- No error.
- `Rule` is `**/node_modules`.
- `Flags` is `vendor`.
- `Owner` is empty.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Missing vendor flag or wrong rule fails.

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
	if resp.Rule != "**/node_modules" {
		t.Fatalf("Rule: got %q, want **/node_modules", resp.Rule)
	}
	if resp.Flags != "vendor" {
		t.Fatalf("Flags: got %q, want vendor", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
