## Expected

- No error.
- `Rule` is `.codex/logs_2.sqlite`.
- `Flags` is `logs`.
- `Owner` is `codex`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule or flags fails.

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
	if resp.Rule != ".codex/logs_2.sqlite" {
		t.Fatalf("Rule: got %q, want .codex/logs_2.sqlite", resp.Rule)
	}
	if resp.Flags != "logs" {
		t.Fatalf("Flags: got %q, want logs", resp.Flags)
	}
	if resp.Owner != "codex" {
		t.Fatalf("Owner: got %q, want codex", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
