package pathflag

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// Result is the attribute outcome of Classify.
type Result struct {
	Rule   string
	Reason string
	Flags  Flag
	Owner  Owner
}

type attrRule struct {
	rule   string
	flags  Flag
	owner  Owner
	reason string
}

// Path catalog rules (exact path prefix match). Longest rule wins.
var attributeRules = []attrRule{
	{".bun", FlagCache, OwnerBun, "Bun install cache"},
	{".grok/downloads", FlagCache, OwnerGrok, "Grok downloads cache"},
	{".grok/marketplace-cache", FlagCache, OwnerGrok, "Grok plugin marketplace git cache"},
	{".grok/vendor", FlagCache | FlagVendor, OwnerGrok, "Grok vendored dependencies cache"},
	{".grok/projects", FlagCache, OwnerGrok, "Grok per-workspace agent project state"},
	{".grok/bundled", FlagCache | FlagBinary | FlagVendor, OwnerGrok, "Grok bundled skills/roles (reinstallable)"},
	{".grok/logs", FlagLogs, OwnerGrok, "Grok application logs"},
	{".android/cache", FlagCache, OwnerAndroid, "Android SDK manager index cache"},
	{".config/chromium", FlagCache, OwnerChromium, "Chromium profile cache"},
	{".cache", FlagCache, OwnerNone, "temporary application cache"},
	{".sandbox", FlagCache, OwnerNone, "sandbox working cache"},
	{".wrk", FlagCache, OwnerNone, "worktree/workspace cache"},
	{".npm", FlagCache, OwnerNpm, "npm cache"},
	{".cargo/registry", FlagCache, OwnerCargo, "Cargo registry cache"},
	{".codex/.tmp", FlagTmp | FlagCache, OwnerCodex, "Codex temporary plugin cache"},
	{".codex/skills/.system", FlagCache, OwnerCodex, "Codex system skills cache"},
	{".opencode/bin", FlagBinary, OwnerOpenCode, "OpenCode binary (reinstallable)"},
	{".local/share/cursor-agent/versions", FlagBinary | FlagCache, OwnerCursor, "Cursor agent version cache"},
	{".local/share/opencode/repos", FlagCache, OwnerOpenCode, "OpenCode repo clone cache"},
	{".local/share/opencode/snapshot", FlagCache, OwnerOpenCode, "OpenCode snapshot cache"},
	{".local/share/opencode/log", FlagLogs, OwnerOpenCode, "OpenCode application logs"},
	{".local/share/containers", FlagCache, OwnerNone, "container image/layer cache"},
	{".local/share/claude", FlagCache | FlagBinary, OwnerClaude, "Claude local share cache/binaries"},
	{".grok/sessions", FlagHistory, OwnerGrok, "Grok conversation history"},
	{".codex/sessions", FlagHistory, OwnerCodex, "Codex conversation history"},
	{".codex/logs_2.sqlite", FlagLogs, OwnerCodex, "Codex SQLite logs"},
	{".expo", FlagCache, OwnerNone, "Expo cache"},
	{".V2rayU/v2ray-core", FlagBinary, OwnerNone, "V2ray core binary"},
	{".gem", FlagCache, OwnerNone, "RubyGems cache"},
	{".ccgauge", FlagCache, OwnerNone, "ccgauge cache"},
	{".bundle", FlagCache, OwnerNone, "Bundler cache"},
	{".Trash", FlagTrash, OwnerNone, "macOS trash"},
	{".local/share/Trash", FlagTrash, OwnerNone, "Linux trash"},
	{".backup", FlagMeta, OwnerNone, "machine backup metadata (injected at pack time)"},
	{".nvm", FlagCache | FlagBinary, OwnerNvm, "nvm Node versions / toolchain"},
	{".tmp", FlagTmp | FlagCache, OwnerNone, "user temp dir"},
	{".claude/backups", FlagCache, OwnerClaude, "Claude config backups"},
	{".claude/cache", FlagCache, OwnerClaude, "Claude cache"},
	{".claude/downloads", FlagCache, OwnerClaude, "Claude downloads"},
	{".claude/projects", FlagCache, OwnerClaude, "Claude per-workspace project state"},
	{".cisco/vpn/log", FlagLogs, OwnerCisco, "Cisco VPN UI logs"},
	{".commandcode/projects", FlagCache, OwnerCommandcode, "Commandcode per-workspace project state"},
}

// Owner prefix table (longest match). Always resolved.
var ownerPrefixes = []struct {
	prefix string
	owner  Owner
}{
	{".codex", OwnerCodex},
	{".opencode", OwnerOpenCode},
	{".local/share/opencode", OwnerOpenCode},
	{".local/share/cursor-agent", OwnerCursor},
	{".grok", OwnerGrok},
	{".android", OwnerAndroid},
	{".bun", OwnerBun},
	{".npm", OwnerNpm},
	{".cargo", OwnerCargo},
	{".config/chromium", OwnerChromium},
	{".nvm", OwnerNvm},
	{".claude", OwnerClaude},
	{".cisco", OwnerCisco},
	{".commandcode", OwnerCommandcode},
}

// Classify attributes a home-relative path. No I/O.
// Zero Result if no attribute rule and no owner prefix.
// Error on empty, absolute, or ".." segments after normalize.
func Classify(relPath string) (Result, error) {
	norm, err := normalize(relPath)
	if err != nil {
		return Result{}, err
	}

	var res Result

	// 1. Longest path-prefix attribute rule
	if rule, ok := matchLongestAttr(norm); ok {
		res.Rule = rule.rule
		res.Reason = rule.reason
		res.Flags = rule.flags
		res.Owner = rule.owner
	} else if segRule, ok := matchSegment(norm); ok {
		// 2. Segment rules
		res.Rule = segRule.rule
		res.Reason = segRule.reason
		res.Flags = segRule.flags
		res.Owner = segRule.owner
	} else if strings.HasSuffix(filepath.Base(norm), ".log") {
		// 3. Basename .log suffix (case-sensitive)
		res.Rule = "**/*.log"
		res.Reason = "log files"
		res.Flags = FlagLogs
		res.Owner = OwnerNone
	}

	// 5. Always fill Owner from longest owner prefix table
	if o, ok := matchLongestOwner(norm); ok {
		res.Owner = o
	}

	return res, nil
}

func normalize(relPath string) (string, error) {
	s := strings.TrimSpace(relPath)
	s = filepath.ToSlash(s)
	for strings.HasPrefix(s, "./") {
		s = s[2:]
	}
	s = strings.TrimSpace(s)
	if s == "" || s == "." {
		return "", fmt.Errorf("empty path")
	}
	if isAbsolute(s) {
		return "", fmt.Errorf("absolute path: %q", relPath)
	}
	parts := strings.Split(s, "/")
	for _, p := range parts {
		if p == ".." {
			return "", fmt.Errorf("path contains parent segment: %q", relPath)
		}
	}
	// Drop empty segments from double slashes but keep structure simple:
	// design does not require collapsing //; tests do not use them.
	return s, nil
}

func isAbsolute(s string) bool {
	if strings.HasPrefix(s, "/") {
		return true
	}
	// Windows drive: C:/ or C:\
	if len(s) >= 2 && unicode.IsLetter(rune(s[0])) && s[1] == ':' {
		return true
	}
	return false
}

func pathPrefixMatch(path, rule string) bool {
	return path == rule || strings.HasPrefix(path, rule+"/")
}

func matchLongestAttr(path string) (attrRule, bool) {
	var best attrRule
	bestLen := -1
	for _, r := range attributeRules {
		if pathPrefixMatch(path, r.rule) && len(r.rule) > bestLen {
			best = r
			bestLen = len(r.rule)
		}
	}
	if bestLen < 0 {
		return attrRule{}, false
	}
	return best, true
}

func matchSegment(path string) (attrRule, bool) {
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if p == "node_modules" {
			return attrRule{
				rule:   "**/node_modules",
				flags:  FlagVendor,
				owner:  OwnerNone,
				reason: "node_modules directories",
			}, true
		}
	}
	for _, p := range parts {
		if p == "upload-chunks" {
			return attrRule{
				rule:   "**/upload-chunks",
				flags:  FlagTmp,
				owner:  OwnerNone,
				reason: "incomplete upload temp state",
			}, true
		}
	}
	return attrRule{}, false
}

func matchLongestOwner(path string) (Owner, bool) {
	var best Owner
	bestLen := -1
	for _, e := range ownerPrefixes {
		if pathPrefixMatch(path, e.prefix) && len(e.prefix) > bestLen {
			best = e.owner
			bestLen = len(e.prefix)
		}
	}
	if bestLen < 0 {
		return OwnerNone, false
	}
	return best, true
}
