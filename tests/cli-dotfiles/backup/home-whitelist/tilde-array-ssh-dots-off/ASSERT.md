# `"~": [".ssh"]` + dots off → only .ssh, not .bashrc, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`.ssh/config`) exists with **KeepContent**.
- **ExcludedBackupPath** (`.bashrc`) does **not** exist.

## Side Effects

- Only expanded `~` array entries (here `.ssh`) are jobs when dots are off.

## Errors

- Missing `.ssh` backup or present `.bashrc` under store fails.

## Exit Code

- **0**

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	got := readFileOrEmpty(req.KeepBackupPath)
	if got != req.KeepContent {
		t.Fatalf("keep %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.KeepBackupPath, got, req.KeepContent, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("dots-off must not backup unlisted .bashrc: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
}
```
