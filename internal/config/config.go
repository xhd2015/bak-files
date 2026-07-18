// Package config loads bak.config JSON, validates required envs, and resolves
// file entries to logical mapping paths for list/backup/restore.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// DefaultConfigName is used when --config is omitted.
const DefaultConfigName = "bak.config.json"

// Config is the bak.config document.
type Config struct {
	Validate  []ValidateRule   `json:"validate"`
	Files     map[string]any   `json:"files"`
	TargetDir string           `json:"targetDir"`
	Mapping   map[string]string `json:"mapping"`
	Global    *GlobalConfig    `json:"global"`

	// FileKeys preserves declaration order of "files" object keys.
	FileKeys []string
}

// ValidateRule lists env names that must be set.
type ValidateRule struct {
	Env []string `json:"env"`
}

// GlobalConfig holds shared excludes etc.
type GlobalConfig struct {
	Excludes []string `json:"excludes"`
	// IncludeDotFiles controls auto-discovery of top-level $HOME dots.
	// nil/absent means enabled (default true).
	IncludeDotFiles *bool    `json:"includeDotFiles"`
	DotExcludes     []string `json:"dotExcludes"`
	DotIncludes     []string `json:"dotIncludes"`
}

// IncludeDotFilesEnabled reports whether home top-level dot discovery is on.
// Default is true when global is nil or includeDotFiles is absent.
func (c *Config) IncludeDotFilesEnabled() bool {
	if c == nil || c.Global == nil || c.Global.IncludeDotFiles == nil {
		return true
	}
	return *c.Global.IncludeDotFiles
}

// Load reads and parses a bak.config JSON file, preserving files key order.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	cfg, err := Parse(data)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Parse parses bak.config JSON bytes.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: invalid JSON: %w", err)
	}
	keys, err := orderedObjectKeys(data, "files")
	if err != nil {
		// If files is missing or not an object, leave FileKeys empty.
		keys = nil
	}
	cfg.FileKeys = keys
	if cfg.Files == nil {
		cfg.Files = map[string]any{}
	}
	if cfg.Mapping == nil {
		cfg.Mapping = map[string]string{}
	}
	return &cfg, nil
}

// BuiltinEnvs are always required, even if not listed in validate.
var BuiltinEnvs = []string{"HOME", "WORKING_ROLE"}

// ValidateEnvs ensures builtin envs and config.validate[].env are non-empty.
func (c *Config) ValidateEnvs() error {
	required := make([]string, 0, len(BuiltinEnvs)+4)
	required = append(required, BuiltinEnvs...)
	for _, rule := range c.Validate {
		required = append(required, rule.Env...)
	}
	seen := map[string]struct{}{}
	for _, name := range required {
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if strings.TrimSpace(os.Getenv(name)) == "" {
			return fmt.Errorf("missing env %s", name)
		}
	}
	return nil
}

type mapEntry struct {
	rawKey   string
	expKey   string
	rawValue string
}

// MappingPaths returns resolved mapping paths for each files entry, in
// declaration order. Boolean false values are skipped; other values include the key.
func (c *Config) MappingPaths() ([]string, error) {
	keys := c.FileKeys
	if len(keys) == 0 && len(c.Files) > 0 {
		// Fallback: unstable order only if order could not be recovered.
		for k := range c.Files {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	// Expand mapping keys/values for matching ($ENV and leading ~).
	entries := make([]mapEntry, 0, len(c.Mapping))
	for k, v := range c.Mapping {
		entries = append(entries, mapEntry{
			rawKey:   k,
			expKey:   expandPath(k),
			rawValue: v,
		})
	}
	// Prefer longer keys (more specific prefixes).
	sort.SliceStable(entries, func(i, j int) bool {
		return len(entries[i].rawKey) > len(entries[j].rawKey)
	})

	var out []string
	for _, key := range keys {
		val, ok := c.Files[key]
		if !ok {
			continue
		}
		if b, isBool := val.(bool); isBool && !b {
			continue
		}
		mp, err := applyMapping(key, entries)
		if err != nil {
			return nil, err
		}
		out = append(out, mp)
	}
	return out, nil
}

func applyMapping(fileKey string, entries []mapEntry) (string, error) {
	// Match on raw template form first (TS-compatible ~/Scripts vs ~).
	for _, e := range entries {
		if prefixMatch(fileKey, e.rawKey) {
			rest := fileKey[len(e.rawKey):]
			// Value may contain $ENV (e.g. HOME/$WORKING_ROLE); expand after join.
			joined := e.rawValue + rest
			return expandEnvOnly(joined), nil
		}
	}
	// No mapping: expand env/~ on the file key itself.
	return expandPath(fileKey), nil
}

// prefixMatch reports whether fileKey starts with mapKey as a path prefix.
// Exact match or mapKey followed by '/' (or mapKey empty).
func prefixMatch(fileKey, mapKey string) bool {
	if mapKey == "" {
		return true
	}
	if fileKey == mapKey {
		return true
	}
	if strings.HasPrefix(fileKey, mapKey) {
		// Avoid matching "~" against something that is not a path under ~
		// when mapKey does not end with /: require next char is '/'.
		if strings.HasSuffix(mapKey, "/") {
			return true
		}
		rest := fileKey[len(mapKey):]
		return strings.HasPrefix(rest, "/")
	}
	return false
}

// ExpandPath expands leading ~ to $HOME and $VAR / ${VAR} elsewhere.
// Exported for backup/restore engine path resolution.
func ExpandPath(s string) (string, error) {
	return expandPath(s), nil
}

// MappingPathFor returns the resolved mapping path for a single files key.
func (c *Config) MappingPathFor(fileKey string) (string, error) {
	entries := make([]mapEntry, 0, len(c.Mapping))
	for k, v := range c.Mapping {
		entries = append(entries, mapEntry{
			rawKey:   k,
			expKey:   expandPath(k),
			rawValue: v,
		})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return len(entries[i].rawKey) > len(entries[j].rawKey)
	})
	return applyMapping(fileKey, entries)
}

// expandPath expands leading ~ to $HOME and $VAR / ${VAR} elsewhere.
func expandPath(s string) string {
	if s == "~" {
		s = os.Getenv("HOME")
	} else if strings.HasPrefix(s, "~/") {
		// "~/.bashrc" → "$HOME/.bashrc" (HOME may be empty in tests for this path)
		s = os.Getenv("HOME") + s[1:] // home + "/.bashrc"
	}
	return expandEnvOnly(s)
}

// expandEnvOnly substitutes $VAR and ${VAR} using the process environment.
// Unknown vars expand to empty string (like os.ExpandEnv).
func expandEnvOnly(s string) string {
	return os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})
}

// orderedObjectKeys returns top-level object field "field"'s child key order.
func orderedObjectKeys(data []byte, field string) ([]string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil, fmt.Errorf("expected object")
	}
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key, _ := keyTok.(string)
		if key == field {
			return readObjectKeyOrder(dec)
		}
		// Skip value.
		if err := skipValue(dec); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("field %q not found", field)
}

func readObjectKeyOrder(dec *json.Decoder) ([]string, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		// null or other — no keys
		return nil, nil
	}
	var keys []string
	for dec.More() {
		kTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		k, _ := kTok.(string)
		keys = append(keys, k)
		if err := skipValue(dec); err != nil {
			return nil, err
		}
	}
	// consume closing }
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return keys, nil
}

func skipValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := tok.(json.Delim); ok {
		switch d {
		case '{', '[':
			for dec.More() {
				if d == '{' {
					// key
					if _, err := dec.Token(); err != nil {
						return err
					}
				}
				if err := skipValue(dec); err != nil {
					return err
				}
			}
			// closing delim
			_, err = dec.Token()
			return err
		}
	}
	return nil
}
