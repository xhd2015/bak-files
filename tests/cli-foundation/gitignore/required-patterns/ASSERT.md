# .gitignore required patterns

## Expected

- `.gitignore` is readable (Run returns content, exit 0).
- Content includes (as line or path fragment, flexible matching):
  - payload dir: `files/` or `/files/`
  - stats: `bak.stats`
  - index: `sum.index`
  - secrets: `.env`
  - OS junk: `.DS_Store`
  - binary/build output: at least one of `bak-files`, `/bin/`, `bin/`, or `*.exe`

## Side Effects

- Read-only.

## Errors

- Missing file or missing required pattern fails.

## Exit Code

- **0** from successful read

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	content := resp.GitignoreContent
	if strings.TrimSpace(content) == "" {
		t.Fatal(".gitignore missing or empty")
	}

	required := []struct {
		name string
		ok   func(string) bool
	}{
		{"files/ payload dir", func(c string) bool {
			return strings.Contains(c, "files/") || strings.Contains(c, "/files")
		}},
		{"bak.stats", func(c string) bool { return strings.Contains(c, "bak.stats") }},
		{"sum.index", func(c string) bool { return strings.Contains(c, "sum.index") }},
		{".env", func(c string) bool { return strings.Contains(c, ".env") }},
		{".DS_Store", func(c string) bool { return strings.Contains(c, ".DS_Store") }},
		{"binary or bin/", func(c string) bool {
			return strings.Contains(c, "bak-files") ||
				strings.Contains(c, "/bin/") ||
				strings.Contains(c, "\nbin/") ||
				strings.HasPrefix(c, "bin/") ||
				strings.Contains(c, "*.exe") ||
				strings.Contains(c, "/bin\n")
		}},
	}

	var missing []string
	for _, r := range required {
		if !r.ok(content) {
			missing = append(missing, r.name)
		}
	}
	if len(missing) > 0 {
		t.Fatalf(".gitignore missing patterns %v\ncontent:\n%s", missing, content)
	}
}
```
