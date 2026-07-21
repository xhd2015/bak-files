# Scenario

**Feature**: bak-files preserves OS symlinks on backup (self-dir + file links)

```
# backup real
operator -> bak-files backup --config …
  -> walk ~/Scripts; ModeSymlink → copySymlink (same target string)
  -> self symlink-to-dir: exit 0; no "is a directory"

# backup dry-run
operator -> bak-files backup --config … --dry-run
  -> log "would symlink"; no targetDir writes
```

## Preconditions

- Module root: `d.DOCTEST_ROOT/../..` (feature root is `tests/cli-symlink`).
- Production entrypoint: `cmd/bak-files` (symlink preserve already implemented).
- Process-local binary/cache via in-memory mutex (one-process suite; not in-memory mutex)
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — in-memory once build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`; no writes
  under the module root.
- Fixtures: mapping `"~": "HOME/$WORKING_ROLE"`, `WORKING_ROLE=alice`,
  `targetDir` `./files`, `includeDotFiles: false`, explicit `~/Scripts` only.
  Layout: `WorkDir/files/HOME/alice/Scripts/…`.

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Grouping/leaf Setup builds temp WorkDir, config, Scripts tree + symlinks,
   Env, Args.
3. `Run` builds the binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code, Readlink targets, regular file content, logs.

## Context

- **Coverage backfill:** GREEN expected for covered paths.
- Sealed sibling trees must not be edited here.
- Prefer Lstat/Readlink helpers — do not follow self-dir links with ReadFile
  blindly when fingerprinting trees that contain symlink-to-dir entries.

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

// pathExists reports whether path exists (file, dir, or symlink; follows for Stat).
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// pathLexists reports whether path exists via Lstat (true for dangling links too).
func pathLexists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// isSymlink reports whether path is a symlink.
func isSymlink(path string) bool {
	st, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeSymlink != 0
}

// readlinkOrEmpty returns Readlink result or "".
func readlinkOrEmpty(path string) string {
	s, err := os.Readlink(path)
	if err != nil {
		return ""
	}
	return s
}

// symlink creates a symlink (removes existing path first if present).
func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("mkdir for symlink: %v", err)
	}
	_ = os.RemoveAll(link)
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink %s -> %s: %v", link, target, err)
	}
}

// treeFingerprint is a stable multi-line listing under root.
// Regular files: rel + tab + body. Symlinks: rel + " -> " + target.
// Does not follow symlink-to-dir (Lstat; only real dirs recurse).
func treeFingerprint(t *testing.T, root string) string {
	t.Helper()
	if _, err := os.Lstat(root); err != nil {
		return ""
	}
	var lines []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			tgt, err := os.Readlink(p)
			if err != nil {
				return err
			}
			lines = append(lines, rel+" -> "+tgt)
			return nil
		}
		if info.IsDir() {
			return nil
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

// scriptsConfig backs up ~/Scripts only; dots off; mapping ~ → HOME/$WORKING_ROLE.
func scriptsConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/Scripts": true
  },
  "targetDir": "./files",
  "mapping": {
    "~": "HOME/$WORKING_ROLE"
  },
  "global": {
    "includeDotFiles": false,
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

// setupScriptsWorld creates WorkDir, home, config, Env, TargetDir, HomeDir.
// command is typically "backup"; dryRun appends --dry-run.
// extraArgs are appended after --config path (e.g. --verbose).
func setupScriptsWorld(t *testing.T, req *Request, command string, dryRun bool, extraArgs ...string) (dir, home, cfgPath string) {
	t.Helper()
	dir = t.TempDir()
	home = filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(home, "Scripts"), 0o755); err != nil {
		t.Fatalf("mkdir Scripts: %v", err)
	}
	cfgPath = filepath.Join(dir, "bak.config.json")
	writeFile(t, cfgPath, scriptsConfig())
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
```
