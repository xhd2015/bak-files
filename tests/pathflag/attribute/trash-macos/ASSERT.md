## Expected

- No error.
- `Rule` is `.Trash`.
- `Flags` is `trash`.
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
	if resp.Rule != ".Trash" {
		t.Fatalf("Rule: got %q, want .Trash", resp.Rule)
	}
	if resp.Flags != "trash" {
		t.Fatalf("Flags: got %q, want trash", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
