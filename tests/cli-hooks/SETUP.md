# Scenario

**Feature**: bak-files P4 — hooks (`beforeCopy`/`afterCopy`/`beforeRestore`) and
`cmd` / `restoreCmd` process execution; dry-run never execs scripts

```
# backup: cmd generates file; beforeCopy runs around copy
operator -> bak-files backup --config … [--dry-run]
  -> (real) run beforeCopy / cmd / afterCopy; write targetDir
  -> (dry-run) log would…; no shell exec; no marker; no generated file
  -> beforeCopy exit ≠ 0 → bak-files exit ≠ 0

# restore: restoreCmd / beforeRestore
operator -> bak-files restore --config …
  -> (real) execute restoreCmd → side-effect file
```

## Preconditions

- Module root: `d.DOCTEST_ROOT/../..` (feature root is `tests/cli-hooks`).
- Production entrypoint: `cmd/bak-files` (hooks/cmd may be RED until implementer).
- Process-local binary/cache via in-memory mutex (one-process suite; not in-memory mutex)
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — in-memory once build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`; scripts only
  touch paths under that WorkDir.
- Safe scripts only: `printf`/`echo`, `touch`, `exit 1`.

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Leaf Setup builds temp WorkDir, config with hooks/cmd, Env, Args, path fields.
3. `Run` builds binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code and filesystem side effects (or their absence).

## Context

- Classic TDD: this tree is RED until P4 hook/`cmd` support lands.
- Sealed: do not edit `tests/cli-foundation`, `tests/cli-list`, `tests/cli-backup`.
- Mapping fixtures: prefer `~` → `HOME` and logical keys mapped under `HOME/…`
  so paths stay deterministic: `WorkDir/files/HOME/…`.

```go
import (
	"encoding/json"
	"fmt"
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

// shellQuote wraps path for embedding in a single-quoted shell fragment.
// Rejects paths containing single quotes (temp dirs never have them).
func shellQuote(path string) string {
	if strings.Contains(path, "'") {
		panic("path contains single quote: " + path)
	}
	return "'" + path + "'"
}

// baseConfigSkeleton returns common validate/targetDir/mapping/global JSON
// fields as a map for building hook/cmd configs.
func baseConfigSkeleton() map[string]any {
	return map[string]any{
		"validate": []any{
			map[string]any{"env": []any{"HOME", "WORKING_ROLE"}},
		},
		"targetDir": "./files",
		"mapping": map[string]any{
			"~": "HOME",
		},
		"global": map[string]any{
			"excludes": []any{".DS_Store"},
		},
	}
}

// writeConfigJSON marshals cfg to WorkDir/bak.config.json and returns path.
func writeConfigJSON(t *testing.T, workDir string, cfg map[string]any) string {
	t.Helper()
	path := filepath.Join(workDir, "bak.config.json")
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	writeFile(t, path, string(b)+"\n")
	return path
}

// initHookWorld creates WorkDir + home, sets Env/TargetDir/WorkDir.
// Leaves fill files entry, Args, and path fields.
func initHookWorld(t *testing.T, req *Request) (workDir, home string) {
	t.Helper()
	workDir = t.TempDir()
	home = filepath.Join(workDir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	req.WorkDir = workDir
	req.Env = bakEnv(home, "alice")
	req.TargetDir = filepath.Join(workDir, "files")
	return workDir, home
}

// mustFormat is fmt.Sprintf that fatals on empty format misuse (helper name).
func mustFormat(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
```
