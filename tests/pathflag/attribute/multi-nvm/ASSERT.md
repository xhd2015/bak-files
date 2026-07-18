## Expected

- No error.
- `Rule` is `.nvm`.
- `Flags` is `cache|binary`.
- `Owner` is `nvm`.
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
	if resp.Rule != ".nvm" {
		t.Fatalf("Rule: got %q, want .nvm", resp.Rule)
	}
	if resp.Flags != "cache|binary" {
		t.Fatalf("Flags: got %q, want cache|binary", resp.Flags)
	}
	if resp.Owner != "nvm" {
		t.Fatalf("Owner: got %q, want nvm", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
