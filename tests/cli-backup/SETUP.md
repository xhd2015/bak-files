# Scenario

**Feature**: bak-files P3 — simple backup/restore with `--dry-run` and excludes

```
# backup: filesystem → targetDir
operator -> bak-files backup [--config FILE] [--dry-run]
  -> load config + validate env + resolve mapping
  -> (real) copy sources under targetDir  |  (dry-run) log would…, no writes
  -> missing source: skip info, exit 0

# restore: targetDir → filesystem
operator -> bak-files restore [--config FILE] [--dry-run]
  -> (real) write destinations  |  (dry-run) no dest writes

# help
bak-files backup|restore --help -> Usage on stdout (exit 0)

# excludes
global.excludes + entry excludes → matching names not copied
```

## Preconditions

- Module root: `DOCTEST_ROOT/../..` (feature root is `tests/cli-backup`).
- Production entrypoint: `cmd/bak-files` (backup/restore may still be stubs
  until implementer greens this tree — classic TDD RED first).
- Session cache: `$TMPDIR/bak-files-cli-backup-<DOCTEST_SESSION_ID>/`
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — flock one-time build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`; no writes
  under the module root.
- Fixtures: config + source trees written in leaf Setup; absolute paths recorded
  on `Request` for Assert side effects.

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Grouping/leaf Setup builds temp WorkDir, config, source/seed trees, Env, Args.
3. `Run` builds the binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code, optional log tokens, and filesystem trees.

## Context

- Classic TDD: leaves are RED until backup/restore (+ dry-run, excludes) land.
- Sealed `tests/cli-list` and `tests/cli-foundation` must not be edited here
  (foundation `stub/backup` / `stub/restore` may be reconciled later by implementer).
- Prefer mapping path `HOME/…` under `targetDir` `./files` so backup layout is
  deterministic: `WorkDir/files/HOME/notes.txt`.

```go
import (
	"os"
	"path/filepath"
	"strings"
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

// readFileOrEmpty returns file contents or "" if missing.
func readFileOrEmpty(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// pathExists reports whether path exists (file or dir).
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// treeFingerprint is a stable multi-line listing of all files under root
// with contents (relative paths). Empty string if root does not exist.
func treeFingerprint(t *testing.T, root string) string {
	t.Helper()
	if !pathExists(root) {
		return ""
	}
	var lines []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		lines = append(lines, rel+"\t"+string(b))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return strings.Join(lines, "\n")
}

// minimalEnv returns PATH + TMPDIR (+ optional KEY=value).
func minimalEnv(extra ...string) []string {
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"TMPDIR=" + os.TempDir(),
	}
	return append(env, extra...)
}

// bakEnv sets HOME and WORKING_ROLE for happy-path backup/restore.
func bakEnv(home, role string, extra ...string) []string {
	base := minimalEnv(
		"HOME="+home,
		"WORKING_ROLE="+role,
	)
	return append(base, extra...)
}

// simpleFileConfig is a bak.config with one file entry ~/notes.txt,
// mapping ~ → HOME, targetDir ./files.
// Expected backup path under WorkDir: files/HOME/notes.txt
func simpleFileConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/notes.txt": true
  },
  "targetDir": "./files",
  "mapping": {
    "~": "HOME"
  },
  "global": {
    "excludes": [".DS_Store"]
  }
}
`
}

// dirWithExcludesConfig backs up ~/proj directory; excludes *.tmp via global.
func dirWithExcludesConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/proj": true
  },
  "targetDir": "./files",
  "mapping": {
    "~": "HOME"
  },
  "global": {
    "excludes": ["*.tmp", ".DS_Store"]
  }
}
`
}

// missingSourceConfig lists a path that Setup will not create.
func missingSourceConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/does-not-exist.txt": true
  },
  "targetDir": "./files",
  "mapping": {
    "~": "HOME"
  },
  "global": {
    "excludes": []
  }
}
`
}

// setupSimpleBackupWorld creates WorkDir with config + home/notes.txt content.
// Sets SourcePath, BackupPath, TargetDir, Content, WorkDir, Env.
// command is "backup" or "restore"; dryRun appends --dry-run.
func setupSimpleBackupWorld(t *testing.T, req *Request, command string, dryRun bool, content string) {
	t.Helper()
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	src := filepath.Join(home, "notes.txt")
	cfgPath := filepath.Join(dir, "bak.config.json")
	target := filepath.Join(dir, "files")
	backup := filepath.Join(target, "HOME", "notes.txt")

	writeFile(t, cfgPath, simpleFileConfig())
	if content != "" && command == "backup" {
		writeFile(t, src, content)
	}
	// For restore real/dry-run, leaf may seed backup and optionally dest.

	req.WorkDir = dir
	req.Env = bakEnv(home, "alice")
	req.TargetDir = target
	req.SourcePath = src
	req.BackupPath = backup
	req.Content = content

	args := []string{command, "--config", cfgPath}
	if dryRun {
		args = append(args, "--dry-run")
	}
	req.Args = args
}
```
