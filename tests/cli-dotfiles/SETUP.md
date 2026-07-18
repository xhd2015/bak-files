# Scenario

**Feature**: bak-files home whitelist (`"~": [names]`), **general files array
whitelist** (any `PREFIX: [basenames]`), default home dots (`includeDotFiles` ≡
`~/.*`), pathflag skips, CLI/config overrides

```
# whitelist + discovery + walk skip policy
operator -> bak-files backup|restore|list [--config FILE] [--dry-run]
           [--no-dot-files] [--include PATH]… [--exclude PATH]…
  -> files key "~": [names] → each $HOME/name job (NEVER full-home Source=$HOME)
  -> files key PREFIX: [names] → each expand(PREFIX)/name job (NEVER full PREFIX)
  -> discover $HOME top-level dots when includeDotFiles (default true) ≡ ~/.*
  -> walk under included jobs: force-include → force-exclude → pathflag → basename excludes
  -> (real) copy keep paths  |  (dry-run) log would… / skip…, no writes

# help
bak-files backup --help -> Usage documents --no-dot-files, --include, --exclude
```

## Preconditions

- Module root: `DOCTEST_ROOT/../..` (feature root is `tests/cli-dotfiles`).
- Production entrypoint: `cmd/bak-files` (home-whitelist ResolveJobs may be RED
  until implementer greens — classic TDD).
- Session cache: `$TMPDIR/bak-files-cli-dotfiles-<DOCTEST_SESSION_ID>/`
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — flock one-time build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`; no writes
  under the module root.
- Fixtures: prefer mapping `"~": "HOME/$WORKING_ROLE"`, `WORKING_ROLE=alice`,
  `targetDir` `./files` so layout is `WorkDir/files/HOME/alice/…`.

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Grouping/leaf Setup builds temp WorkDir, config, home trees, Env, Args.
3. `Run` builds the binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code, mapping/list output, logs, and filesystem trees.

## Context

- Classic TDD: home-whitelist leaves require `"~": [names]` expand (no full-home).
  New **files-array-whitelist** leaves are RED until ResolveJobs expands **any**
  string-array files value (e.g. `"$W0/proj": [".vscode"]`) into per-name jobs
  (no full-PREFIX copyDir) — not only bare `"~"`.
- Sealed sibling trees must not be edited here.
- Empty `"files": {}` isolates auto-discovery; `"~": [names]` and PREFIX arrays
  test array expand.

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

// emptyFilesConfigDotsOff sets includeDotFiles false.
func emptyFilesConfigDotsOff() string {
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
    "includeDotFiles": false,
    "excludes": [".DS_Store"]
  }
}
`
}

// emptyFilesConfigBasenameExcludes adds global excludes for *.tmp.
func emptyFilesConfigBasenameExcludes() string {
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
    "excludes": ["*.tmp", ".DS_Store"]
  }
}
`
}

// explicitBashrcConfig lists ~/.bashrc so discovery must not double-schedule it.
func explicitBashrcConfig() string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~/.bashrc": true
  },
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

// tildeArrayConfig: files key "~" is a whitelist of home basenames (never full-home).
// includeDotFiles default true (omitted). namesJSON is a JSON array body, e.g. [".bashrc"].
func tildeArrayConfig(namesJSON string) string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~": ` + namesJSON + `
  },
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

// tildeArrayConfigDotsOff is tildeArrayConfig with global.includeDotFiles false.
func tildeArrayConfigDotsOff(namesJSON string) string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE"]
    }
  ],
  "files": {
    "~": ` + namesJSON + `
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

// explicitScriptsConfig lists only ~/Scripts (whitelist); dots default on.
func explicitScriptsConfig() string {
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

// projectArrayConfig: non-home PREFIX "$W0/proj" is a basename whitelist (never full PREFIX).
// namesJSON is a JSON array body, e.g. [".vscode"]. Dots off for store isolation.
func projectArrayConfig(namesJSON string) string {
	return `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE", "W0"]
    }
  ],
  "files": {
    "$W0/proj": ` + namesJSON + `
  },
  "targetDir": "./files",
  "mapping": {
    "$W0": "W",
    "~": "HOME/$WORKING_ROLE"
  },
  "global": {
    "includeDotFiles": false,
    "excludes": [".DS_Store"]
  }
}
`
}

// setupDotsWorld creates WorkDir, home, config, Env, TargetDir, HomeDir.
// command is backup|restore|list; dryRun appends --dry-run for backup/restore.
// extraArgs are appended after --config path (e.g. --no-dot-files).
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

// setupProjectArrayWorld creates WorkDir, W0 root, home, config, Env with W0 set.
// Returns dir, w0 absolute path, cfgPath. TargetDir = dir/files.
func setupProjectArrayWorld(t *testing.T, req *Request, cfgJSON, command string, dryRun bool, extraArgs ...string) (dir, w0, cfgPath string) {
	t.Helper()
	dir = t.TempDir()
	home := filepath.Join(dir, "home")
	w0 = filepath.Join(dir, "w0")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	if err := os.MkdirAll(w0, 0o755); err != nil {
		t.Fatalf("mkdir w0: %v", err)
	}
	cfgPath = filepath.Join(dir, "bak.config.json")
	writeFile(t, cfgPath, cfgJSON)
	target := filepath.Join(dir, "files")

	req.WorkDir = dir
	req.HomeDir = home
	req.Env = bakEnv(home, "alice", "W0="+w0)
	req.TargetDir = target

	args := []string{command, "--config", cfgPath}
	if dryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, extraArgs...)
	req.Args = args
	return dir, w0, cfgPath
}
```
