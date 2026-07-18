# bak-files — symlink preserve (coverage backfill)

**Coverage backfill** for existing correct symlink handling in
`internal/engine` (`copySymlink` / `sameContent` via Lstat+Readlink). Real
backup **preserves** OS symlinks (same target string) and does **not** follow
self-referential symlink-to-dir into `Open` / `"is a directory"` failures.

**Intent:** document and lock GREEN behavior; RED not required.

Out of scope: portable absolute→relative rewrite on restore; Windows symlink
semantics; production code changes (designer).

Do **not** modify sealed sibling trees unless an implementer reconciles them.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** walks directory sources with **Lstat** semantics for entries that
are **symlinks**:

1. **ModeSymlink** → **`copySymlink`**: `Readlink` source; on real run create
   an OS symlink at the mapping destination with the **same target string**;
   do **not** `Open`/follow the link.
2. **Self-referential symlink-to-dir** (e.g. `proj/proj` → abs path of `proj`)
   must **exit 0**; destination link preserved; combined logs must **not**
   contain the failure token **`is a directory`**.
3. **File symlink** (link → regular file): destination is a symlink (not
   necessarily the target’s bytes written as a non-link open-copy).
4. **`--dry-run`**: log intent containing **`would symlink`** (or dry-run +
   symlink wording); **zero writes** under **targetDir**.
5. **Second backup** when the backup store already has the same link target:
   exit **0**; optional identical skip when verbose (`identical symlink`).

Minimal fixtures: **`includeDotFiles: false`**, **`files: {"~/Scripts": true}`**,
mapping **`~` → `HOME/$WORKING_ROLE`**, validate **HOME** + **WORKING_ROLE**.

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup; inspects logs and targetDir links |
| bak-files CLI | Walks sources; copySymlink / dry-run would symlink |
| bak.config | files ~/Scripts, mapping, targetDir, includeDotFiles false |
| Source tree | Simulated HOME with Scripts tree + symlinks |
| targetDir | Backup store under WorkDir (`./files`) |
| Environment | HOME, WORKING_ROLE for expand + validate |

## Version

0.0.2

## Decision Tree

```
tests/cli-symlink/                                 [Request{Args,Env,WorkDir,paths…}]
│                                                  Run: session-build + exec bak-files
└── backup/
    ├── real/                                      # may write targetDir
    │   ├── self-link-dir/                         # dir tree + self symlink-to-dir
    │   ├── file-symlink/                          # file link preserved as symlink
    │   └── identical-skip/                        # second backup of same links
    └── dry-run/
        └── symlink-msg/                           # would symlink; no writes
```

**Significance order:** write mode (real vs `--dry-run`) → symlink kind
(self-dir | file) → second-run stability.

## Test Index

| Leaf | Description |
|------|-------------|
| `backup/real/self-link-dir` | Self-referential symlink-to-dir under Scripts: exit 0; link + file preserved; no `is a directory` |
| `backup/real/file-symlink` | File symlink recreated as symlink under mapping path |
| `backup/real/identical-skip` | Second backup after first: exit 0; symlink target unchanged |
| `backup/dry-run/symlink-msg` | Dry-run mentions would symlink; targetDir unchanged |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-symlink
doctest test ./tests/cli-symlink
```

Single leaf:

```bash
doctest test ./tests/cli-symlink/backup/real/self-link-dir
```

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

// Request describes a bak-files backup invocation for symlink leaves.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path, "--dry-run").
	Args []string
	// Env is the complete child environment (KEY=value). Leaves must set Env.
	Env []string
	// WorkDir is process cwd (temp dir with config + trees). Empty → module root.
	WorkDir string

	// Absolute paths for filesystem side-effect assertions.
	TargetDir string // resolved targetDir under WorkDir (e.g. …/files)
	HomeDir   string // simulated $HOME
	// Regular file under Scripts that must be copied as a normal file.
	SourceFilePath string
	BackupFilePath string
	Content        string
	// Symlink path on operator side and expected path under targetDir.
	SourceLinkPath string
	BackupLinkPath string
	// LinkTarget is the exact Readlink string expected at BackupLinkPath.
	LinkTarget string
	// Fingerprint of TargetDir before dry-run (empty string = dir must stay absent).
	TargetDirBefore string
}

// Response captures CLI I/O from the harness.
type Response struct {
	Stdout     string
	Stderr     string
	ExitCode   int
	BinaryPath string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	bin := ensureBakFilesBinary(t)

	args := req.Args
	if args == nil {
		args = []string{}
	}

	cmd := exec.Command(bin, args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	} else {
		cmd.Dir = moduleRoot(t)
	}
	if req.Env != nil {
		cmd.Env = req.Env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			return &Response{
				Stdout:     stdout.String(),
				Stderr:     stderr.String(),
				BinaryPath: bin,
			}, err
		}
	}
	return &Response{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		ExitCode:   code,
		BinaryPath: bin,
	}, nil
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	// DOCTEST_ROOT is tests/cli-symlink → module root is two levels up.
	root, err := filepath.Abs(filepath.Join(DOCTEST_ROOT, "../.."))
	if err != nil {
		t.Fatalf("module root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("module root missing go.mod at %s: %v", root, err)
	}
	return root
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "bak-files-cli-symlink-"+DOCTEST_SESSION_ID)
}

func withFileLock(t *testing.T, lockPath string, fn func() error) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		t.Fatalf("mkdir cache: %v", err)
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open lock: %v", err)
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatalf("flock: %v", err)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err := fn(); err != nil {
		t.Fatal(err)
	}
}

func ensureBakFilesBinary(t *testing.T) string {
	t.Helper()
	cache := sessionCacheDir()
	bin := filepath.Join(cache, "bak-files")
	ready := filepath.Join(cache, "binaries.ready")
	lock := filepath.Join(cache, "build.lock")
	root := moduleRoot(t)

	withFileLock(t, lock, func() error {
		if fileExists(ready) && fileExists(bin) {
			return nil
		}
		if err := os.MkdirAll(cache, 0o755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/bak-files")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("go build ./cmd/bak-files: %w\n%s", err, out)
		}
		return os.WriteFile(ready, []byte("ok"), 0o644)
	})
	return bin
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}
```
