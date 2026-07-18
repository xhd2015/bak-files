## Expected

- No error.
- `Rule` is `.local/share/cursor-agent/versions`.
- `Flags` is `cache|binary`.
- `Owner` is `cursor`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong order (`binary|cache`) or missing bits fails.

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
	if resp.Rule != ".local/share/cursor-agent/versions" {
		t.Fatalf("Rule: got %q, want .local/share/cursor-agent/versions", resp.Rule)
	}
	if resp.Flags != "cache|binary" {
		t.Fatalf("Flags: got %q, want cache|binary", resp.Flags)
	}
	if resp.Owner != "cursor" {
		t.Fatalf("Owner: got %q, want cursor", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
