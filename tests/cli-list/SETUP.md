# Scenario

**Feature**: bak-files P2 — config load and `list` of resolved mapping paths

```
# operator runs list; CLI loads bak.config, validates env, prints mapping paths
operator -> bak-files list [--config FILE] [NAMES...]
  -> (load config + validate env + resolve mapping)
  -> stdout: mappingPath per line  |  stderr error + non-zero

# help skips config
bak-files list --help
  -> Usage on stdout (exit 0)

# env
HOME, WORKING_ROLE required; config.validate may require more (e.g. W0)

# mapping (TS-compatible)
"~/Scripts" -> "HOME/Scripts"
"~" -> "HOME/$WORKING_ROLE"
```

## Preconditions

- Module root: `DOCTEST_ROOT/../..` (feature root is `tests/cli-list`).
- Production entrypoint: `cmd/bak-files` (list body may still be stub until
  implementer greens this tree).
- Session cache: `$TMPDIR/bak-files-cli-list-<DOCTEST_SESSION_ID>/`
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — flock one-time build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir` for config
  fixtures; no writes under the module root.
- Reference behavior (private TS): list prints `mappingPath` for each resolved
  item; optional NAMES filter on `mappingPath`; `--config` selects the file
  (default `bak.config.json`).

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Grouping/leaf Setup writes fixtures, sets controlled `Env` and `WorkDir`,
   and completes `Args` (`list`, flags, names).
3. `Run` builds the binary once per session and executes it with `req.Env` and
   `req.WorkDir`.
4. Leaf Assert checks exit code and stdout/stderr (v3 output templates for
   successful path lists).

## Context

- Classic TDD: leaves are RED until list + config load are implemented.
- Sealed `tests/cli-foundation` must not be edited by this phase; foundation
  `stub/list` may need a later reconcile when list becomes real (out of scope
  for the designer).
- Prefer **stable** list order: same order as `files` object keys / resolved
  items in the fixture (not arbitrary map iteration).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}

// writeFile writes path with content (creates parent dirs).
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// minimalEnv returns an environ with PATH (and optional extras). Does not set HOME.
func minimalEnv(extra ...string) []string {
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"TMPDIR=" + os.TempDir(),
	}
	return append(env, extra...)
}

// listEnv sets HOME and WORKING_ROLE (and optional KEY=value extras) for happy paths.
func listEnv(home, role string, extra ...string) []string {
	base := minimalEnv(
		"HOME="+home,
		"WORKING_ROLE="+role,
	)
	return append(base, extra...)
}

// standardListFixture is a small bak.config with two file entries and mapping
// that expands ~ and $WORKING_ROLE into deterministic mapping paths.
//
// With HOME=<home> and WORKING_ROLE=alice expected mapping paths (files order):
//
//	HOME/alice/.bashrc
//	HOME/Scripts/tool.sh
func standardListFixture() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/.bashrc": true,
    "~/Scripts/tool.sh": true
  },
  "targetDir": "./files",
  "mapping": {
    "~/Scripts": "HOME/Scripts",
    "~": "HOME/$WORKING_ROLE"
  },
  "global": {
    "excludes": [".DS_Store"]
  }
}
`
}
```
