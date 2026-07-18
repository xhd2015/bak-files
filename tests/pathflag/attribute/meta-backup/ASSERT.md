## Expected

- No error.
- `Rule` is `.backup`.
- `Flags` is `meta`.
- `Owner` is empty.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule/flags fails.

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
	if resp.Rule != ".backup" {
		t.Fatalf("Rule: got %q, want .backup", resp.Rule)
	}
	if resp.Flags != "meta" {
		t.Fatalf("Flags: got %q, want meta", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
