# Scenario

**Feature**: any `files` key `PREFIX` with a **string array** value is a basename
whitelist under that PREFIX — never a full-PREFIX copy job

```
# general array whitelist (non-home PREFIX)
files: { "$W0/proj": [".vscode"] }   # or any PREFIX: [basename…]
mapping: { "$W0": "W" }
operator -> bak-files backup --config …
  -> one job per name: Source=expand(PREFIX)/name, mapping PREFIX/name
  -> NEVER Job{Source: expand(PREFIX)} / copyDir(full project tree)
  -> vendor / poison trees under PREFIX stay out of targetDir
```

## Preconditions

- Non-home PREFIX required (reproduce production `$W0/…` config).
- Env: `W0=<workdir>/w0`, plus HOME/WORKING_ROLE for validate builtins.
- Poison sibling under PREFIX (e.g. `vendor/poison.go`) proves no full-PREFIX walk.
- `includeDotFiles: false` so home discovery does not confuse store asserts.

## Steps

1. Leaves use `projectArrayConfig` / `setupProjectArrayWorld`.
2. Real backup (not dry-run) so filesystem asserts prove no vendor poison.

## Context

- Classic TDD RED until `ResolveJobs` expands **any** string-array files value
  the same way bare `"~": [names]` already expands (not only `key == "~"`).
- Today non-`~` array values still schedule Source=full PREFIX → vendor poison.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	} else if req.Args[0] != "backup" {
		req.Args = append([]string{"backup"}, req.Args...)
	}
	return nil
}
```
