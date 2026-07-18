## Expected

- No error.
- `Rule` is `.grok/logs` (path catalog, not `**/*.log`).
- `Flags` is `logs`.
- `Owner` is `grok`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule (e.g. only `**/*.log`) or flags fails.

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
	if resp.Rule != ".grok/logs" {
		t.Fatalf("Rule: got %q, want .grok/logs", resp.Rule)
	}
	if resp.Flags != "logs" {
		t.Fatalf("Flags: got %q, want logs", resp.Flags)
	}
	if resp.Owner != "grok" {
		t.Fatalf("Owner: got %q, want grok", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
