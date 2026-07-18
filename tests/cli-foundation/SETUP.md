# Scenario

**Feature**: bak-files P1 foundation — multi-level help, stub subcommands, build, .gitignore

```
# operator invokes bak-files binary; CLI dispatches help / stub / unknown
operator -> bak-files [args] -> stdout / stderr + exit code

# root help
bak-files | bak-files --help | bak-files -h
  -> Usage listing backup, restore, list (exit 0)

# subcommand help
bak-files backup|restore|list --help
  -> command Usage (exit 0)

# stub run (not yet implemented)
bak-files backup|restore|list
  -> not implemented message (exit non-zero)

# unknown
bak-files <garbage>
  -> Error: or bak-files: on stderr (exit non-zero)

# build
go build -o <bin> ./cmd/bak-files -> binary (exit 0)

# safety
repo .gitignore covers binaries, files/, bak.stats, sum.index, .env, .DS_Store
```

## Preconditions

- Module root: `DOCTEST_ROOT/../..` (repo root with `go.mod`; feature root is
  `tests/cli-foundation`).
- Production entrypoint (implementer): `cmd/bak-files` (may not exist yet — leaves
  that need the binary will fail `go build` until implemented).
- Flags library (implementer): prefer `github.com/xhd2015/less-flags`.
- Session cache (shared across parallel leaves of one `doctest test` run):
  `$TMPDIR/bak-files-cli-foundation-<DOCTEST_SESSION_ID>/`
  - `bak-files` — built binary
  - `binaries.ready` — sentinel after successful build
  - `build.lock` — flock for first-time population
- Per-leaf isolation: CLI leaves only read stdout/stderr/exit; no shared mutable
  workspace under the module root.
- Default `Request.Mode` is `"cli"`. Leaves set `Mode` to `"build"` or
  `"gitignore"` when not invoking the binary.

## Steps

1. Root Setup leaves `Mode`/`Args` defaults unset for grouping nodes to fill.
2. For `Mode=cli`, `Run` builds the binary once per session (flock) then
   executes `bin Args...` with cwd = module root; captures stdout, stderr, exit.
3. For `Mode=build`, `Run` runs a fresh `go build -o <temp> ./cmd/bak-files`.
4. For `Mode=gitignore`, `Run` reads module-root `.gitignore` into
   `Response.GitignoreContent`.
5. Leaf Assert checks exit codes, usage tokens, error prefixes, or patterns.

## Context

- Classic TDD greenfield: only `go.mod` exists initially; implementer adds
  `cmd/bak-files` and `.gitignore` to make leaves pass.
- Help must not require config files or environment variables.
- Stub subcommands prefer **non-zero** exit without `--help` so later phases
  can replace stubs without flipping success semantics for help.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	if req.Mode == "" {
		req.Mode = "cli"
	}
	return nil
}
```
