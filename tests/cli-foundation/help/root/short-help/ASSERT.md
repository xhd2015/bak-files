# `-h` → root Usage

## Expected

- Exit code **0**.
- **Stdout** contains Usage and subcommands `backup`, `restore`, `list`.

## Side Effects

- None.

## Errors

- Harness error or wrong exit / missing tokens → fail.

## Exit Code

- **0**

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout
	if !containsFold(out, "usage") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	for _, cmd := range []string{"backup", "restore", "list"} {
		if !strings.Contains(out, cmd) {
			t.Fatalf("stdout missing subcommand %q:\n%s", cmd, out)
		}
	}
}
```
