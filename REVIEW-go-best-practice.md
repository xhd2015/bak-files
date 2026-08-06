# Go Best Practice Review — bak-files

**Date:** 2026-08-06  
**Branch / worktree:** `master-2026-08-06-use-go-best-practice-to-review-current-project`  
**Module:** `github.com/xhd2015/bak-files`  
**Scope:** codebase structure, CLI design, flags handling, package layout  
**Reference:** [go-best-practice](https://github.com/xhd2015) skill topics: `cli/*`, `flags-parsing/*`, `cmd-exec`, `kool-create`, `go-embed-assets`  
**Mode:** review only (no implementation)

---

## Executive summary

`bak-files` is a config-driven backup/restore tool with a solid **library core** (`pathflag`, `internal/config`, `internal/engine`) and an extensive **doctest contract** for a multi-command CLI. Against go-best-practice, the largest gap is not style: the **CLI entrypoint is missing entirely** (`cmd/bak-files` never exists in the tree or git history), while README, CI, and ~all `tests/cli-*` trees assume it. `less-flags` is present only as an orphan module dependency.

What already exists under `internal/` is largely well-shaped for a Go CLI product: pure public classifier, internal config/engine, dry-run *intent* in the engine, streaming-friendly log lines. The review below orders findings by severity and maps each recommended change to a concrete skill topic.

---

## Project snapshot

| Area | Current state |
|------|----------------|
| Packages | `pathflag` (public), `internal/config`, `internal/engine` |
| CLI | **Absent** — no `cmd/`, no `main` package |
| Flags lib | `github.com/xhd2015/less-flags v1.0.2` in `go.mod`/`go.sum` as **indirect**, `go mod why` → main module does not need it |
| External cmds | `os/exec` for `git` and `sh -c` (hooks / `cmd` / restore) |
| Embed / UI | None — `go-embed-assets` N/A |
| Scaffold | Hand-built CLI module — `kool-create` N/A as product template |
| Tests | No `*_test.go`; doctest trees under `tests/` (pathflag pure + CLI E2E) |
| CI | `go test ./...` then `doctest test -v --label-all ./...` (doctest will fail without binary) |
| Build today | `go build ./...` succeeds (library packages only) |

Documented layout (README) vs reality:

```text
cmd/bak-files/     ← MISSING
internal/config/   ← present
internal/engine/   ← present (~1.4k LOC single file)
pathflag/          ← present
tests/             ← present (cli-*, pathflag)
```

---

## Findings (by severity)

### Critical

#### C1. Missing CLI entrypoint `cmd/bak-files`

**Evidence:** README install/build paths, `.github/workflows/test.yml` doctest step, and every `tests/cli-*` SETUP/ASSERT expect `go build ./cmd/bak-files` and subcommands `backup` / `restore` / `list`. No `cmd/` directory and no `main.go` in history.

**Why it matters:** The product is unusable as a CLI; doctests cannot green; `go install github.com/xhd2015/bak-files/cmd/bak-files@latest` is broken.

**Topic:** `flags-parsing/subcommand` (dispatcher + per-level help), `flags-parsing` (less-flags wiring), README/layout conventions for Go CLIs.

**Recommended change:**

1. Add `cmd/bak-files/main.go` (thin `main` → `run(os.Args[1:]) error`).
2. Dispatch with the less-flags **subcommand** pattern (`StopOnFirstArg` if root gains global flags, or the “no toplevel flags” variant if root only dispatches).
3. Implement **help at every level** (root, `backup`, `restore`, `list`) — already required by `tests/cli-foundation` and `flags-parsing/subcommand`.
4. Wire handlers to existing packages:
   - `list` → load config → `ValidateEnvs` → `engine.MappingPaths` (not only `config.MappingPaths`; see M3)
   - `backup` / `restore` → `engine.Backup` / `engine.Restore` with `Options` filled from flags
5. Prefer non-zero exit + stderr prefix `Error:` or `bak-files:` for unknown commands (foundation unknown ASSERT).

Sketch (aligned with skill recipe):

```go
// cmd/bak-files/main.go — conceptual
func run(args []string) error {
    if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
        fmt.Print(topHelp)
        return nil
    }
    switch args[0] {
    case "backup":
        return handleBackup(args[1:])
    case "restore":
        return handleRestore(args[1:])
    case "list":
        return handleList(args[1:])
    default:
        return fmt.Errorf("unknown command: %s", args[0])
    }
}
```

---

#### C2. `less-flags` is declared but unused (orphan dependency)

**Evidence:** `go.mod` marks `less-flags` as `// indirect`; `go mod why` reports the main module does not need the package. Zero Go imports of `less-flags`. Doctest foundation SETUP explicitly prefers `github.com/xhd2015/less-flags`.

**Topic:** `flags-parsing`, `flags-parsing/types`, `flags-parsing/subcommand`.

**Recommended change:**

- When adding `cmd/bak-files`, import less-flags and run `go get github.com/xhd2015/less-flags@latest` so it is a **direct** require.
- Do **not** hand-roll flag parsing for production CLI; doctests and skill both assume less-flags.
- Per-command flag sets (from README + doctests):

| Command | Flags (expected) |
|---------|------------------|
| `backup` | `--config`, `--dry-run`, `--no-dot-files`, `--include` (repeat), `--exclude` (repeat), `-v/--verbose`, `-h/--help` |
| `restore` | same family as backup |
| `list` | `--config`, optional name filter args, `-h/--help` (and likely same dot filters if list uses discovery) |

Example handler shape (`flags-parsing/types` + `subcommand`):

```go
var (
    configPath string
    dryRun, verbose, noDotFiles bool
    include, exclude []string
)
remain, err := lessflags.
    String("--config", &configPath).
    Bool("--dry-run", &dryRun).
    Bool("--no-dot-files", &noDotFiles).
    StringSlice("--include", &include).
    StringSlice("--exclude", &exclude).
    Bool("-v,--verbose", &verbose).
    Help("-h,--help", backupHelp).
    Parse(args)
```

Default `--config` to `config.DefaultConfigName` (`bak.config.json`) when empty.

---

### High

#### H1. Dry-run uses **separate pipelines** in places (drift risk)

**Topic:** `cli/dry-run` — *one path, gate side effects*; anti-pattern is `handleDryRun()` that reimplements discovery/plan.

**Evidence in `internal/engine/engine.go`:**

| Path | Pattern | Assessment |
|------|---------|------------|
| `Backup` plain jobs | Early `if opt.DryRun { dryRunBackupJob(...); continue }` then live hooks/copy | **Separate dry-run function** — skill anti-pattern |
| `Restore` | Large inline dry-run block before live path | Parallel control flow (same discovery via `ResolveJobs`, but copy/hooks reimplemented) |
| `backupGitTree` / `backupGitPatch` | Inline gate after dirty check | **Closer to preferred** (same dirty plan; gate writes/hooks) |
| `copyDir` / `copyFile` / `copySymlink` | Inline `if opt.DryRun` gates | Preferred, but **bypassed** by `Backup`’s early branch |

**Concrete drift already visible:**

- Plain dry-run logs `beforeCopy` / `afterCopy` / `cmd` via `dryRunBackupJob`.
- `backupGitTree` dry-run logs only “would copy dirty gitTree” + bak.stats — **no** would-run for `beforeCopy`/`afterCopy` even if set.
- Live `copyDir` preserves symlinks as links; dry-run walk logs `would symlink` without target string (live logs `symlinked dst -> target`). Minor UX drift.
- Live path runs identical-content skip (`sameContent`); dry-run always “would copy” — acceptable for dry-run, but documents that dry-run is not a perfect mirror of apply.

**Recommended change (when implementing / refactoring — not trivial docs):**

1. Prefer **one job pipeline**:
   - `ResolveJobs` (already shared) → for each job: plan steps → gate side effects.
2. Collapse `dryRunBackupJob` / restore dry-run block into the live structure:
   - `if !opt.DryRun { runShell(beforeCopy) } else { log would-run }`
   - Reuse `backupCopySource` / `copyDir` / `copyFile` which already understand `opt.DryRun`.
3. For git modes: keep dirty check **outside** the gate; log would-run hooks consistently with plain backup.
4. Output convention (skill): prefer `[dry-run]` prefix on **stdout** for planned lines; soft warnings on **stderr**. Today engine uses `dry-run: …` (doctests only require substring `dry-run` or `would` — either prefix is OK if tests stay green; aligning to `[dry-run]` is preferred for skill consistency).

---

#### H2. Dual `MappingPaths` APIs invite wrong `list` implementation

**Evidence:**

- `config.MappingPaths()` — expands **files keys only** (no home-dot discovery, no array whitelist expansion to child jobs).
- `engine.MappingPaths(cfg, opt)` — uses `ResolveJobs` (whitelist arrays, home guards, discovery, skip filters).

Doctests for list under `tests/cli-dotfiles` expect discovery and array expansion. Using config-only mapping for list would silently diverge from backup/restore.

**Topic:** package layout / CLI UX (list should answer “what would this run back up?” — same discovery as `cli/dry-run` spirit).

**Recommended change:**

- Document SSoT: **CLI `list` must call `engine.MappingPaths`** (or a renamed `engine.ListJobs` / stream jobs).
- Consider deprecating or narrowing `config.MappingPaths` to “raw files-key mapping debug” or fold it into engine only to avoid two truths.
- Stream lines as produced (`cli/streaming`): print each mapping path when ready; do not buffer entire list unless filtering/sorting requires it.

---

#### H3. Single log sink merges progress, warnings, and planned output

**Evidence:** `engine.Options.Log io.Writer` receives INFO, dry-run would-lines, warnings (`logWarn` falls back to `os.Stderr` only when Log is nil), hook stdout, and success copy lines. CLI will be tempted to set `Log = os.Stdout` for everything.

**Topic:** `cli/streaming` (results → stdout, diagnostics → stderr), `cli/dry-run` (planned lines stdout; soft-fail warnings stderr).

**Recommended change:**

- Split writers in `Options`, e.g. `Out io.Writer` (user-facing plan / list / copy reports) and `Err io.Writer` (warnings, verbose diagnostics), or keep `Log` for progress and add `Warn`.
- CLI wiring:
  - `list` paths → stdout only
  - dry-run would-lines → stdout
  - `warning:` / skip noise under verbose → stderr
- Avoid buffering whole runs; keep current line-at-a-time `fmt.Fprintf` style (`cli/streaming` already matched by engine).

---

### Medium

#### M1. External commands use raw `os/exec` instead of `xgo/support/cmd`

**Topic:** `cmd-exec`.

**Evidence:** `gitOutput`, `runShell`, `runCmdToFile` all use `exec.Command` + manual buffers.

**Assessment:** Not wrong for capture-heavy cases (`git status --porcelain`, cmd stdout → file). Skill prefers `github.com/xhd2015/xgo/support/cmd` for fluent Dir/Env/Debug and default inherit of stdio.

**Recommended change:**

| Call site | Prefer |
|-----------|--------|
| `git -C repo …` with captured stdout | `cmd.Output("git", "-C", repo, …)` or `cmd.Debug().Env(...).Dir` if Dir replaces `-C`; keep `GIT_CONFIG_*` env overrides |
| Hooks `sh -c` with stdout → `opt.Log` | `cmd` with `.Stdout(opt.Log)` when no full buffer needed; keep buffer when writing atomic files from stdout |
| Verbose “print command being run” | `cmd.Debug()` instead of hand-rolled `run cmd: …` only when Verbose |

Do **not** force Debug print of every git status during large walks (noise). Use Debug/verbose selectively.

Dry-run gate (`cli/dry-run` + `cmd-exec`): never execute mutating shell hooks in dry-run (already true); git **read** probes (`status`, `rev-parse`) may still run so dry-run can report dirty vs clean accurately — that matches “same plan” better than inventing a fake dirty state. Document that dry-run may still invoke read-only git.

---

#### M2. Monolithic `internal/engine` (~1393 LOC, 40+ funcs)

**Not a named skill topic**, but it blocks clean application of dry-run, cmd-exec, and testability.

**Recommended layout (when refactoring):**

```text
internal/engine/
  options.go      // Options, Job
  resolve.go      // ResolveJobs, discovery, MappingPaths
  backup.go       // Backup + git modes
  restore.go      // Restore
  copy.go         // copyDir/File/Symlink, sameContent
  skip.go         // shouldSkipPath, logSkip
  shell.go        // runShell, runCmdToFile  (+ cmd-exec adoption)
  git.go          // gitIsDirty, gitOutput, bak.stats
```

Keep package `engine` so CLI imports stay stable.

---

#### M3. No package-level Go tests for config/engine

**Evidence:** `go test ./...` reports `[no test files]` for all packages. Correctness is gated on CLI doctests that cannot run without C1.

**Recommended change:**

- Add focused `internal/config` tests: parse order, env validate, mapping prefix match.
- Add `internal/engine` tests with temp dirs: ResolveJobs array whitelist, dry-run no-write, skip policy (pathflag + excludes).
- Keep pathflag coverage in doctests (already strong) or add table-driven `Classify` tests for fast feedback.

This is orthogonal to skill topics but required for a sound less-flags CLI (handlers stay thin).

---

#### M4. README / schema drift (`global.gitExcludes`)

**Evidence:** README documents `global.gitExcludes`; `GlobalConfig` has no such field. Git modes hardcode `.git` exclude for `gitTree` only.

**Recommended change:** Either implement `gitExcludes` in config + engine, or remove/reword README until real. Docs-only fix is appropriate and low risk.

---

#### M5. Help text contracts not yet realizable

**Topic:** `flags-parsing/subcommand` — every level needs `--help`; root help should say “run `bak-files <command> --help`…”.

**Evidence:** Foundation + backup/dotfiles help ASSERTs require Usage tokens and flag spellings (`--no-dot-files`, `--include`, `--exclude`). No help strings exist in Go yet.

**Recommended change:** Keep help as `const` strings next to each handler (skill style); empty root args → print root help exit 0; never require config/env for help.

---

### Low / informational

#### L1. Color policy not applicable yet

**Topic:** `cli/color`.

No ANSI styling today. If status lines later use color, adopt three-mode policy: auto (TTY + `NO_COLOR`) / `--color` / `--no-color`, mutual exclusion error. Not required for v1 if output stays plain.

---

#### L2. `go-embed-assets` not applicable

No SPA/extension assets, no `//go:embed` of generated trees. Do not introduce embed placeholders without a product need.

---

#### L3. `kool-create` not applicable as scaffold driver

This is a pure Go CLI module, not a kool `go-react` / server template. Keep hand layout; use skill recipes for **CLI/flags/cmd** only.

---

#### L4. `flags-parsing/cut` and `collect` not needed yet

No foreign command-line tail (`--exec …`) and no parent→child argv reconstruct. Revisit `Cut` if you add “run user command after backup” style flags; revisit `CollectParsedFlags` if a parent wrapper strips flags before re-exec.

---

#### L5. `pathflag` package layout is good

Public pure classifier (no I/O), bitmask + owners + catalog, consumed by engine for `BackupSkipMask`. Matches “small reusable library + internal product engine”. Keep it free of CLI and config concerns.

---

#### L6. Dry-run / streaming strengths already present

- Engine logs **as it walks** (good `cli/streaming`).
- Git modes share dirty detection between dry-run and live (good `cli/dry-run` fragment).
- Atomic writes via temp+rename for files and bak.stats (sound apply path).
- Access-denied soft-continue is a product feature (doctests under `cli-access-warn`); keep warnings on stderr when writers are split (H3).

---

## Recommended implementation order (for a later change pass)

Do **not** implement in this review. Suggested sequence when you do:

1. **C1 + C2:** `cmd/bak-files` with less-flags subcommands, multi-level help, env/config load, wire `list`/`backup`/`restore`.
2. **H2 + H3:** list uses `engine.MappingPaths`; split Out/Err writers; stream list lines.
3. **H1:** refactor dry-run onto shared copy/hook path; align git-mode would-hook logs.
4. **M1:** adopt `xgo/support/cmd` where capture/inherit is clearer; keep dry-run from executing mutating shells.
5. **M2 + M3:** split engine files; add unit tests so doctest is not the only safety net.
6. **M4:** docs-only README fix for `gitExcludes` (or implement if product-needed).
7. **L1:** color only if UX wants it.

---

## Checklist vs go-best-practice topics

| Topic | Status in this project | Action |
|-------|------------------------|--------|
| `flags-parsing` / `types` | Dependency present, unused | Direct require + parse in `cmd` |
| `flags-parsing/subcommand` | No dispatcher | Implement root + 3 commands + help everywhere |
| `flags-parsing/cut` | N/A | — |
| `flags-parsing/collect` | N/A | — |
| `cli/dry-run` | Partial; separate dry-run funcs | Unify pipelines; gate side effects only |
| `cli/streaming` | Engine OK if CLI streams | Print list/progress incrementally; split stderr |
| `cli/color` | Not used | Optional later |
| `cli/skill-cli` | N/A (not a skill host) | — |
| `cmd-exec` | Raw `os/exec` | Prefer `xgo/support/cmd` for git/shell where fit |
| `kool-create` | N/A | — |
| `go-embed-assets` | N/A | — |

---

## Positive notes

1. **Clear package boundaries:** pure `pathflag`, config load/mapping, engine orchestration — good Go layout once CLI is thin.
2. **Doctest-first CLI contract** is unusually complete (help, dry-run no-write, hooks, gitTree/gitPatch, pathflag catalog) — implement against tests rather than inventing UX.
3. **Home safety guards** (never full-`$HOME` job, pathflag skip mask including history for backup) are product-critical and well encoded.
4. **Whitelist array semantics** for `files` keys are explicit and tested in design docs.
5. **Atomic file writes** and symlink-as-link copy behavior show careful filesystem practice.

---

## Out of scope for this review

- Implementing `cmd/bak-files` or refactoring dry-run (explicitly deferred).
- Security audit of `sh -c` hooks (config is trusted operator input by design).
- Performance of `sameContent` (full-file reads).
- Cross-platform shell (`sh -c` assumes POSIX shell).

---

## Appendix A — Flag map for implementers

Grounded in README + doctests; use less-flags types from `flags-parsing/types`.

```text
bak-files [--help]
bak-files backup  [--config FILE] [--dry-run] [--no-dot-files]
                  [--include PATH]... [--exclude PATH]... [-v|--verbose] [--help]
bak-files restore [--config FILE] [--dry-run] [--no-dot-files]
                  [--include PATH]... [--exclude PATH]... [-v|--verbose] [--help]
bak-files list    [--config FILE] [NAMES...] [--help]
```

Map to `engine.Options`:

| Flag | Field |
|------|--------|
| `--dry-run` | `DryRun` |
| `-v/--verbose` | `Verbose` |
| `--no-dot-files` | `NoDotFiles` |
| `--include` | `DotInclude` |
| `--exclude` | `DotExclude` |
| (CLI Log) | `Log` → prefer split Out/Err (H3) |

---

## Appendix B — Severity legend

| Level | Meaning |
|-------|---------|
| Critical | Product unusable or skill baseline violated with no workaround |
| High | Behavioral drift / wrong API choice for core CLI UX recipes |
| Medium | Maintainability or secondary skill fit; fix in follow-up |
| Low | Optional / N/A topics; polish |

---

*End of review. No code changes were made beyond this document.*
