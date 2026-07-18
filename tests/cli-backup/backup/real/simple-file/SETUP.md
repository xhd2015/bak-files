# Scenario

**Feature**: backup copies a single source file into `targetDir`

```
# operator
HOME=…/home  WORKING_ROLE=alice
write home/notes.txt = "hello-backup\n"
operator -> bak-files backup --config bak.config.json
  -> files/HOME/notes.txt == "hello-backup\n"  (exit 0)
```

## Preconditions

- Fixture: `simpleFileConfig()`; source `~/notes.txt` exists with known body.
- Env: HOME, WORKING_ROLE=alice.

## Steps

1. Create WorkDir with config and `home/notes.txt`.
2. Args: `backup --config <path>` (no dry-run).
3. Record SourcePath, BackupPath, Content for Assert.

## Context

- Primary exit criterion for P3 backup happy path.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupSimpleBackupWorld(t, req, "backup", false, "hello-backup\n")
	return nil
}
```
