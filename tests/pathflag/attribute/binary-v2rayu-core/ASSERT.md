## Expected

- No error.
- `Rule` is `.V2rayU/v2ray-core`.
- `Flags` is `binary`.
- `Owner` is empty.
- `Reason` non-empty.

## Side Effects

- None.

## Errors

- Flagging whole `.V2rayU` as rule or wrong flags fails.

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
	if resp.Rule != ".V2rayU/v2ray-core" {
		t.Fatalf("Rule: got %q, want .V2rayU/v2ray-core", resp.Rule)
	}
	if resp.Flags != "binary" {
		t.Fatalf("Flags: got %q, want binary", resp.Flags)
	}
	if resp.Owner != "" {
		t.Fatalf("Owner: got %q, want empty", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
