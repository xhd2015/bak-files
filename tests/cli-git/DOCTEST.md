# bak-files — git-backed modes (P5)

Plan phase **P5**: **working-tree-aware backup** (`gitTree`), **patch
generation** (`gitPatch`), **`bak.stats`**, and **`--dry-run`** that never
writes patch/stats or mutates remotes.

**Exit criteria (pragmatic MVP):**

1. **`gitTree` dirty:** uncommitted change → copy files under `targetDir`; write
   **`bak.stats`** with commit hash / `hasChange`
2. **`gitTree` clean:** clean worktree → **skip** backup (INFO); no file churn
3. **`--dry-run` `gitTree`:** no writes under `targetDir`, no `bak.stats` change
4. **`gitPatch` dirty:** produces a **`.patch`** (content from `git diff HEAD`)
   under the mapping path
5. **`--dry-run` `gitPatch`:** no patch file written
6. **Restore after gitTree backup:** plain file restore of backed content works

**Out of scope:** full apply-patch ecosystem; private-repo migration; **network
remotes** (`git push` / `pull` / `fetch` must not be required for green tests).
Fixtures use **local-only** `git init` repos (no `origin`).

Do **not** modify sealed trees: `tests/cli-foundation`, `tests/cli-list`,
`tests/cli-backup`, `tests/cli-hooks`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a CLI that loads **bak.config** and, for directory **files**
entries, may treat the source as a **Git repository**:

- **`gitTree`** (`true` or `{ "excludes": [...] }`): inspect the **worktree**.
  If there are **uncommitted changes** (dirty), copy the relevant files under
  **targetDir**/**mappingPath** and record status in **`bak.stats`** (cwd file)
  keyed by mapping path — including **commitHash** and **hasChange**. If the
  tree is **clean**, **skip** the copy (prefer INFO log mentioning clean /
  skip) and avoid churn under targetDir / stats for that entry.
- **`gitPatch`** (`true` or object): when the worktree differs from the chosen
  base (MVP: **`HEAD`** for local fixtures), write a **single patch file** under
  the mapping path (content compatible with `git diff HEAD`, e.g. contains
  `diff --git` or unified-diff markers). When clean / no diff, may skip writing
  a patch (dirty path is required for this tree’s success leaf).

**`--dry-run`:** resolve and log intent (`dry-run` / `would`); **must not**
create or change files under **targetDir**, **must not** write/update
**`bak.stats`**, and **must not** call network git remotes.

**Restore (minimal):** after a prior real **gitTree** backup has placed files
under targetDir, **`bak-files restore`** copies them back to the operator path
like a plain directory restore (no FF/checkout required in P5 MVP).

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup/restore; inspects targetDir, bak.stats, logs |
| bak-files CLI | Loads config; git status / diff; copy or skip; stats I/O |
| bak.config | files entries with `gitTree` / `gitPatch` |
| Local git repo | `git init` fixture under simulated HOME (no remote) |
| targetDir | Backup store (`./files` under WorkDir) |
| bak.stats | JSON stats in process cwd (WorkDir) after real gitTree dirty |
| Environment | HOME, WORKING_ROLE; PATH includes `git` |

Network policy: tests never configure `origin` or require connectivity.

## Version

0.0.2

## Decision Tree

```
tests/cli-git/                                      [Request{Args,Env,WorkDir,git paths…}]
│                                                   Run: session-build + exec bak-files
├── backup/                                         # backup + git modes
│   ├── gitTree/                                    # entry gitTree: true
│   │   ├── real/                                   # may write targetDir + bak.stats
│   │   │   ├── dirty/                              # uncommitted change → copy + stats
│   │   │   └── clean/                              # clean → skip, no churn
│   │   └── dry-run/                                # --dry-run
│   │       └── dirty-no-write/                     # dirty but zero writes / no stats
│   └── gitPatch/                                   # entry gitPatch: true
│       ├── real/
│       │   └── dirty/                              # writes .patch under mapping
│       └── dry-run/
│           └── dirty-no-write/                     # no patch file
└── restore/
    └── real/
        └── after-gitTree-backup/                   # plain restore of backed files
```

**Significance order:** subcommand (backup vs restore) → entry mode (`gitTree`
vs `gitPatch`) → write mode (real vs `--dry-run`) → worktree state
(dirty vs clean) / restore scenario.

## Test Index

| Leaf | Description |
|------|-------------|
| `backup/gitTree/real/dirty` | Dirty worktree: copies tracked content under targetDir; writes `bak.stats` with `hasChange` + commit hash |
| `backup/gitTree/real/clean` | Clean worktree: skip (INFO); no targetDir churn; no new `hasChange:true` stats for mapping |
| `backup/gitTree/dry-run/dirty-no-write` | Dirty + `--dry-run`: exit 0; targetDir + bak.stats unchanged |
| `backup/gitPatch/real/dirty` | Dirty: patch file under mapping with diff-like content; exit 0 |
| `backup/gitPatch/dry-run/dirty-no-write` | Dirty + `--dry-run`: no patch under targetDir |
| `restore/real/after-gitTree-backup` | Seeded targetDir content restored to source path; exit 0 |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-git
doctest test ./tests/cli-git
```

Single leaf:

```bash
doctest test ./tests/cli-git/backup/gitTree/real/dirty
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

// Request describes a bak-files git-mode backup/restore invocation.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path, "--dry-run").
	Args []string
	// Env is the complete child environment (KEY=value).
	Env []string
	// WorkDir is process cwd (temp dir with config + git fixtures).
	WorkDir string

	// Absolute paths for filesystem side-effect assertions.
	TargetDir  string // resolved targetDir under WorkDir (…/files)
	SourcePath string // git repo root (operator source directory)
	BackupPath string // expected path under TargetDir after backup (file)
	Content    string // expected file body for success copies / restore
	// StatsPath is WorkDir/bak.stats (cwd-relative file written by real gitTree).
	StatsPath string
	// MappingKey is the bak.stats object key (mapping path, e.g. "HOME/repo").
	MappingKey string
	// CommitHash is HEAD of the fixture repo at Setup time (full or abbreviated).
	CommitHash string
	// Fingerprints before dry-run / clean assertions.
	TargetDirBefore string
	StatsBefore     string
	// DestBefore is source-side fingerprint before restore dry-run (unused MVP).
	DestBefore string
	// PatchGlobDir is directory expected to contain a .patch for gitPatch real.
	// When empty, BackupPath itself is the patch file path.
	PatchDir string
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
	// DOCTEST_ROOT is tests/cli-git → module root is two levels up.
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
	return filepath.Join(os.TempDir(), "bak-files-cli-git-"+DOCTEST_SESSION_ID)
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
