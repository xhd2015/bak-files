# Scenario

**Feature**: bare `files` key `"~"` is a **home basename whitelist**, never a full-`$HOME` job

```
# home whitelist
files: { "~": [names…] }  (+ optional explicit ~/path keys)
operator -> bak-files backup --config …
  -> each name → job source $HOME/name (mapping HOME/$ROLE/name)
  -> NEVER Job{Source: $HOME} / copyDir(full home)
  -> non-whitelisted non-dot trees (Library, Downloads, …) stay out of targetDir
```

## Preconditions

- Simulated HOME under WorkDir; mapping `~` → `HOME/$WORKING_ROLE`, role `alice`.
- Leaves plant both desired whitelist paths and poison non-dot trees (e.g. `Library`).

## Steps

1. Leaves use `tildeArrayConfig` / `tildeArrayConfigDotsOff` / `explicitScriptsConfig`.
2. Real backup (not dry-run) so filesystem asserts prove no full-home copy.

## Context

- Classic TDD RED until `ResolveJobs` expands `"~": [names]` into per-name jobs.
- Today the engine expands key `"~"` to `$HOME` and ignores the array → full-home walk.

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
