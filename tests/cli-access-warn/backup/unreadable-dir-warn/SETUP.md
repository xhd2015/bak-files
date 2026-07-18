# Scenario

**Feature**: chmod 000 unreadable home dir → warning: skip, exit 0, bashrc ok

```
# operator
HOME=…/home  files={}
write home/.bashrc = "export OK=1\n"
mkdir home/.private
write home/.private/secret = "nope\n"
chmod 000 home/.private
operator -> bak-files backup --config …
  -> exit 0
  -> combined log has "warning:" and "skip" and path (.private)
  -> files/HOME/alice/.bashrc == body
  -> no fatal Error: abort of whole backup
  -> secret not required under target (dir unreadable)
```

## Preconditions

- Empty files; includeDotFiles default true.
- `.private` is **not** pathflag DefaultSkipMask — it is scheduled as a job so
  Walk / open hits permission denied (portable stand-in for TCC).
- Nested file under unreadable dir; parent home stays readable.
- Tests run as normal user; root CI that ignores chmod 000 is out of band.
- `t.Cleanup` restores mode so TempDir removal works.

## Steps

1. setupDotsWorld emptyFilesConfig, real backup.
2. Write `.bashrc`; create `.private/secret`; chmod 000 `.private`.
3. Record BackupPath, UnreadableDir, ExcludedBackupPath for secret mapping.

## Context

- Asserts `logWarn` / `walkAccessErr` path: `warning: skip <path>: …`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false)

	const body = "export OK=1\n"
	writeFile(t, filepath.Join(home, ".bashrc"), body)

	private := filepath.Join(home, ".private")
	makeUnreadableDir(t, private, "secret", "nope\n")

	req.Content = body
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.UnreadableDir = private
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".private", "secret")
	return nil
}
```
