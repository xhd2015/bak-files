## Expected

- No error.
- `Rule` is `.local/share/opencode/log` (not a shorter prefix and not only `**/*.log`).
- `Flags` is `logs`.
- `Owner` is `opencode`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Owner-only zero flags or wrong rule fails.

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
	if resp.Rule != ".local/share/opencode/log" {
		t.Fatalf("Rule: got %q, want .local/share/opencode/log", resp.Rule)
	}
	if resp.Flags != "logs" {
		t.Fatalf("Flags: got %q, want logs", resp.Flags)
	}
	if resp.Owner != "opencode" {
		t.Fatalf("Owner: got %q, want opencode", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
