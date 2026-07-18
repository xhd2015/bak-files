# Scenario

**Feature**: backup of content already present under targetDir leaves bytes unchanged

```
# operator
source home/notes.txt = "same-bytes\n"
pre-seed files/HOME/notes.txt = "same-bytes\n"   # prior backup
operator -> bak-files backup --config …
  -> files/HOME/notes.txt still "same-bytes\n"
  -> exit 0
  (optional: log skip/identical; not required)
```

## Preconditions

- Source and pre-seeded backup file have identical content.
- Real backup (no dry-run).

## Steps

1. Create simple-file world with content `same-bytes\n`.
2. Pre-write BackupPath with the same bytes (simulates a prior backup).
3. Snapshot TargetDir fingerprint into `TargetDirBefore`.
4. Leaf Run executes `backup` again.

## Context

- Optional P3 criterion: assert **content** (and preferably whole target tree
  fingerprint) unchanged; implementer may skip rewrite when identical.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const body = "same-bytes\n"
	setupSimpleBackupWorld(t, req, "backup", false, body)
	// Prior backup already on disk with identical content.
	writeFile(t, req.BackupPath, body)
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	return nil
}
```
