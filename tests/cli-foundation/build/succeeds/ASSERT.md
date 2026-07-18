# go build ./cmd/bak-files succeeds

## Expected

- Build process exit code **0**.
- `BuildSucceeded` is true (output binary exists and is a file).
- `BinaryPath` is non-empty.

## Side Effects

- Binary written only under the test temp dir (not module root).

## Errors

- Compile failure or missing binary fails the leaf.

## Exit Code

- **0** (from `go build`)

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v\nstdout:\n%s\nstderr:\n%s", err, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("go build exit %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !resp.BuildSucceeded {
		t.Fatalf("build claimed success but binary missing: %q\nstderr:\n%s",
			resp.BinaryPath, resp.Stderr)
	}
	if resp.BinaryPath == "" {
		t.Fatal("BinaryPath empty")
	}
}
```
