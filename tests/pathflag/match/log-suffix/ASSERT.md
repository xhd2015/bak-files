## Expected

- No error.
- `Rule` is `**/*.log`.
- `Flags` is `logs`.
- `Owner` is empty.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Zero Result or wrong rule fails.

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
	if resp.Rule != "**/*.log" {
		t.Fatalf("Rule: got %q, want **/*.log", resp.Rule)
	}
	if resp.Flags != "logs" {
		t.Fatalf("Flags: got %q, want logs", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
