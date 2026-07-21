# bak-files — access-denied → warning (coverage backfill)

**Coverage backfill** for existing permission / pathflag catalog behavior in
`internal/engine` and `pathflag`:

- Pathflag-skipped top-level dots (e.g. **`.Trash`**) are **not scheduled** as
  jobs; discovery logs **`INFO: skip … (macOS trash)`**
- Unreadable dirs/files during walk (`isAccessDenied` / `walkAccessErr` /
  `logWarn`): **`warning: skip <path>: …`** and **exit 0** (backup continues)
- Other jobs still copy successfully

**Intent:** document and lock GREEN behavior; RED not required.

Out of scope: changing warn format beyond current `warning: skip …`; TCC-specific
macOS sandbox (chmod 000 is the portable stand-in); production code changes
(designer).

Do **not** modify sealed sibling trees unless an implementer reconciles them.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** loads **bak.config** and, with **includeDotFiles** default **true**,
**auto-discovers** top-level **dot names** under **`$HOME`**.

During discovery, if a top-level dot is **pathflag DefaultSkipMask** (e.g.
**`.Trash` → FlagTrash**, reason **macOS trash**), the engine **logs skip** and
**does not** append a Job — so nothing under `.Trash` is walked or copied.

During **walk / copy**, **permission denied** / **operation not permitted**
(including portable **chmod 000** stand-ins for TCC) are classified by
**`isAccessDenied`**: **`logWarn`** emits **`warning: skip <path>: …`**, the
walk **skips** that subtree, and the overall backup **still exits 0**. Readable
siblings (e.g. **`.bashrc`**, **`.other`**) still land under **targetDir**.

An **operator** invokes:

- **`bak-files backup`** — discover dots; skip catalog trash; walk jobs; warn+continue on access denied

Minimal fixtures: **empty `files`**, mapping **`~` → `HOME/$WORKING_ROLE`**,
validate **HOME** + **WORKING_ROLE**, role **`alice`**, **targetDir** `./files`.

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup; inspects logs and targetDir |
| bak-files CLI | Discovers dots; pathflag catalog skip; walkAccessErr / logWarn |
| bak.config | empty files, mapping, targetDir, default includeDotFiles |
| pathflag | `.Trash` → Trash / macOS trash; DefaultSkipMask |
| $HOME tree | Simulated HOME with .Trash, .bashrc, unreadable dirs |
| targetDir | Backup store under WorkDir (`./files`) |
| Environment | HOME, WORKING_ROLE for expand + validate |

## Version

0.0.2

## Decision Tree

```
tests/cli-access-warn/                              [Request{Args,Env,WorkDir,paths…}]
│                                                   Run: session-build + exec bak-files
└── backup/                                         # real backup (may write targetDir)
    ├── trash-catalog-skip/                         # .Trash not scheduled; INFO skip; bashrc ok
    ├── unreadable-dir-warn/                        # chmod 000 dir → warning: skip; exit 0
    └── continues-after-warn/                       # after warn, other good files still backed up
```

**Significance order:** backup surface → concern (catalog skip vs walk access
warn vs continue-after-warn) → concrete fixture.

## Test Index

| Leaf | Description |
|------|-------------|
| `backup/trash-catalog-skip` | HOME `.Trash` + `.bashrc`: exit 0; `INFO: skip` trash; bashrc copied; no `.Trash` under targetDir |
| `backup/unreadable-dir-warn` | HOME `.private/` chmod 000 + `.bashrc`: exit 0; `warning:` skip path; bashrc copied; no fatal `Error:` |
| `backup/continues-after-warn` | Unreadable dir + good `.bashrc` / `.other`: exit 0; good files present under mapping |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-access-warn
doctest test ./tests/cli-access-warn
```

Single leaf:

```bash
doctest test ./tests/cli-access-warn/backup/trash-catalog-skip
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

// Request describes a bak-files backup invocation for access-warn leaves.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path).
	Args []string
	// Env is the complete child environment (KEY=value). Leaves must set Env.
	Env []string
	// WorkDir is process cwd (temp dir with config + trees). Empty → module root.
	WorkDir string

	// Absolute paths for filesystem side-effect assertions.
	TargetDir string // resolved targetDir under WorkDir (e.g. …/files)
	HomeDir   string // simulated $HOME

	// Primary good file (typically .bashrc).
	SourcePath string
	BackupPath string
	Content    string

	// Second good file for continues-after-warn.
	OtherSourcePath string
	OtherBackupPath string
	OtherContent    string

	// Unreadable dir / path under HOME (chmod 000); used for warn leaves.
	UnreadableDir string
	// Expected path that must NOT appear under targetDir (e.g. .Trash or secret).
	ExcludedBackupPath string
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
