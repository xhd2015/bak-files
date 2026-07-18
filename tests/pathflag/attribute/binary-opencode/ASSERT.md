## Expected

- No error.
- `Rule` is `.opencode/bin`.
- `Flags` is `binary`.
- `Owner` is `opencode`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule/flags/owner fails.

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
	if resp.Rule != ".opencode/bin" {
		t.Fatalf("Rule: got %q, want .opencode/bin", resp.Rule)
	}
	if resp.Flags != "binary" {
		t.Fatalf("Flags: got %q, want binary", resp.Flags)
	}
	if resp.Owner != "opencode" {
		t.Fatalf("Owner: got %q, want opencode", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
