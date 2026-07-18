## Expected

- No error.
- `Rule` is `.android/cache`.
- `Flags` is `cache`.
- `Owner` is `android`.
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
	if resp.Rule != ".android/cache" {
		t.Fatalf("Rule: got %q, want .android/cache", resp.Rule)
	}
	if resp.Flags != "cache" {
		t.Fatalf("Flags: got %q, want cache", resp.Flags)
	}
	if resp.Owner != "android" {
		t.Fatalf("Owner: got %q, want android", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
