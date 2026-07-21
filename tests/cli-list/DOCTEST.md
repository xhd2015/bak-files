# bak-files — list + config load (P2)

Plan phase **P2**: load `bak.config` JSON (validate envs, `files`, `mapping`),
resolve each entry to a **mapping path**, and implement **`bak-files list`** so
it prints those paths one per line.

Out of scope for this tree: backup/restore, hooks, git modes, dry-run, writing
under `targetDir`. Do **not** modify sealed `tests/cli-foundation`.

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a **command-line tool** that reads a **bak.config** JSON file
and resolves a **file map** into backup **entries**.

An **operator** invokes **`bak-files list`** with optional **`--config FILE`**
and optional positional **NAMES** filters. The **CLI**:

1. Locates the config (`--config` or default `bak.config.json` under cwd)
2. Parses JSON into **config** (`validate`, `files`, `mapping`, `targetDir`, …)
3. **Validates environment** — always requires **`HOME`** and **`WORKING_ROLE`**;
   also any env names listed under `config.validate[].env`
4. Expands path templates (`~`, `$ENV`) on file keys and mapping keys/values
5. For each file entry, computes a **mappingPath** (logical path under the
   backup store after applying `mapping` prefixes)
6. Optionally **filters** entries whose `mappingPath` matches any **NAME**
   (exact match, or `prefix*` / `*suffix` / `*` as in the private TS tool)
7. Prints each remaining **mappingPath** on its own **stdout** line (stable
   order: declaration order of `files` keys / resolved items), exit **0**

Failures (missing config file, invalid JSON, missing required env) print a clear
message on **stderr** and exit **non-zero**. **`list --help`** prints usage on
stdout and exits **0** without requiring config or env.

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes `bak-files list` with args/env; reads paths or errors |
| bak-files CLI | Parses flags; loads config; validates env; lists mapping paths |
| bak.config | JSON: validate, files, mapping, targetDir |
| Environment | Supplies HOME, WORKING_ROLE, and any validate[] envs for expand/check |

## Version

0.0.2

## Decision Tree

```
tests/cli-list/                         [Request{Args, Env, WorkDir, …}]
│                                       Run: build once + exec bak-files in WorkDir
├── help/
│   └── long-help/                      # list --help → exit 0, Usage (no config)
├── errors/                             # list run fails before listing
│   ├── missing-config-file/            # --config path does not exist
│   ├── invalid-json/                   # config is not valid JSON
│   ├── missing-env-builtin/            # HOME/WORKING_ROLE unset
│   └── missing-env-validate/           # config.validate requires extra env
└── success/                            # valid config + env → mapping paths
    ├── all-entries/                    # --config fixture; full list, exit 0
    ├── filter-names/                   # positional NAME filters mappingPath
    └── default-config/                 # cwd bak.config.json without --config
```

**Significance order:** invocation kind (help vs list-run) → failure class
(config I/O | JSON | env) vs success → config location (`--config` vs default) /
NAMES filter on success paths.

## Test Index

| Leaf | Description |
|------|-------------|
| `help/long-help` | `list --help` → exit 0; Usage mentions list (and preferably `--config`) |
| `errors/missing-config-file` | `--config` missing path → non-zero; error about config/file |
| `errors/invalid-json` | Bad JSON file → non-zero; parse/config error |
| `errors/missing-env-builtin` | Valid JSON, `WORKING_ROLE` unset → non-zero; mentions env |
| `errors/missing-env-validate` | Builtin envs set; `validate` needs `W0` missing → non-zero |
| `success/all-entries` | Fixture + env; stdout is expected mapping paths (one per line) |
| `success/filter-names` | Same fixture; NAME filters to subset of mapping paths |
| `success/default-config` | `bak.config.json` in cwd; `list` without `--config` lists paths |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-list
doctest test ./tests/cli-list
```

Or a single leaf:

```bash
doctest test ./tests/cli-list/success/all-entries
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

// Request describes a bak-files list (or list --help) invocation.
// Setup fills Args, Env, and WorkDir; Run executes the binary.
type Request struct {
	// Args is the full argv after the binary name (e.g. "list", "--config", path).
	Args []string
	// Env is the complete child environment (KEY=value entries). Empty means a
	// minimal default (PATH only) — leaves that need HOME must set Env.
	Env []string
	// WorkDir is the process cwd (temp dir with fixtures). Empty → module root.
	WorkDir string
}

// Response captures CLI I/O from the list harness.
type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	// BinaryPath is the session-built bak-files binary.
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
