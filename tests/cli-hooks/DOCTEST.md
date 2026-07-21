# bak-files — hooks & command-generated entries (P4)

Plan phase **P4**: **shell hooks** (`beforeCopy` / `afterCopy` on backup;
`beforeRestore` / related on restore) and **command-generated** file entries
(`cmd` on backup, `restoreCmd` on restore) via process execution.

**Exit criteria (from plan):**

1. Fixture with **`cmd`**: real backup writes mapping-path file from command **stdout**
2. Fixture with **`beforeCopy`**: real backup **runs** the script (marker file appears)
3. **`--dry-run`**: does **not** execute hooks/`cmd`; no marker, no generated file
4. **Failing** `beforeCopy` → process exit **non-zero**
5. Optional easy path: **`restoreCmd`** (or `beforeRestore`) produces a side effect

Out of scope: `gitTree`, `gitPatch`, `bak.stats`, MD5 history. Do **not** modify
sealed trees `tests/cli-foundation`, `tests/cli-list`, or `tests/cli-backup`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a CLI that loads **bak.config** and, for each **files** entry,
either copies a filesystem path or **generates** content by running a shell
command. Entries may also attach **lifecycle hooks** that run as shell snippets
around copy/generate.

An **operator** invokes:

- **`bak-files backup`** — for each entry:
  - object with **`cmd`**: run shell command; capture **stdout** → file under
    **targetDir**/**mappingPath** (real run only)
  - object with **`file`** (or path-like entry) plus **`beforeCopy`** /
    **`afterCopy`**: run those scripts around the copy (real run only)
  - simple path entries: plain copy (covered by P3; not re-tested here)
- **`bak-files restore`** — may run **`beforeRestore`** / **`restoreCmd`**
  (execute; stdout may be ignored) on real restore
- **`--dry-run`**: resolve and **log intent** (prefer `dry-run` / `would` /
  would-run style messages); **must not** spawn hooks or `cmd` / `restoreCmd`;
  **must not** write generated or copied artifacts for those entries

Shared pipeline (backup/restore, not help):

1. Load config (`--config` or default)
2. Validate env (`HOME`, `WORKING_ROLE`, `validate[]`)
3. Expand paths; resolve mapping under `targetDir`
4. For each entry: if dry-run → log would-run / would-write, **no exec, no write**;
   else run hooks/`cmd`/`restoreCmd` and perform side effects / file writes
5. Hook non-zero exit → whole command fails **non-zero**

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup/restore; inspects markers, generated files, exit codes |
| bak-files CLI | Loads config; spawns shells for hooks/cmd when not dry-run |
| bak.config | `files` entries with `cmd`, `beforeCopy`, `afterCopy`, `restoreCmd`, `beforeRestore` |
| Shell | Executes safe scripts (`echo`, `touch`, `exit 1`) in WorkDir / temp |
| Marker / side-effect files | Prove hooks actually ran (or did not on dry-run) |
| targetDir | Receives `cmd` stdout files and normal copies |
| Environment | HOME, WORKING_ROLE for expand + validate |

Safe scripts only: `echo` / `printf`, `touch`, `exit 1` — no network, no sudo.

## Version

0.0.2

## Decision Tree

```
tests/cli-hooks/                                   [Request{Args,Env,WorkDir,paths…}]
│                                                  Run: session-build + exec bak-files
├── backup/                                        # backup + hooks/cmd
│   ├── real/                                      # no --dry-run → may exec + write
│   │   ├── cmd-stdout/                            # {cmd} stdout → mapping path file
│   │   ├── beforeCopy-marker/                     # beforeCopy runs → marker exists
│   │   └── beforeCopy-fail/                       # beforeCopy non-zero → exit ≠ 0
│   └── dry-run/                                   # --dry-run → no exec, no write
│       └── no-exec/                               # cmd + beforeCopy both inert
└── restore/                                       # restore + restoreCmd
    └── real/
        └── restoreCmd-side-effect/                # restoreCmd runs → side-effect file
```

**Significance order:** subcommand (backup vs restore) → write/exec mode
(real vs `--dry-run`) → entry mechanism / outcome (`cmd` success, `beforeCopy`
success, `beforeCopy` fail, `restoreCmd` success).

## Test Index

| Leaf | Description |
|------|-------------|
| `backup/real/cmd-stdout` | Real backup: `cmd` stdout becomes file under `targetDir`/mapping |
| `backup/real/beforeCopy-marker` | Real backup: `beforeCopy` creates marker file; exit 0 |
| `backup/real/beforeCopy-fail` | Real backup: failing `beforeCopy` → non-zero exit |
| `backup/dry-run/no-exec` | `backup --dry-run`: no marker from beforeCopy; no cmd-generated file |
| `restore/real/restoreCmd-side-effect` | Real restore: `restoreCmd` creates side-effect file; exit 0 |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-hooks
doctest test ./tests/cli-hooks
```

Single leaf:

```bash
doctest test ./tests/cli-hooks/backup/real/cmd-stdout
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

// Request describes a bak-files backup/restore invocation for hooks/cmd leaves.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path, "--dry-run").
	Args []string
	// Env is the complete child environment (KEY=value). Leaves must set Env.
	Env []string
	// WorkDir is process cwd (temp dir with config). Empty → module root.
	WorkDir string

	// Absolute paths for filesystem side-effect assertions.
	TargetDir  string // resolved targetDir under WorkDir (e.g. …/files)
	SourcePath string // operator-side source file when copying
	BackupPath string // expected path under TargetDir after backup (cmd or copy)
	Content    string // expected body of BackupPath after successful cmd/copy

	// MarkerPath is a file that only a real hook (beforeCopy/afterCopy) creates.
	MarkerPath string
	// MarkerBefore is fingerprint of MarkerPath before run (usually "" absent).
	MarkerBefore string

	// SideEffectPath is a file created by restoreCmd / beforeRestore (restore leaves).
	SideEffectPath string
	// SideEffectBefore fingerprint before restore (usually "").
	SideEffectBefore string

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
