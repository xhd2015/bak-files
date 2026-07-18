# Scenario

**Feature**: global excludes prevent matching names from being backed up

```
# operator
HOME=…/home
write home/proj/keep.txt = "keep-me\n"
write home/proj/noise.tmp = "noise\n"
config: files["~/proj"]=true, global.excludes=["*.tmp"]
operator -> bak-files backup --config …
  -> files/HOME/proj/keep.txt exists
  -> files/HOME/proj/noise.tmp MUST NOT exist
  -> exit 0
```

## Preconditions

- Directory entry `~/proj` with one kept and one excluded file.
- Exclude pattern `*.tmp` in `global.excludes`.

## Steps

1. Write `dirWithExcludesConfig()` and source tree under home/proj.
2. Args: `backup --config <path>`.
3. Record KeepBackupPath and ExcludedBackupPath.

## Context

- Entry-level excludes would be equivalent; global is the shared rule for P3.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	proj := filepath.Join(home, "proj")
	writeFile(t, filepath.Join(proj, "keep.txt"), "keep-me\n")
	writeFile(t, filepath.Join(proj, "noise.tmp"), "noise\n")
	cfgPath := filepath.Join(dir, "bak.config.json")
	writeFile(t, cfgPath, dirWithExcludesConfig())

	target := filepath.Join(dir, "files")
	req.WorkDir = dir
	req.Env = bakEnv(home, "alice")
	req.TargetDir = target
	req.KeepBackupPath = filepath.Join(target, "HOME", "proj", "keep.txt")
	req.ExcludedBackupPath = filepath.Join(target, "HOME", "proj", "noise.tmp")
	req.Content = "keep-me\n"
	req.Args = []string{"backup", "--config", cfgPath}
	return nil
}
```
