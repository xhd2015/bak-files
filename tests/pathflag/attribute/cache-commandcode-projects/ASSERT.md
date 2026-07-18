## Expected

- No error.
- `Rule` is `.commandcode/projects`.
- `Flags` is `cache`.
- `Owner` is `commandcode`.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Wrong rule/flags/owner fails.

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
	if resp.Rule != ".commandcode/projects" {
		t.Fatalf("Rule: got %q, want .commandcode/projects", resp.Rule)
	}
	if resp.Flags != "cache" {
		t.Fatalf("Flags: got %q, want cache", resp.Flags)
	}
	if resp.Owner != "commandcode" {
		t.Fatalf("Owner: got %q, want commandcode", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
