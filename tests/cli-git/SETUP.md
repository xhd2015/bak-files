# Scenario

**Feature**: bak-files P5 — `gitTree` / `gitPatch` / `bak.stats` with local-only
git fixtures; dry-run never writes patch or stats

```
# gitTree backup
operator -> bak-files backup --config …
  -> (dirty) copy worktree files under targetDir; write bak.stats {commitHash, hasChange}
  -> (clean) INFO skip; no targetDir / stats churn
  -> (--dry-run) log would…; no targetDir writes; bak.stats unchanged

# gitPatch backup
operator -> bak-files backup --config …
  -> (dirty, real) write .patch under mapping from git diff HEAD
  -> (--dry-run) no patch file

# restore (minimal)
operator -> bak-files restore --config …
  -> plain copy targetDir → source path after prior gitTree backup seed
```

## Preconditions

- Module root: `d.DOCTEST_ROOT/../..` (feature root is `tests/cli-git`).
- Production entrypoint: `cmd/bak-files` (git modes RED until implementer).
- **`git` on PATH** — leaves fail Setup if `git` is missing (not skipped).
- Process-local binary/cache via in-memory mutex (one-process suite; not in-memory mutex)
  - `bak-files` — built binary
  - `binaries.ready`, `build.lock` — in-memory once build
- Per-leaf isolation: each leaf uses `t.TempDir()` as `WorkDir`.
- Git fixtures: `git init` under `HOME/repo` **without** remotes; identity
  set via `user.email` / `user.name` locally. No `fetch`/`push`/`pull`.
- `bak.stats` path: **`WorkDir/bak.stats`** (process cwd).

## Steps

1. Root Setup ensures `Args` is non-nil.
2. Leaf Setup builds WorkDir, local git repo (clean or dirty), config, Env, Args.
3. `Run` builds the binary once per session and executes with `req.Env` / `WorkDir`.
4. Leaf Assert checks exit code, logs (optional), targetDir, bak.stats, patches.

## Context

- Classic TDD: this tree is RED until gitTree/gitPatch/bak.stats land.
- Sealed: do not edit `tests/cli-foundation`, `tests/cli-list`, `tests/cli-backup`,
  `tests/cli-hooks`.
- Mapping: `~` → `HOME`; repo at `~/repo` → backup under `files/HOME/repo/…`.
- MVP gitPatch base: **`HEAD`** (requirement: prefer single file from
  `git diff HEAD`). Filename may be the mapping path ending in `.patch` or a
  file under the mapping directory matching `*.patch`.

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
		// Ignore .git objects for operator source fingerprints when needed;
		// for targetDir/stats fingerprints walk everything.
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

// requireGit fatals if git is not on PATH.
func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Fatalf("git required on PATH for cli-git fixtures: %v", err)
	}
}

// gitRun runs git -C dir args... and fatals on non-zero.
func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	// Avoid global hooks / templates surprises.
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=bak-files-test",
		"GIT_AUTHOR_EMAIL=bak-files-test@example.com",
		"GIT_COMMITTER_NAME=bak-files-test",
		"GIT_COMMITTER_EMAIL=bak-files-test@example.com",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v in %s: %v\nstderr: %s", args, dir, err, stderr.String())
	}
	return strings.TrimSpace(stdout.String())
}

// initLocalRepo creates a local-only git repo at repoDir with one committed file
// README.md = committedContent. No remotes. Returns full HEAD commit hash.
func initLocalRepo(t *testing.T, repoDir, committedContent string) (commitHash string) {
	t.Helper()
	requireGit(t)
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	gitRun(t, repoDir, "init")
	// Default branch name for portability.
	gitRun(t, repoDir, "checkout", "-b", "main")
	gitRun(t, repoDir, "config", "user.email", "bak-files-test@example.com")
	gitRun(t, repoDir, "config", "user.name", "bak-files-test")
	writeFile(t, filepath.Join(repoDir, "README.md"), committedContent)
	gitRun(t, repoDir, "add", "README.md")
	gitRun(t, repoDir, "commit", "-m", "initial")
	return gitRun(t, repoDir, "rev-parse", "HEAD")
}

// dirtyWorktree appends/modifies README.md after the initial commit.
func dirtyWorktree(t *testing.T, repoDir, newContent string) {
	t.Helper()
	writeFile(t, filepath.Join(repoDir, "README.md"), newContent)
}

// baseConfigSkeleton returns common validate/targetDir/mapping/global JSON.
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

// initGitWorld creates WorkDir + HOME, sets Env/TargetDir/WorkDir/StatsPath.
func initGitWorld(t *testing.T, req *Request) (workDir, home string) {
	t.Helper()
	workDir = t.TempDir()
	home = filepath.Join(workDir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	req.WorkDir = workDir
	req.Env = bakEnv(home, "alice")
	req.TargetDir = filepath.Join(workDir, "files")
	req.StatsPath = filepath.Join(workDir, "bak.stats")
	req.MappingKey = "HOME/repo"
	return workDir, home
}

// setupGitTreeBackup configures ~/repo with gitTree:true, optional dirty content.
// committed = content at commit; workContent = current worktree file (may equal
// committed for clean). dryRun appends --dry-run.
func setupGitTreeBackup(t *testing.T, req *Request, committed, workContent string, dryRun bool) {
	t.Helper()
	workDir, home := initGitWorld(t, req)
	repo := filepath.Join(home, "repo")
	hash := initLocalRepo(t, repo, committed)
	if workContent != committed {
		dirtyWorktree(t, repo, workContent)
	}
	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/repo": map[string]any{
			"gitTree": true,
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = repo
	req.CommitHash = hash
	req.Content = workContent
	// Expected backed file for dirty copy (MVP: README.md under mapping).
	req.BackupPath = filepath.Join(req.TargetDir, "HOME", "repo", "README.md")
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	req.StatsBefore = readFileOrEmpty(req.StatsPath)

	args := []string{"backup", "--config", cfgPath}
	if dryRun {
		args = append(args, "--dry-run")
	}
	req.Args = args
}

// setupGitPatchBackup configures ~/repo with gitPatch:true and dirty worktree.
func setupGitPatchBackup(t *testing.T, req *Request, committed, workContent string, dryRun bool) {
	t.Helper()
	workDir, home := initGitWorld(t, req)
	repo := filepath.Join(home, "repo")
	hash := initLocalRepo(t, repo, committed)
	if workContent != committed {
		dirtyWorktree(t, repo, workContent)
	}
	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/repo": map[string]any{
			"gitPatch": true,
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = repo
	req.CommitHash = hash
	req.Content = workContent
	// Patch lives under mapping directory (or as a .patch file path).
	req.PatchDir = filepath.Join(req.TargetDir, "HOME", "repo")
	req.BackupPath = filepath.Join(req.PatchDir, "worktree.patch")
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	req.StatsBefore = readFileOrEmpty(req.StatsPath)

	args := []string{"backup", "--config", cfgPath}
	if dryRun {
		args = append(args, "--dry-run")
	}
	req.Args = args
}

// findPatchFiles returns paths of *.patch files under dir (non-recursive + recursive).
func findPatchFiles(t *testing.T, dir string) []string {
	t.Helper()
	if !pathExists(dir) {
		return nil
	}
	var out []string
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".patch") {
			out = append(out, p)
		}
		return nil
	})
	return out
}

// looksLikeDiff reports whether content resembles a unified / git diff.
func looksLikeDiff(s string) bool {
	if strings.Contains(s, "diff --git") {
		return true
	}
	if strings.Contains(s, "\n+++ ") || strings.Contains(s, "\n--- ") {
		return true
	}
	if strings.Contains(s, "@@") && (strings.Contains(s, "+") || strings.Contains(s, "-")) {
		return true
	}
	return false
}

// statsHasChangeTrue returns true if bak.stats JSON has mappingKey with hasChange true.
// Accepts either nested {"status":{"hasChange":true,"commitHash":"…"}} or flat fields.
func statsHasChangeTrue(statsJSON, mappingKey, commitHash string) error {
	if strings.TrimSpace(statsJSON) == "" {
		return fmt.Errorf("bak.stats empty or missing")
	}
	var root map[string]any
	if err := json.Unmarshal([]byte(statsJSON), &root); err != nil {
		return fmt.Errorf("bak.stats not JSON: %w", err)
	}
	entry, ok := root[mappingKey]
	if !ok {
		// Also try with leading ./ or slash variants — fail with keys listed.
		keys := make([]string, 0, len(root))
		for k := range root {
			keys = append(keys, k)
		}
		return fmt.Errorf("bak.stats missing key %q; have %v", mappingKey, keys)
	}
	m, ok := entry.(map[string]any)
	if !ok {
		return fmt.Errorf("bak.stats[%s] not object", mappingKey)
	}
	// Nested status object (TS shape) or flat.
	status := m
	if st, ok := m["status"].(map[string]any); ok {
		status = st
	}
	hc, ok := status["hasChange"].(bool)
	if !ok || !hc {
		return fmt.Errorf("bak.stats[%s] hasChange want true, got %v", mappingKey, status["hasChange"])
	}
	ch, _ := status["commitHash"].(string)
	if ch == "" {
		return fmt.Errorf("bak.stats[%s] missing commitHash", mappingKey)
	}
	// Accept full hash or prefix match.
	if commitHash != "" && !strings.HasPrefix(commitHash, ch) && !strings.HasPrefix(ch, commitHash) && ch != commitHash {
		// Allow if commitHash equals either way after trim; soft check prefix 7.
		if len(ch) >= 7 && len(commitHash) >= 7 {
			if ch[:7] != commitHash[:7] {
				return fmt.Errorf("bak.stats commitHash %q does not match HEAD %q", ch, commitHash)
			}
		}
	}
	return nil
}
```
