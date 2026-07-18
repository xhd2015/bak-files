## Expected

- No error.
- `Rule` is `**/upload-chunks`.
- `Flags` is `tmp`.
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
	if resp.Rule != "**/upload-chunks" {
		t.Fatalf("Rule: got %q, want **/upload-chunks", resp.Rule)
	}
	if resp.Flags != "tmp" {
		t.Fatalf("Flags: got %q, want tmp", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
