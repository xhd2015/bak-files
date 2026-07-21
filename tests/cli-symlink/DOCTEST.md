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
	"sync"
	"strings"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
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

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	bin := ensureBakFilesBinary(t, d)

	args := req.Args
	if args == nil {
		args = []string{}
	}

	cmd := exec.Command(bin, args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	} else {
		cmd.Dir = moduleRoot(t, d)
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

func moduleRoot(t *testing.T, d *session.Doctest) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join(d.DOCTEST_ROOT, "../.."))
	if err != nil {
		t.Fatalf("module root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("module root missing go.mod at %s: %v", root, err)
	}
	return root
}

// Process-local binary (one-process suite; in-memory mutex, not session flock).
var (
	ensureBakFilesBinaryMu   sync.Mutex
	ensureBakFilesBinaryPath string
	ensureBakFilesBinaryErr  error
)

func ensureBakFilesBinary(t *testing.T, d *session.Doctest) string {
	t.Helper()
	ensureBakFilesBinaryMu.Lock()
	defer ensureBakFilesBinaryMu.Unlock()
	if ensureBakFilesBinaryPath != "" || ensureBakFilesBinaryErr != nil {
		if ensureBakFilesBinaryErr != nil {
			t.Fatal(ensureBakFilesBinaryErr)
		}
		return ensureBakFilesBinaryPath
	}
	dir, err := os.MkdirTemp("", "ensureBakFilesBinary-")
	if err != nil {
		ensureBakFilesBinaryErr = err
		t.Fatal(err)
	}
	binPath := filepath.Join(dir, "bak-files")
	root := moduleRoot(t, d)
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/bak-files")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		ensureBakFilesBinaryErr = fmt.Errorf("go build ./cmd/bak-files: %w\n%s", err, strings.TrimSpace(string(out)))
		t.Fatal(ensureBakFilesBinaryErr)
	}
	ensureBakFilesBinaryPath = binPath
	return binPath
}


```
