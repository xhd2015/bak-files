## Expected

- No error.
- `Rule` is `.cisco/vpn/log`.
- `Flags` is `logs`.
- `Owner` is `cisco`.
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
	if resp.Rule != ".cisco/vpn/log" {
		t.Fatalf("Rule: got %q, want .cisco/vpn/log", resp.Rule)
	}
	if resp.Flags != "logs" {
		t.Fatalf("Flags: got %q, want logs", resp.Flags)
	}
	if resp.Owner != "cisco" {
		t.Fatalf("Owner: got %q, want cisco", resp.Owner)
	}
	if resp.Reason == "" {
		t.Fatalf("Reason: want non-empty")
	}
}
```
