# bak-files — default home dots + pathflag skips + CLI/config overrides

**Classic TDD** for automatic **home top-level dotfile discovery**, **pathflag**
walk skips under `$HOME`, and **`--no-dot-files` / `--include` / `--exclude`**
plus config `global.includeDotFiles` / `dotIncludes` / `dotExcludes`.

Out of scope: binary content detection, private bak.config, hooks/git modes.
Do **not** modify sealed `tests/cli-*` or `tests/pathflag`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a CLI that loads **bak.config** and, beyond explicit **files**
entries, may **auto-discover** top-level **dot names** under **`$HOME`** when
**includeDotFiles** is enabled (default **true** when the key is omitted).

An **operator** invokes:

- **`bak-files backup`** — resolve explicit jobs + optional auto-dot jobs; walk
  sources under `$HOME` applying **skip policy**; copy or **`--dry-run`** log
- **`bak-files restore`** — reverse copy with the **same skip policy** on
  home-relative paths during the walk
- **`bak-files list`** — print **mapping paths** for explicit **and** discovered
  jobs (same discovery rules as backup)
- **`bak-files backup|restore|list --help`** — document flags including
  **`--no-dot-files`**, **`--include`**, **`--exclude`**

**Discovery** (when dots enabled: default / config true; disabled by
`--no-dot-files` or `global.includeDotFiles: false`):

1. Read `$HOME` top-level entries whose names start with `.`
2. If no explicit job already covers that source, add a job with key `~/name`,
   source `$HOME/name`, mapping via `~` → e.g. `HOME/$WORKING_ROLE`

**Skip policy** (any walk under `$HOME`, home-relative path):

1. force-include (`--include` ∪ `global.dotIncludes`) → keep
2. force-exclude (`--exclude` ∪ `global.dotExcludes`) → skip (**exclude wins**
   if both match)
3. `pathflag.Classify(rel)` with `Flags & DefaultSkipMask != 0` → skip
   (log prefer `INFO: skip <rel> (<reason>)`)
4. existing basename/glob **excludes** → skip
5. else copy / would-copy

**Config** may set `global.includeDotFiles`, `global.dotIncludes`,
`global.dotExcludes`, and existing `global.excludes`. Secrets such as `.ssh`
are **included by default** (not pathflag-skipped unless excluded).

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup/restore/list; sets CLI include/exclude/no-dot-files |
| bak-files CLI | Discovers dots; walks jobs; applies skip policy; copies or dry-runs |
| bak.config | files, mapping, targetDir, global.includeDotFiles / excludes / dots |
| pathflag | Classifies home-relative paths; DefaultSkipMask drives walk skips |
| $HOME tree | Top-level dots and nested paths under simulated HOME |
| targetDir | Backup store under WorkDir (e.g. `./files`) |
| Environment | HOME, WORKING_ROLE for expand + validate |

## Version

0.0.2

## Decision Tree

```
tests/cli-dotfiles/                              [Request{Args,Env,WorkDir,paths…}]
│                                                Run: session-build + exec bak-files
├── help/                                        # usage documents new flags
│   └── backup-flags/                            # backup --help → --no-dot-files, --include, --exclude
├── list/                                        # list uses same discovery as backup
│   └── discovers-dots/                          # empty files + home .bashrc → mapping path listed
├── backup/                                      # backup + discovery / filters / walk
│   ├── discovery/                               # auto-dot jobs (empty or sparse files)
│   │   ├── default-on/                          # dots on: would/copy .bashrc; skip .cache
│   │   ├── config-off/                          # includeDotFiles false → no auto dots
│   │   ├── flag-off/                            # --no-dot-files → no auto dots
│   │   └── dedupe-explicit/                     # explicit ~/.bashrc + discover → one copy
│   ├── filters/                                 # CLI force include/exclude
│   │   ├── include-cache/                       # --include .cache keeps pathflag-cache path
│   │   ├── exclude-ssh/                         # --exclude .ssh skips default-included secret
│   │   └── exclude-wins/                        # include + exclude same path → skip
│   └── walk-skip/                               # nested walk under a discovered/home dir
│       ├── pathflag-partial/                    # .codex/config keep; .codex/.tmp skip
│       └── basename-tmp/                        # global excludes *.tmp still apply with dots
└── restore/                                     # restore skip policy
    └── dry-run-pathflag-skip/                   # dry-run restore logs skip for pathflag tmp
```

**Significance order:** subcommand / surface (help | list | backup | restore) →
backup concern (discovery enablement | force filters | walk-time skip) →
concrete scenario (default / config / flag / dedupe / include / exclude / wins /
pathflag partial / basename / restore dry-run).

## Test Index

| Leaf | Description |
|------|-------------|
| `help/backup-flags` | `backup --help` mentions `--no-dot-files`, `--include`, `--exclude` |
| `list/discovers-dots` | Empty `files` + home `.bashrc`; list prints mapping path for it |
| `backup/discovery/default-on` | Dry-run: would handle `.bashrc`; skip `.cache`; no targetDir writes |
| `backup/discovery/config-off` | `includeDotFiles: false` → no auto-dot backup of `.bashrc` |
| `backup/discovery/flag-off` | `--no-dot-files` → no auto-dot backup of `.bashrc` |
| `backup/discovery/dedupe-explicit` | Explicit `~/.bashrc` + discover → single copy under mapping path |
| `backup/filters/include-cache` | `--include .cache` force-keeps cache tree content |
| `backup/filters/exclude-ssh` | `--exclude .ssh` skips `.ssh` while other dots backup |
| `backup/filters/exclude-wins` | Both `--include` and `--exclude` on `.cache` → skip |
| `backup/walk-skip/pathflag-partial` | Under `.codex`: config copied; `.tmp` not |
| `backup/walk-skip/basename-tmp` | With dots on, `global.excludes` `*.tmp` still excludes basenames |
| `restore/dry-run-pathflag-skip` | Restore `--dry-run` prefers skip log for pathflag-tmp path; no dest write |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-dotfiles
doctest test ./tests/cli-dotfiles
```

Single leaf:

```bash
doctest test ./tests/cli-dotfiles/backup/discovery/default-on
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

// Request describes a bak-files invocation for dotfile discovery / skip leaves.
// Setup fills Args, Env, WorkDir, and absolute paths used by Assert.
type Request struct {
	// Args is argv after the binary (e.g. "backup", "--config", path, "--dry-run").
	Args []string
	// Env is the complete child environment (KEY=value). Leaves must set Env.
	Env []string
	// WorkDir is process cwd (temp dir with config + trees). Empty → module root.
	WorkDir string

	// Absolute paths for filesystem side-effect assertions (optional per leaf).
	TargetDir  string // resolved targetDir under WorkDir (e.g. …/files)
	HomeDir    string // simulated $HOME
	SourcePath string // operator-side path of primary source of interest
	BackupPath string // expected path under TargetDir for primary keep path
	Content    string // expected body for successful copy of primary path

	// KeepBackupPath / ExcludedBackupPath for multi-path leaves.
	KeepBackupPath     string
	ExcludedBackupPath string
	KeepContent        string

	// Fingerprint of TargetDir before dry-run (empty string = dir must stay absent).
	TargetDirBefore string
	// Fingerprint / content of restore dest before dry-run.
	DestBefore string
	DestPath   string // restore destination that must stay absent/unchanged

	// ExpectedMappingPath is a substring or line expected in list stdout.
	ExpectedMappingPath string
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
	// DOCTEST_ROOT is tests/cli-dotfiles → module root is two levels up.
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
	return filepath.Join(os.TempDir(), "bak-files-cli-dotfiles-"+DOCTEST_SESSION_ID)
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
