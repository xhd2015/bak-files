# Scenario

**Feature**: bak-files access-denied → warning; pathflag trash catalog skip

```
# catalog skip at discovery
HOME has .Trash + .bashrc; files={}
operator -> bak-files backup --config …
  -> INFO: skip .Trash (macOS trash); .bashrc copied; no .Trash under targetDir

# walk access deny
HOME has .private/ chmod 000 (+ nested file) + good dots
operator -> bak-files backup --config …
  -> warning: skip <path>: …; exit 0; good files still under mapping
```

## Preconditions

- Module root: `d.DOCTEST_ROOT/../..` (feature root is `tests/cli-access-warn`).
- Production entrypoint: `cmd/bak-files` (warn + trash skip already implemented).
- Process-local binary/cache via in-memory mutex (one-process suite; not in-memory mutex)
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — in-memory once build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`; no writes
  under the module root.
- Fixtures: mapping `"~": "HOME/$WORKING_ROLE"`, `WORKING_ROLE=alice`,
  `targetDir` `./files`, empty `files`, default includeDotFiles.
  Layout: `WorkDir/files/HOME/alice/…`.
- chmod 000 leaves must restore mode via `t.Cleanup` so TempDir cleanup works.

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Grouping/leaf Setup builds temp WorkDir, config, home trees, Env, Args.
3. `Run` builds the binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code, logs (`INFO: skip` / `warning: skip`), and trees.

## Context

- **Coverage backfill:** GREEN expected for covered paths.
- Sealed sibling trees must not be edited here.
- chmod 000 is a portable stand-in for macOS TCC on `.Trash` open failure;
  pathflag catalog skip of `.Trash` is a separate (preferred) discovery path.

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

// pathLexists reports whether path exists via Lstat.
func pathLexists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// treeFingerprint is a stable multi-line listing of all regular files under root
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

// bakEnv sets HOME and WORKING_ROLE for happy-path invocations.
func bakEnv(home, role string, extra ...string) []string {
	base := minimalEnv(
		"HOME="+home,
		"WORKING_ROLE="+role,
	)
	return append(base, extra...)
}

// emptyFilesConfig: no explicit files; dots default-on; mapping ~ → HOME/$WORKING_ROLE.
func emptyFilesConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {},
  "targetDir": "./files",
  "mapping": {
    "~": "HOME/$WORKING_ROLE"
  },
  "global": {
    "excludes": [".DS_Store"]
  }
}
`
}

// mappingBackup joins targetDir with HOME/<role>/rel parts.
func mappingBackup(target, role string, relParts ...string) string {
	parts := append([]string{target, "HOME", role}, relParts...)
	return filepath.Join(parts...)
}

// setupDotsWorld creates WorkDir, home, config, Env, TargetDir, HomeDir.
// command is typically "backup"; dryRun appends --dry-run.
// extraArgs are appended after --config path.
func setupDotsWorld(t *testing.T, req *Request, cfgJSON, command string, dryRun bool, extraArgs ...string) (dir, home, cfgPath string) {
	t.Helper()
	dir = t.TempDir()
	home = filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	cfgPath = filepath.Join(dir, "bak.config.json")
	writeFile(t, cfgPath, cfgJSON)
	target := filepath.Join(dir, "files")

	req.WorkDir = dir
	req.HomeDir = home
	req.Env = bakEnv(home, "alice")
	req.TargetDir = target

	args := []string{command, "--config", cfgPath}
	if dryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, extraArgs...)
	req.Args = args
	return dir, home, cfgPath
}

// combinedOut lowercases stdout+stderr for token checks.
func combinedOut(resp *Response) string {
	return strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
}

// makeUnreadableDir creates dir (mode 0755), writes nested file, then chmod 000.
// Registers t.Cleanup to restore u+rwx so TempDir can remove the tree.
// Nested dir under a readable parent so Walk hits open errors as a normal user.
func makeUnreadableDir(t *testing.T, dir, nestedName, body string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir unreadable parent: %v", err)
	}
	nested := filepath.Join(dir, nestedName)
	if err := os.WriteFile(nested, []byte(body), 0o644); err != nil {
		t.Fatalf("write nested under unreadable: %v", err)
	}
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatalf("chmod 000 %s: %v", dir, err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0o700)
	})
}

// restoreUnreadable ensures chmod for cleanup mid-assert (also covered by Cleanup).
func restoreUnreadable(dir string) {
	if dir == "" {
		return
	}
	_ = os.Chmod(dir, 0o700)
}
```
