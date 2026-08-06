# CI workflow note

## Status

**updated** (workflow existed; expanded to full doctest-pattern CI and pushed)

## Branch / remote / push

| Field | Value |
|-------|--------|
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote | `origin` → `ssh://git@github.com/xhd2015/bak-files.git` |
| Push result | success (`6437e49..139ff20`, upstream set) |
| Commit SHA | `139ff20591ad9a7944a56276e4a53e58622491da` (short: `139ff20`) |

## Paths changed (in push)

- `.github/workflows/test.yml` — full pattern: setup-go, coveraged `go test`, doctest discovery + e2e stages, xgo merge, step summary, artifacts
- `script/ci/coverage-package-table.py` — package coverage markdown table for `github.com/xhd2015/bak-files/`

## How to view Actions for this push

1. Repo: https://github.com/xhd2015/bak-files  
2. Actions: https://github.com/xhd2015/bak-files/actions  
3. Filter by branch:  
   https://github.com/xhd2015/bak-files/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project  
4. Workflow name: **Test** (`.github/workflows/test.yml`)  
5. Or open the commit: https://github.com/xhd2015/bak-files/commit/139ff20591ad9a7944a56276e4a53e58622491da and use the checks / Actions tab

## How this differs from doctest’s workflow

| Aspect | doctest reference | bak-files (this push) |
|--------|-------------------|------------------------|
| `COVERPKG` | `github.com/xhd2015/doctest/...` | `github.com/xhd2015/bak-files/...` |
| Install doctest | `go install ./cmd/doctest` (checkout under test) | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Discovery / e2e | both stages | both stages (e2e may match zero leaves today) |
| Package table | `script/ci/coverage-package-table.py` for doctest module | same helper, module prefix `github.com/xhd2015/bak-files/` |
| Merge / artifacts | gotest + discovery + e2e | same profile names and merge logic |

## Prior state

A minimal workflow already ran `go test` + `doctest --label-all` in a `golang:` container without coverage merge/summary/artifacts. That was replaced by the aligned workflow above.
