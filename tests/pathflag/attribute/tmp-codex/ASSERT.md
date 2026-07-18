## Expected

- No error.
- `Rule` is `.codex/.tmp`.
- `Flags` is `tmp|cache` (ascending bit order).
- `Owner` is `codex`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong flag order or missing bits fails.

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
	if resp.Rule != ".codex/.tmp" {
		t.Fatalf("Rule: got %q, want .codex/.tmp", resp.Rule)
	}
	if resp.Flags != "tmp|cache" {
		t.Fatalf("Flags: got %q, want tmp|cache", resp.Flags)
	}
	if resp.Owner != "codex" {
		t.Fatalf("Owner: got %q, want codex", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
