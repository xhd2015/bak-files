# pathflag — Classify path attributes

Classic TDD for the pure path classifier library
`github.com/xhd2015/bak-files/pathflag`. `Classify(relPath)` attributes a
**home-relative** path (no I/O): optional **Rule** / **Reason** / **Flags**, plus
**Owner** from a longest-prefix owner table.

Out of scope: CLI/engine wiring, walk/discover, binary content detection
(`**(binary)`), user config merge, tar/backup policy.

Module: `github.com/xhd2015/bak-files`  
Package under test: `pathflag` (not implemented yet — tree must be RED)

# DSN (Domain Specific Notion)

**pathflag** is a pure **library** that classifies a single **relative path**
string as it would appear under a user home directory (e.g. `.cache/foo`,
`.codex/sessions/x`).

**Participants**

| Participant | Role |
|-------------|------|
| Caller | Supplies a home-relative path string (no filesystem access) |
| Classify | Normalizes the path, matches attribute/segment/suffix rules, resolves Owner |
| Result | Carries Rule, Reason, Flags (bitmask), and Owner (enum string) |
| Flag | `uint64` bitmask with `String()` (`tmp\|cache\|…`) and `Names()` helpers |
| Owner | Known tool/home owners (`codex`, `opencode`, `grok`, …) or empty |

**Behaviors**

1. **Normalize** — trim space, `ToSlash`, strip leading `./`; reject empty,
   absolute (`/` or Windows drive), or any `..` segment.
2. **Attribute catalog** — longest path-prefix rule wins (`rel == rule` or
   `rule + "/"`); sets Rule, Reason, Flags, and often Owner.
3. **Segment rules** — if no path rule, any component `node_modules` →
   `**/node_modules` (Vendor); `upload-chunks` → `**/upload-chunks` (Tmp).
4. **Log suffix** — if still unmatched, basename ending in `.log`
   (case-sensitive) → `**/*.log` (Logs).
5. **Owner prefix table** — longest owner prefix is always applied (even when
   Flags are zero); attribute-set Owner must agree with the table for that path.
6. **No hit** — no rule and no owner prefix → zero `Result` (empty fields), nil error.
7. **Flag helpers** — `Flag.String()` joins lower-case names with `|` in ascending
   bit order (`tmp`, `cache`, `logs`, `binary`, `trash`, `meta`, `vendor`);
   empty string if zero. `Names()` returns the same names as a slice.

## Version

0.0.2

## Decision Tree

```
tests/pathflag/                          [Request{Op, RelPath, FlagBits}]
│                                        Run: pathflag.Classify | Flag.String/Names
├── invalid/                             # Classify → error
│   ├── empty/                           # "" after normalize
│   ├── absolute/                        # leading /
│   └── parent-segment/                  # any ".." segment
├── no-hit/                              # nil error, zero Result
│   └── ordinary-path/                   # e.g. Documents/notes.txt
├── owner-only/                          # Owner set, Flags empty, no Rule
│   ├── codex/                           # .codex/sessions/...
│   └── opencode-local-share/            # .local/share/opencode (prefix only)
├── attribute/                           # catalog path rules (representative)
│   ├── cache-bun/                       # .bun → Cache, bun
│   ├── cache-home/                      # .cache → Cache
│   ├── logs-grok/                       # .grok/logs → Logs, grok
│   ├── tmp-codex/                       # .codex/.tmp → Tmp|Cache, codex
│   ├── binary-opencode/                 # .opencode/bin → Binary, opencode
│   ├── trash-macos/                     # .Trash → Trash
│   ├── meta-backup/                     # .backup → Meta
│   ├── multi-cursor-versions/           # …/cursor-agent/versions → Binary|Cache
│   └── multi-grok-vendor/               # .grok/vendor → Cache|Vendor, grok
├── match/                               # priority / segment / suffix
│   ├── longest-prefix/                  # …/opencode/log beats shorter prefixes
│   ├── segment-node-modules/            # **/node_modules
│   ├── segment-upload-chunks/           # **/upload-chunks
│   ├── log-suffix/                      # **/*.log
│   └── log-suffix-case/                 # .LOG does not match (case-sensitive)
├── normalize/                           # input sugar still classifies
│   ├── leading-dot-slash/               # ./.bun
│   └── trim-space/                      # "  .cache  "
└── flag/                                # Flag.String / Names (helpers)
    ├── string-empty/                    # 0 → ""
    ├── string-multi/                    # tmp|cache|logs order
    └── names-order/                     # Names() ascending bit order
```

**Significance order:** outcome class (invalid | no-hit | owner-only | attribute |
match-priority | normalize | flag helpers) → concrete path/rule or helper op.

## Test Index

| Leaf | Description |
|------|-------------|
| `invalid/empty` | Empty / whitespace-only → error |
| `invalid/absolute` | Absolute unix path → error |
| `invalid/parent-segment` | Path with `..` segment → error |
| `no-hit/ordinary-path` | Unrelated relative path → zero Result |
| `owner-only/codex` | Under `.codex` without catalog rule → Owner codex only |
| `owner-only/opencode-local-share` | `.local/share/opencode/...` owner prefix only |
| `attribute/cache-bun` | `.bun` → Cache, owner bun, rule `.bun` |
| `attribute/cache-home` | `.cache` → Cache, no owner |
| `attribute/logs-grok` | `.grok/logs` → Logs, grok |
| `attribute/tmp-codex` | `.codex/.tmp` → `tmp\|cache`, codex |
| `attribute/binary-opencode` | `.opencode/bin` → Binary, opencode |
| `attribute/trash-macos` | `.Trash` → Trash |
| `attribute/meta-backup` | `.backup` → Meta |
| `attribute/multi-cursor-versions` | cursor-agent versions → `cache\|binary` |
| `attribute/multi-grok-vendor` | `.grok/vendor` → `cache\|vendor`, grok |
| `match/longest-prefix` | OpenCode log path uses longest catalog rule |
| `match/segment-node-modules` | Nested `node_modules` segment → Vendor |
| `match/segment-upload-chunks` | Nested `upload-chunks` → Tmp |
| `match/log-suffix` | Basename `*.log` → Logs |
| `match/log-suffix-case` | `.LOG` does not take log rule |
| `normalize/leading-dot-slash` | `./.bun` classifies as `.bun` |
| `normalize/trim-space` | Spaced `.cache` classifies as cache |
| `flag/string-empty` | `Flag(0).String()` → `""` |
| `flag/string-multi` | Combined bits → `tmp\|cache\|logs` |
| `flag/names-order` | `Names()` ascending bit-order names |

## How to Run

From module root:

```bash
doctest vet ./tests/pathflag
doctest test ./tests/pathflag
```

Single leaf:

```bash
doctest test ./tests/pathflag/attribute/cache-bun
```

```go
import (
	"fmt"
	"testing"

	"github.com/xhd2015/bak-files/pathflag"
)

// Op values for Request.Op.
const (
	OpClassify   = "classify"
	OpFlagString = "flag_string"
	OpFlagNames  = "flag_names"
)

// Request drives Classify or Flag helper operations.
// Setup fills Op / RelPath / FlagBits; Run calls the real pathflag API.
type Request struct {
	// Op selects the harness path: classify (default), flag_string, flag_names.
	Op string
	// RelPath is the home-relative path for OpClassify.
	RelPath string
	// FlagBits is the raw Flag value for flag helper ops.
	FlagBits uint64
}

// Response mirrors Classify Result fields as assert-friendly strings,
// plus Err for validation failures and Names for Flag.Names().
type Response struct {
	Rule   string
	Reason string
	Flags  string   // Flag.String()
	Owner  string   // Owner.String()
	Names  []string // Flag.Names()
	// Err is the error message when Classify fails; empty on success.
	Err string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	op := req.Op
	if op == "" {
		op = OpClassify
	}
	switch op {
	case OpClassify:
		res, err := pathflag.Classify(req.RelPath)
		if err != nil {
			return &Response{Err: err.Error()}, nil
		}
		return &Response{
			Rule:   res.Rule,
			Reason: res.Reason,
			Flags:  res.Flags.String(),
			Owner:  res.Owner.String(),
		}, nil
	case OpFlagString:
		f := pathflag.Flag(req.FlagBits)
		return &Response{Flags: f.String()}, nil
	case OpFlagNames:
		f := pathflag.Flag(req.FlagBits)
		return &Response{Names: f.Names()}, nil
	default:
		return nil, fmt.Errorf("unknown Op %q", op)
	}
}
```
