# bak-files — CLI foundation (P1)

Plan phase **P1**: public module installable as `bak-files` with multi-level
help, stub subcommands, successful `go build` of `./cmd/bak-files`, and a safe
repo-root `.gitignore`.

Out of scope for this tree: config parsing, file I/O, dry-run, real
backup/restore/list logic (stubs only).

Target binary:

```text
cmd/bak-files → bak-files
```

Module: `github.com/xhd2015/bak-files`

# DSN (Domain Specific Notion)

**bak-files** is a **command-line tool** for backup and restore workflows.
An **operator** (human or shell) invokes the **`bak-files` binary** with an
**argv** of flags and optional subcommand.

The **CLI** (built with `github.com/xhd2015/less-flags`) **dispatches** on the
first token:

- **Root help** — empty argv, `--help`, or `-h` → print **Usage** listing
  subcommands **`backup`**, **`restore`**, **`list`** on **stdout**, exit **0**
- **Subcommand help** — `bak-files <cmd> --help` or `-h` → that command’s
  **Usage** on **stdout**, exit **0**
- **Known stub command** without help — `backup` / `restore` / `list` with no
  help flag → **not implemented** message (stdout or stderr), exit **non-zero**
  (stubs until later phases)
- **Unknown command** → **error** on **stderr** (prefix `Error:` or
  `bak-files:`), exit **non-zero**

The **module root** also ships a **`.gitignore`** that keeps binaries,
coverage artifacts, editor/OS junk, secrets, accidental payload dirs
(`files/`), and index/stats files (`bak.stats`, `sum.index`) out of git.
**Build** (`go build -o … ./cmd/bak-files`) must succeed from the module root.

Participants:

| Participant | Role |
|-------------|------|
| Operator | Invokes `bak-files` with args; reads stdout/stderr and exit code |
| bak-files CLI | Parses flags/subcommands; prints help or stub/error messages |
| Module / go build | Compiles `cmd/bak-files` into a runnable binary |
| .gitignore | Declares ignore patterns so private/junk paths are not committed |

## Version

0.0.2

## Decision Tree

```
tests/cli-foundation/                    [Request{Mode, Args}]
│                                        Run: build once + cli | go build | read .gitignore
├── help/                                # usage paths → exit 0
│   ├── root/
│   │   ├── empty-args/                  # no args → root Usage + backup/restore/list
│   │   ├── long-help/                   # --help
│   │   └── short-help/                  # -h
│   ├── backup/                          # backup --help
│   ├── restore/                         # restore --help
│   └── list/                            # list --help
├── stub/                                # known cmd, no help → not implemented, non-zero
│   ├── backup/
│   ├── restore/
│   └── list/
├── unknown/
│   └── command/                         # garbage subcommand → stderr Error, non-zero
├── build/
│   └── succeeds/                        # go build ./cmd/bak-files → exit 0, binary
└── gitignore/
    └── required-patterns/               # .gitignore covers bin/files/stats/secrets/OS
```

**Significance order:** surface under test (`cli` help | stub | unknown | build |
gitignore) → help target (root vs subcommand) → help form (empty / `--help` / `-h`)
or which stub command.

## Test Index

| Leaf | Description |
|------|-------------|
| `help/root/empty-args` | No args → exit 0; Usage lists backup, restore, list |
| `help/root/long-help` | `--help` → exit 0; same root usage expectations |
| `help/root/short-help` | `-h` → exit 0; same root usage expectations |
| `help/backup` | `backup --help` → exit 0; backup-specific usage |
| `help/restore` | `restore --help` → exit 0; restore-specific usage |
| `help/list` | `list --help` → exit 0; list-specific usage |
| `stub/backup` | `backup` alone → non-zero; not-implemented signal |
| `stub/restore` | `restore` alone → non-zero; not-implemented signal |
| `stub/list` | `list` alone → non-zero; not-implemented signal |
| `unknown/command` | Unknown subcommand → non-zero; Error/bak-files: on stderr |
| `build/succeeds` | `go build -o … ./cmd/bak-files` succeeds; binary exists |
| `gitignore/required-patterns` | Root `.gitignore` has required safety patterns |

## How to Run

From module root:

```bash
doctest vet ./tests/cli-foundation
doctest test ./tests/cli-foundation
```

Or a single leaf:

```bash
doctest test ./tests/cli-foundation/help/root/empty-args
```

```go
import (
	"sync"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Request selects what the harness exercises.
// Mode: "cli" (default), "build", "gitignore".
type Request struct {
	Mode string
	Args []string
}

// Response captures CLI I/O, build result, or .gitignore content.
type Response struct {
	Stdout            string
	Stderr            string
	ExitCode          int
	BinaryPath        string
	GitignoreContent  string
	BuildSucceeded    bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	mode := req.Mode
	if mode == "" {
		mode = "cli"
	}
	root := moduleRoot(t, d)

	switch mode {
	case "gitignore":
		p := filepath.Join(root, ".gitignore")
		b, err := os.ReadFile(p)
		if err != nil {
			return &Response{ExitCode: 1}, fmt.Errorf("read .gitignore: %w", err)
		}
		return &Response{
			GitignoreContent: string(b),
			ExitCode:         0,
		}, nil

	case "build":
		bin := filepath.Join(t.TempDir(), "bak-files-build-check")
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/bak-files")
		cmd.Dir = root
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
					Stdout: stdout.String(),
					Stderr: stderr.String(),
				}, err
			}
		}
		ok := false
		if code == 0 {
			if st, e := os.Stat(bin); e == nil && !st.IsDir() {
				ok = true
			}
		}
		return &Response{
			Stdout:         stdout.String(),
			Stderr:         stderr.String(),
			ExitCode:       code,
			BinaryPath:     bin,
			BuildSucceeded: ok,
		}, nil

	default: // cli
		bin := ensureBakFilesBinary(t, d)
		cmd := exec.Command(bin, req.Args...)
		cmd.Dir = root
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

// containsFold is a case-insensitive substring check for Assert helpers.
func containsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}
```
