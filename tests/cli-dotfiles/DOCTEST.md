# bak-files — home whitelist, general files array whitelist, default dots, pathflag, CLI

**Classic TDD** for **whitelist-driven** backup policy, automatic **top-level
dot discovery** (`includeDotFiles` ≡ synthetic **`~/.*`**), **bare `"~":
[names]`** expansion (never full-home `copyDir`), **any `PREFIX: [names]`
string-array whitelist** (never full-PREFIX `copyDir`), **pathflag** walk skips
under `$HOME`, and **`--no-dot-files` / `--include` / `--exclude`** plus config
`global.includeDotFiles` / `dotIncludes` / `dotExcludes`.

Out of scope: binary content detection, private bak.config, hooks/git modes,
pathflag catalog changes, ExpandPath changes for `~/foo`.
Do **not** modify sealed sibling `tests/cli-*` or `tests/pathflag`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a CLI that loads **bak.config** and resolves **files** into
**copy jobs** under a **whitelist-only** model: only paths matching explicit
entries (after expand) plus optional synthetic home-dot discovery are backed up.

An **operator** invokes:

- **`bak-files backup`** — resolve explicit jobs + optional auto-dot jobs; walk
  sources under `$HOME` applying **skip policy**; copy or **`--dry-run`** log
- **`bak-files restore`** — reverse copy with the **same skip policy** on
  home-relative paths during the walk
- **`bak-files list`** — print **mapping paths** for explicit **and** discovered
  jobs (same ResolveJobs rules as backup)
- **`bak-files backup|restore|list --help`** — document flags including
  **`--no-dot-files`**, **`--include`**, **`--exclude`**

**Bare `"~"` key (home basename whitelist)** — **never** a recursive job of
`$HOME`:

| Config | Behavior |
|--------|----------|
| `"~": ["name", …]` | Each **name** → ordinary job key `~/name`, source `$HOME/name`, mapping via `~` prefix (e.g. `HOME/$WORKING_ROLE/name`) |
| `"~": [".bashrc"]` etc. | Dot names are whitelist entries (subset of `~/.*` when dots on) |
| Non-dot in array (e.g. `"Notes"`) | Whitelist addition `~/Notes` |
| Array ignored + `ExpandPath("~")` → `$HOME` | **Bug** (historical): full-home `copyDir` |

There is **no home-root job**. Trees such as `Library/`, `Downloads/` are
**not** backed up unless listed (via `"~"` array or an explicit `~/…` key).

**Any `files` key PREFIX with string array value** — same rule, not only `"~"`:

| Config | Behavior |
|--------|----------|
| `"$W0/proj": [".vscode"]` | One job: key `$W0/proj/.vscode`, source `expand($W0)/proj/.vscode`, mapping e.g. `W/proj/.vscode` |
| `PREFIX: [name, …]` | Each basename → `PREFIX/name` job; **never** Source=full expanded PREFIX alone |
| Object `{"excludes":[…]}` / `true` | Still full PREFIX with optional excludes — **unchanged** |
| Array only for `"~"`, other PREFIX full-tree | **Bug** (current): vendor poison — **files-array-whitelist** leaves RED until fixed |

**Discovery / synthetic `~/.*`** (when dots enabled: default / config true;
disabled by `--no-dot-files` or `global.includeDotFiles: false`):

1. Equivalent to a synthetic whitelist of top-level `$HOME` basenames starting
   with `.` (existing discovery)
2. If no explicit job already covers that source, add a job with key `~/name`,
   source `$HOME/name`, mapping via `~` → e.g. `HOME/$WORKING_ROLE`
3. When dots **off**, only explicit entries remain — including names expanded
   from the `"~"` array

**Skip policy** (any walk under an already-included `$HOME` job tree,
home-relative path) — filters only **narrow** included trees:

1. force-include (`--include` ∪ `global.dotIncludes`) → keep
2. force-exclude (`--exclude` ∪ `global.dotExcludes`) → skip (**exclude wins**
   if both match)
3. `pathflag.Classify(rel)` with `Flags & DefaultSkipMask != 0` → skip
   (log prefer `INFO: skip <rel> (<reason>)`)
4. existing basename/glob **excludes** → skip
5. else copy / would-copy

**Config** may set `global.includeDotFiles`, `global.dotIncludes`,
`global.dotExcludes`, and existing `global.excludes`. Secrets such as `.ssh`
are **included by default** when discovered or listed (not pathflag-skipped
unless excluded).

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes backup/restore/list; sets CLI include/exclude/no-dot-files |
| bak-files CLI | Expands `"~"` / any PREFIX string arrays; discovers dots; walks jobs; skip policy; copies |
| bak.config | files (incl. `"~": [names]`, `"$W0/…": [names]`), mapping, targetDir, global dots/excludes |
| pathflag | Classifies home-relative paths; DefaultSkipMask drives walk skips |
| $HOME tree | Top-level dots, whitelist basenames, nested paths under simulated HOME |
| targetDir | Backup store under WorkDir (e.g. `./files`) |
| Environment | HOME, WORKING_ROLE for expand + validate |

## Version

0.0.4

## Decision Tree

```
tests/cli-dotfiles/                              [Request{Args,Env,WorkDir,paths…}]
│                                                Run: session-build + exec bak-files
├── help/                                        # usage documents new flags
│   └── backup-flags/                            # backup --help → --no-dot-files, --include, --exclude
├── list/                                        # list uses same ResolveJobs as backup
│   ├── discovers-dots/                          # empty files + home .bashrc → mapping path listed
│   ├── tilde-array-expands/                     # "~": [".bashrc"] → HOME/alice/.bashrc, not bare HOME/alice
│   └── project-array-expands/                   # "$W0/proj": [".vscode"] → W/proj/.vscode, not bare W/proj  [RED until general array]
├── backup/                                      # backup + discovery / filters / walk / whitelist
│   ├── discovery/                               # auto-dot jobs (empty or sparse files) ≡ ~/.*
│   │   ├── default-on/                          # dots on: would/copy .bashrc; skip .cache
│   │   ├── config-off/                          # includeDotFiles false → no auto dots
│   │   ├── flag-off/                            # --no-dot-files → no auto dots
│   │   └── dedupe-explicit/                     # explicit ~/.bashrc + discover → one copy
│   ├── home-whitelist/                          # bare "~" array must not full-home copy
│   │   ├── tilde-no-library/                    # "~": [".bashrc"] + Library poison → no Library
│   │   ├── tilde-array-ssh-dots-off/            # dots off + "~": [".ssh"] → .ssh only, not .bashrc
│   │   ├── tilde-array-notes/                   # dots off + "~": ["Notes"] → Notes only, not Library
│   │   └── explicit-scripts/                    # "~/Scripts": true → Scripts; Library out
│   ├── files-array-whitelist/                   # any PREFIX: [names] — not only bare "~"  [RED until general array]
│   │   └── project-vscode-only/                 # "$W0/proj": [".vscode"] + vendor poison → no vendor
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
backup concern (discovery | **home-whitelist** | **files-array-whitelist** |
force filters | walk-time skip) → concrete scenario.

## Test Index

| Leaf | Description |
|------|-------------|
| `help/backup-flags` | `backup --help` mentions `--no-dot-files`, `--include`, `--exclude` |
| `list/discovers-dots` | Empty `files` + home `.bashrc`; list prints mapping path for it |
| `list/tilde-array-expands` | `"~": [".bashrc"]` dots off → list `HOME/alice/.bashrc`, not bare `HOME/alice` |
| `list/project-array-expands` | `"$W0/proj": [".vscode"]` → list `W/proj/.vscode`, not bare `W/proj` (**RED** until general array) |
| `backup/discovery/default-on` | Dry-run: would handle `.bashrc`; skip `.cache`; no targetDir writes |
| `backup/discovery/config-off` | `includeDotFiles: false` → no auto-dot backup of `.bashrc` |
| `backup/discovery/flag-off` | `--no-dot-files` → no auto-dot backup of `.bashrc` |
| `backup/discovery/dedupe-explicit` | Explicit `~/.bashrc` + discover → single copy under mapping path |
| `backup/home-whitelist/tilde-no-library` | `"~": [".bashrc"]` backs bashrc; never Library (no full-home) |
| `backup/home-whitelist/tilde-array-ssh-dots-off` | Dots off + `"~": [".ssh"]` → only `.ssh`; not `.bashrc` |
| `backup/home-whitelist/tilde-array-notes` | Dots off + `"~": ["Notes"]` → Notes only; not Library |
| `backup/home-whitelist/explicit-scripts` | Explicit `~/Scripts` works; Library not pulled |
| `backup/files-array-whitelist/project-vscode-only` | `"$W0/proj": [".vscode"]` keeps settings; never vendor (**RED** until general array) |
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
