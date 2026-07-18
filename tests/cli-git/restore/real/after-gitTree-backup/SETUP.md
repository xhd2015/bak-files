# Scenario

**Feature**: restore plain content after a gitTree backup layout under targetDir

```
# operator
config: "~/repo": { "gitTree": true }  (restore uses same mapping; engine plain-copy)
seed files/HOME/repo/README.md = "from-backup\n"
source ~/repo/README.md missing or different
operator -> bak-files restore --config bak.config.json
  -> ~/repo/README.md == "from-backup\n"
  -> exit 0
```

## Preconditions

- targetDir already contains mapping layout as if a prior gitTree backup ran.
- Source directory exists (may be a non-git dir for restore MVP — still use a
  local git repo path shape for consistency, or plain dir under HOME/repo).
- Config still lists `gitTree: true` so entry resolves the same mapping; P5
  restore does **not** require FF/checkout.

## Steps

1. initGitWorld; create home/repo; seed BackupPath content; write config.
2. Args: `restore --config …` without dry-run.
3. Content = seeded backup body; SourcePath = home/repo/README.md.

## Context

- P5: “restore after gitTree backup” plain content.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, home := initGitWorld(t, req)
	repo := filepath.Join(home, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	const body = "from-backup\n"
	// Seed backup store (as after gitTree dirty backup).
	backupFile := filepath.Join(req.TargetDir, "HOME", "repo", "README.md")
	writeFile(t, backupFile, body)

	// Destination empty/missing — restore should create it.
	srcFile := filepath.Join(repo, "README.md")

	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/repo": map[string]any{
			"gitTree": true,
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = srcFile
	req.BackupPath = backupFile
	req.Content = body
	req.Args = []string{"restore", "--config", cfgPath}
	return nil
}
```
