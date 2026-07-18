# bak-files — simple backup & restore (P3)

Plan phase **P3**: end-to-end **simple file/dir backup and restore** against
`bak.config` (plain paths, optional entry/`global` excludes, `targetDir` +
`mapping`), plus true **`--dry-run`** (compare and log intent; **zero file
writes**).

Out of scope: hooks, `cmd`/`restoreCmd` modes, `gitTree`/`gitPatch`,
`bak.stats`, `sum.index`, MD5 history conflicts. Do **not** modify sealed
`tests/cli-list` or `tests/cli-foundation`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a **command-line tool** that uses a **bak.config** JSON file to
map **source paths** on the operator filesystem to paths under a backup
**targetDir** (default `./files`) via **mapping** prefixes.

An **operator** invokes:

- **`bak-files backup`** — copy (or skip) sources → **targetDir** under each
  entry’s **mappingPath**
- **`bak-files restore`** — copy from **targetDir** → sources (filesystem)
- **`bak-files backup|restore --help`** — usage only (no config required)

Shared behavior for backup/restore (not help):

1. Resolve config (`--config FILE` or default `bak.config.json` in cwd)
2. Validate environment (`HOME`, `WORKING_ROLE`, plus `validate[].env`)
3. Expand path templates (`~`, `$ENV`) on file keys and mapping
4. For each **simple** files entry (boolean true, string path, or object with
   optional `file` / directory + `excludes`), compute source path and
   destination under `targetDir`/`mappingPath`
5. Apply **global.excludes** and entry **excludes** (basename/glob style as
   implemented) so matching names are not copied
6. If **`--dry-run`**: log what **would** happen (prefer messages containing
   `dry-run` or `would`); **must not create, overwrite, or delete** any files
   under targetDir or restore destinations; exit **0** on success
7. If real run: perform copies; **missing source on backup** should **skip**
   with an informational message (prefer INFO / “skip” / “missing”) and still
   exit **0** (TS-compatible), not fail the whole run
8. Optional: if source and backup content are **identical**, second backup may
   skip rewrite (no content change under targetDir)

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup/restore with flags; inspects logs and trees |
| bak-files CLI | Loads config; maps paths; copies or dry-runs |
| bak.config | files, mapping, targetDir, global.excludes, validate |
| Source tree | Operator filesystem paths (under simulated HOME in tests) |
| targetDir | Backup store under WorkDir (e.g. `./files`) |
| Environment | HOME, WORKING_ROLE for expand + validate |

## Version

0.0.2

## Decision Tree

```
tests/cli-backup/                          [Request{Args, Env, WorkDir, paths…}]
│                                          Run: session-build + exec bak-files
├── help/                                  # no config / no writes
│   ├── backup/                            # backup --help → Usage, exit 0
│   └── restore/                           # restore --help → Usage, exit 0
├── backup/                                # backup subcommand (config + env)
│   ├── real/                              # no --dry-run → may write targetDir
│   │   ├── simple-file/                   # src file → targetDir/mapping path
│   │   ├── with-excludes/                 # excluded name not copied
│   │   ├── missing-source/                # missing src → skip, exit 0
│   │   └── identical-skip/                # second backup leaves content same
│   └── dry-run/                           # --dry-run → zero writes
│       └── no-write/                      # targetDir unchanged / absent
└── restore/                               # restore subcommand
    ├── real/                              # writes filesystem dest
    │   └── simple-file/                   # targetDir → dest file content
    └── dry-run/
        └── no-write/                      # dest not created/modified
```

**Significance order:** subcommand / invocation class (help vs backup vs
restore) → write mode (real vs `--dry-run`) → scenario (simple copy, excludes,
missing source, identical).

## Test Index

| Leaf | Description |
|------|-------------|
| `help/backup` | `backup --help` → exit 0; Usage mentions backup |
| `help/restore` | `restore --help` → exit 0; Usage mentions restore |
| `backup/real/simple-file` | Real backup copies fixture file under `targetDir` |
| `backup/real/with-excludes` | Excluded file not present under `targetDir` after backup |
| `backup/real/missing-source` | Config entry path missing → skip (info), exit 0 |
| `backup/real/identical-skip` | Second backup of same bytes does not alter backup file content |
| `backup/dry-run/no-write` | `backup --dry-run` exit 0; no new/changed files under targetDir |
| `restore/real/simple-file` | Real restore writes dest from seeded `targetDir` |
| `restore/dry-run/no-write` | `restore --dry-run` exit 0; dest unchanged / absent |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-backup
doctest test ./tests/cli-backup
```

Single leaf:

```bash
doctest test ./tests/cli-backup/backup/real/simple-file
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

// Request describes a bak-files backup/restore (or --help) invocation.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path, "--dry-run").
	Args []string
	// Env is the complete child environment (KEY=value). Empty → process default
	// is not used; leaves must set Env (PATH at minimum).
	Env []string
	// WorkDir is process cwd (temp dir with config + trees). Empty → module root.
	WorkDir string

	// Absolute paths for filesystem side-effect assertions (optional per leaf).
	TargetDir   string // resolved targetDir under WorkDir (e.g. …/files)
	SourcePath  string // operator-side source file (backup from / restore to)
	BackupPath  string // expected path under TargetDir after backup
	Content     string // expected file body for success copies
	// Extra paths used by exclude / multi-file leaves
	ExcludedBackupPath string // must NOT exist after backup when excludes apply
	KeepBackupPath     string // must exist after directory backup with excludes
	// Fingerprint of TargetDir before dry-run (empty string = dir must stay absent)
	TargetDirBefore string
	// Fingerprint of SourcePath / dest before restore dry-run
	DestBefore string
	// Optional: content planted in backup store before restore
	SeedBackupContent string
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
	// DOCTEST_ROOT is tests/cli-backup → module root is two levels up.
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
	return filepath.Join(os.TempDir(), "bak-files-cli-backup-"+DOCTEST_SESSION_ID)
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
