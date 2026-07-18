## Expected

- No error.
- Zero Result: empty Rule, Reason, Flags, Owner (does not take `**/*.log`).

## Side Effects

- None.

## Errors

- Matching `**/*.log` for `.LOG` fails the test.

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
	if resp.Rule == "**/*.log" {
		t.Fatalf("Classify(%q): .LOG must not match **/*.log", req.RelPath)
	}
	if resp.Rule != "" || resp.Reason != "" || resp.Flags != "" || resp.Owner != "" {
		t.Fatalf("want zero Result for .LOG, got Rule=%q Reason=%q Flags=%q Owner=%q",
			resp.Rule, resp.Reason, resp.Flags, resp.Owner)
	}
}
```
