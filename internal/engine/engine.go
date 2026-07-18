// Package engine performs backup/restore copies and shell hooks using bak.config.
package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xhd2015/bak-files/internal/config"
	"github.com/xhd2015/bak-files/pathflag"
)

// Options control a backup or restore run.
type Options struct {
	DryRun  bool
	Verbose bool
	// Log is used for informational messages (stdout/stderr of the CLI).
	Log io.Writer

	// NoDotFiles disables auto-discovery of top-level $HOME dots for this run.
	NoDotFiles bool
	// DotExclude are home-relative force-exclude prefixes (CLI ∪ config).
	DotExclude []string
	// DotInclude are home-relative force-include prefixes (CLI ∪ config).
	DotInclude []string
}

// Job is one files entry resolved to source and mapping paths, plus hooks/cmds.
type Job struct {
	Key      string
	Source   string // operator filesystem path
	Mapping  string // path under targetDir (mapping path)
	Excludes []string

	// Cmd generates backup content from shell stdout (backup only).
	Cmd string
	// RestoreCmd runs on restore (side effects / optional full restore).
	RestoreCmd string
	// Hooks
	BeforeCopy    string
	AfterCopy     string
	BeforeRestore string

	// GitTree: worktree-aware directory copy + bak.stats when dirty.
	GitTree bool
	// GitPatch: write git diff HEAD as a .patch under the mapping path when dirty.
	GitPatch bool
}

// ResolveJobs expands config file entries into copy jobs (declaration order),
// then optionally merges auto-discovered top-level $HOME dots.
func ResolveJobs(cfg *config.Config, opt Options) ([]Job, error) {
	keys := cfg.FileKeys
	if len(keys) == 0 && len(cfg.Files) > 0 {
		for k := range cfg.Files {
			keys = append(keys, k)
		}
	}

	var globalEx []string
	if cfg.Global != nil {
		globalEx = append(globalEx, cfg.Global.Excludes...)
	}

	var jobs []Job
	for _, key := range keys {
		val, ok := cfg.Files[key]
		if !ok {
			continue
		}
		if b, isBool := val.(bool); isBool && !b {
			continue
		}

		entryEx, err := entryExcludes(val)
		if err != nil {
			return nil, err
		}
		ex := append([]string{}, globalEx...)
		ex = append(ex, entryEx...)

		fields := parseEntryFields(val)

		// Source path: object "file" field if set, else files key.
		srcKey := key
		if fields.File != "" {
			srcKey = fields.File
		}
		src, err := config.ExpandPath(srcKey)
		if err != nil {
			return nil, err
		}
		mapping, err := cfg.MappingPathFor(key)
		if err != nil {
			return nil, err
		}
		// gitTree copies must never pull .git into the backup store.
		if fields.GitTree {
			ex = append(ex, ".git")
		}
		jobs = append(jobs, Job{
			Key:           key,
			Source:        src,
			Mapping:       mapping,
			Excludes:      ex,
			Cmd:           fields.Cmd,
			RestoreCmd:    fields.RestoreCmd,
			BeforeCopy:    fields.BeforeCopy,
			AfterCopy:     fields.AfterCopy,
			BeforeRestore: fields.BeforeRestore,
			GitTree:       fields.GitTree,
			GitPatch:      fields.GitPatch,
		})
	}

	jobs, err := mergeDiscoveredDots(cfg, jobs, opt, globalEx)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// MappingPaths returns mapping paths for explicit and discovered jobs (list).
func MappingPaths(cfg *config.Config, opt Options) ([]string, error) {
	jobs, err := ResolveJobs(cfg, opt)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, j.Mapping)
	}
	return out, nil
}

// mergeDiscoveredDots appends jobs for top-level $HOME dots not already covered
// by an explicit job (same cleaned source path wins). Also discovers dots that
// exist under the mapped backup home root (so restore/list see store-only dots).
func mergeDiscoveredDots(cfg *config.Config, jobs []Job, opt Options, globalEx []string) ([]Job, error) {
	if opt.NoDotFiles || !cfg.IncludeDotFilesEnabled() {
		return jobs, nil
	}
	home := os.Getenv("HOME")
	if home == "" {
		return jobs, nil
	}
	home = filepath.Clean(home)

	covered := make(map[string]struct{}, len(jobs))
	for _, j := range jobs {
		covered[filepath.Clean(j.Source)] = struct{}{}
	}

	addDot := func(name string) error {
		if !strings.HasPrefix(name, ".") || name == "." || name == ".." {
			return nil
		}
		src := filepath.Join(home, name)
		if _, ok := covered[filepath.Clean(src)]; ok {
			return nil // explicit (or prior discovery) wins
		}
		key := "~/" + name
		mapping, err := cfg.MappingPathFor(key)
		if err != nil {
			return err
		}
		jobs = append(jobs, Job{
			Key:      key,
			Source:   src,
			Mapping:  mapping,
			Excludes: append([]string{}, globalEx...),
		})
		covered[filepath.Clean(src)] = struct{}{}
		return nil
	}

	// 1) Live HOME top-level dots.
	if entries, err := os.ReadDir(home); err == nil {
		for _, e := range entries {
			if err := addDot(e.Name()); err != nil {
				return nil, err
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read HOME for dot discovery: %w", err)
	}

	// 2) Backup-store top-level dots under mapped "~" (restore/list when dest absent).
	if bakRoot, err := mappedHomeBackupRoot(cfg); err == nil && bakRoot != "" {
		if entries, err := os.ReadDir(bakRoot); err == nil {
			for _, e := range entries {
				if err := addDot(e.Name()); err != nil {
					return nil, err
				}
			}
		}
	}

	return jobs, nil
}

// mappedHomeBackupRoot returns targetDir/<mapping for "~"> when resolvable.
func mappedHomeBackupRoot(cfg *config.Config) (string, error) {
	mp, err := cfg.MappingPathFor("~")
	if err != nil {
		return "", err
	}
	if mp == "" || mp == "~" {
		// No useful mapping.
		return "", nil
	}
	targetDir := cfg.TargetDir
	if targetDir == "" {
		targetDir = "./files"
	}
	return filepath.Join(targetDir, filepath.FromSlash(mp)), nil
}

type entryFields struct {
	File          string
	Cmd           string
	RestoreCmd    string
	BeforeCopy    string
	AfterCopy     string
	BeforeRestore string
	GitTree       bool
	GitPatch      bool
}

func parseEntryFields(val any) entryFields {
	var f entryFields
	m, ok := val.(map[string]any)
	if !ok {
		return f
	}
	f.File = stringField(m, "file")
	f.Cmd = stringField(m, "cmd")
	f.RestoreCmd = stringField(m, "restoreCmd")
	f.BeforeCopy = stringField(m, "beforeCopy")
	f.AfterCopy = stringField(m, "afterCopy")
	f.BeforeRestore = stringField(m, "beforeRestore")
	f.GitTree = truthyField(m, "gitTree")
	f.GitPatch = truthyField(m, "gitPatch")
	return f
}

// truthyField treats true, non-null objects, and non-empty strings as enabled.
func truthyField(m map[string]any, key string) bool {
	raw, ok := m[key]
	if !ok || raw == nil {
		return false
	}
	switch v := raw.(type) {
	case bool:
		return v
	case map[string]any:
		return true
	case string:
		return v != "" && !strings.EqualFold(v, "false")
	default:
		// numbers / arrays etc. count as present/enabled
		return true
	}
}

func stringField(m map[string]any, key string) string {
	raw, ok := m[key]
	if !ok || raw == nil {
		return ""
	}
	s, ok := raw.(string)
	if !ok {
		return ""
	}
	return s
}

func entryExcludes(val any) ([]string, error) {
	switch v := val.(type) {
	case bool, string, nil:
		return nil, nil
	case map[string]any:
		raw, ok := v["excludes"]
		if !ok || raw == nil {
			return nil, nil
		}
		return asStringSlice(raw)
	default:
		// json numbers etc. — treat as enabled entry with no excludes
		return nil, nil
	}
}

func asStringSlice(raw any) ([]string, error) {
	switch t := raw.(type) {
	case []string:
		return append([]string{}, t...), nil
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			s, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("excludes entry is not a string: %T", e)
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("excludes must be an array, got %T", raw)
	}
}

// Backup copies/generates each job: beforeCopy → cmd|file-copy → afterCopy.
func Backup(cfg *config.Config, opt Options) error {
	if opt.Log == nil {
		opt.Log = io.Discard
	}
	opt = mergeDotFilters(cfg, opt)
	jobs, err := ResolveJobs(cfg, opt)
	if err != nil {
		return err
	}
	targetDir := cfg.TargetDir
	if targetDir == "" {
		targetDir = "./files"
	}

	for _, job := range jobs {
		dst := filepath.Join(targetDir, filepath.FromSlash(job.Mapping))

		// gitTree / gitPatch take precedence over plain copy for backup.
		if job.GitTree {
			if err := backupGitTree(job, dst, opt); err != nil {
				return err
			}
			continue
		}
		if job.GitPatch {
			if err := backupGitPatch(job, dst, opt); err != nil {
				return err
			}
			continue
		}

		if opt.DryRun {
			if err := dryRunBackupJob(job, dst, opt); err != nil {
				return err
			}
			continue
		}

		// beforeCopy
		if job.BeforeCopy != "" {
			if err := runShell(job.BeforeCopy, "beforeCopy", opt); err != nil {
				return err
			}
		}

		// cmd-or-file-copy
		if job.Cmd != "" {
			if err := runCmdToFile(job.Cmd, dst, opt); err != nil {
				return err
			}
		} else {
			if err := backupCopySource(job, dst, opt); err != nil {
				return err
			}
		}

		// afterCopy
		if job.AfterCopy != "" {
			if err := runShell(job.AfterCopy, "afterCopy", opt); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeDotFilters(cfg *config.Config, opt Options) Options {
	var cfgInc, cfgExc []string
	if cfg.Global != nil {
		cfgInc = cfg.Global.DotIncludes
		cfgExc = cfg.Global.DotExcludes
	}
	// CLI filters appended after config so both apply (union).
	opt.DotInclude = append(append([]string{}, cfgInc...), opt.DotInclude...)
	opt.DotExclude = append(append([]string{}, cfgExc...), opt.DotExclude...)
	return opt
}

func dryRunBackupJob(job Job, dst string, opt Options) error {
	if job.BeforeCopy != "" {
		fmt.Fprintf(opt.Log, "dry-run: would run beforeCopy: %s\n", job.BeforeCopy)
	}
	if job.Cmd != "" {
		fmt.Fprintf(opt.Log, "dry-run: would run cmd: %s\n", job.Cmd)
		fmt.Fprintf(opt.Log, "dry-run: would write %s from cmd stdout\n", dst)
	} else {
		src := job.Source
		info, err := os.Lstat(src)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(opt.Log, "INFO: skip missing source %s\n", src)
			} else {
				return fmt.Errorf("stat source %s: %w", src, err)
			}
		} else {
			if info.Mode()&os.ModeSymlink != 0 {
				info, err = os.Stat(src)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Fprintf(opt.Log, "INFO: skip missing source %s\n", src)
					} else {
						return err
					}
					info = nil
				}
			}
			if info != nil {
				if info.IsDir() {
					if err := dryRunWalk(src, dst, job.Excludes, opt, true); err != nil {
						return err
					}
				} else {
					if skip, reason := shouldSkipPath(src, job.Excludes, opt); skip {
						logSkip(opt, src, reason)
					} else {
						fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", src, dst)
					}
				}
			}
		}
	}
	if job.AfterCopy != "" {
		fmt.Fprintf(opt.Log, "dry-run: would run afterCopy: %s\n", job.AfterCopy)
	}
	return nil
}

// dryRunWalk walks srcDir logging would-copy / skip without writing.
// backupSide: true → home-rel from source path; false → home-rel from dest path.
func dryRunWalk(srcDir, dstDir string, excludes []string, opt Options, backupSide bool) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		if rel == "." {
			// Root: evaluate skip on the directory itself.
			checkPath := path
			if !backupSide {
				checkPath = dstDir
			}
			if skip, reason := shouldSkipPath(checkPath, excludes, opt); skip {
				logSkip(opt, checkPath, reason)
				return filepath.SkipDir
			}
			fmt.Fprintf(opt.Log, "dry-run: would ensure dir %s\n", dstDir)
			return nil
		}

		checkPath := path
		if !backupSide {
			checkPath = dst
		}
		if skip, reason := shouldSkipPath(checkPath, excludes, opt); skip {
			logSkip(opt, checkPath, reason)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			fmt.Fprintf(opt.Log, "dry-run: would ensure dir %s\n", dst)
			return nil
		}
		fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", path, dst)
		return nil
	})
}

// backupGitTree backs up a dirty git worktree into dst (excluding .git) and
// writes bak.stats in the process cwd. Clean trees are skipped.
func backupGitTree(job Job, dst string, opt Options) error {
	src := job.Source
	dirty, err := gitIsDirty(src)
	if err != nil {
		return err
	}
	if !dirty {
		fmt.Fprintf(opt.Log, "INFO: skip clean gitTree %s\n", src)
		return nil
	}

	if opt.DryRun {
		fmt.Fprintf(opt.Log, "dry-run: would copy dirty gitTree %s -> %s\n", src, dst)
		fmt.Fprintf(opt.Log, "dry-run: would write bak.stats for %s\n", job.Mapping)
		return nil
	}

	if job.BeforeCopy != "" {
		if err := runShell(job.BeforeCopy, "beforeCopy", opt); err != nil {
			return err
		}
	}

	if err := copyDir(src, dst, job.Excludes, opt, true); err != nil {
		return err
	}

	hash, err := gitRevParseHEAD(src)
	if err != nil {
		return err
	}
	if err := writeBakStats(job.Mapping, hash, true); err != nil {
		return err
	}
	fmt.Fprintf(opt.Log, "wrote bak.stats for %s (hasChange=true, commitHash=%s)\n", job.Mapping, hash)

	if job.AfterCopy != "" {
		if err := runShell(job.AfterCopy, "afterCopy", opt); err != nil {
			return err
		}
	}
	return nil
}

// backupGitPatch writes git diff HEAD to worktree.patch under the mapping dir.
func backupGitPatch(job Job, dst string, opt Options) error {
	src := job.Source
	dirty, err := gitIsDirty(src)
	if err != nil {
		return err
	}
	if !dirty {
		fmt.Fprintf(opt.Log, "INFO: skip clean gitPatch %s\n", src)
		return nil
	}

	patchPath := filepath.Join(dst, "worktree.patch")
	if opt.DryRun {
		fmt.Fprintf(opt.Log, "dry-run: would write git patch %s\n", patchPath)
		return nil
	}

	if job.BeforeCopy != "" {
		if err := runShell(job.BeforeCopy, "beforeCopy", opt); err != nil {
			return err
		}
	}

	diff, err := gitDiffHEAD(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(patchPath), 0o755); err != nil {
		return fmt.Errorf("mkdir for %s: %w", patchPath, err)
	}
	if err := os.WriteFile(patchPath, []byte(diff), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(opt.Log, "wrote git patch %s\n", patchPath)

	if job.AfterCopy != "" {
		if err := runShell(job.AfterCopy, "afterCopy", opt); err != nil {
			return err
		}
	}
	return nil
}

func gitIsDirty(repo string) (bool, error) {
	out, err := gitOutput(repo, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

func gitRevParseHEAD(repo string) (string, error) {
	out, err := gitOutput(repo, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func gitDiffHEAD(repo string) (string, error) {
	// Include unstaged + staged vs HEAD. status --porcelain already established dirty.
	out, err := gitOutput(repo, "diff", "HEAD")
	if err != nil {
		return "", err
	}
	return out, nil
}

func gitOutput(repo string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", repo}, args...)
	c := exec.Command("git", cmdArgs...)
	// Keep local-only; avoid picking up unexpected system config for status/diff.
	c.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
	)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("git %v in %s: %w: %s", args, repo, err, msg)
		}
		return "", fmt.Errorf("git %v in %s: %w", args, repo, err)
	}
	return stdout.String(), nil
}

// writeBakStats merges/updates WorkDir (cwd) bak.stats for mappingKey.
func writeBakStats(mappingKey, commitHash string, hasChange bool) error {
	const statsName = "bak.stats"
	root := map[string]any{}
	if data, err := os.ReadFile(statsName); err == nil && len(bytes.TrimSpace(data)) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			// Overwrite corrupt stats rather than failing hard.
			root = map[string]any{}
		}
	}
	root[mappingKey] = map[string]any{
		"hasChange":  hasChange,
		"commitHash": commitHash,
	}
	b, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	tmp := statsName + ".tmp." + fmt.Sprintf("%d", os.Getpid())
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, statsName); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func backupCopySource(job Job, dst string, opt Options) error {
	src := job.Source
	info, err := os.Lstat(src)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(opt.Log, "INFO: skip missing source %s\n", src)
			return nil
		}
		return fmt.Errorf("stat source %s: %w", src, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		// Treat symlink as file/dir via Stat for simple P3.
		info, err = os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(opt.Log, "INFO: skip missing source %s\n", src)
				return nil
			}
			return err
		}
	}
	if info.IsDir() {
		return copyDir(src, dst, job.Excludes, opt, true)
	}
	if skip, reason := shouldSkipPath(src, job.Excludes, opt); skip {
		logSkip(opt, src, reason)
		return nil
	}
	return copyFile(src, dst, opt)
}

// Restore: beforeRestore → restoreCmd (if set) and/or file restore.
func Restore(cfg *config.Config, opt Options) error {
	if opt.Log == nil {
		opt.Log = io.Discard
	}
	opt = mergeDotFilters(cfg, opt)
	jobs, err := ResolveJobs(cfg, opt)
	if err != nil {
		return err
	}
	targetDir := cfg.TargetDir
	if targetDir == "" {
		targetDir = "./files"
	}

	for _, job := range jobs {
		bak := filepath.Join(targetDir, filepath.FromSlash(job.Mapping))
		dst := job.Source

		if opt.DryRun {
			if job.BeforeRestore != "" {
				fmt.Fprintf(opt.Log, "dry-run: would run beforeRestore: %s\n", job.BeforeRestore)
			}
			if job.RestoreCmd != "" {
				fmt.Fprintf(opt.Log, "dry-run: would run restoreCmd: %s\n", job.RestoreCmd)
			}
			info, err := os.Lstat(bak)
			if err != nil {
				if os.IsNotExist(err) {
					if job.RestoreCmd == "" {
						fmt.Fprintf(opt.Log, "INFO: skip missing backup %s\n", bak)
					}
				} else {
					return fmt.Errorf("stat backup %s: %w", bak, err)
				}
			} else {
				if info.Mode()&os.ModeSymlink != 0 {
					info, err = os.Stat(bak)
					if err != nil {
						if os.IsNotExist(err) {
							if job.RestoreCmd == "" {
								fmt.Fprintf(opt.Log, "INFO: skip missing backup %s\n", bak)
							}
						} else {
							return err
						}
						info = nil
					}
				}
				if info != nil {
					if info.IsDir() {
						// Walk backup; home-rel evaluated from restore dest paths.
						if err := dryRunWalk(bak, dst, job.Excludes, opt, false); err != nil {
							return err
						}
					} else {
						if skip, reason := shouldSkipPath(dst, job.Excludes, opt); skip {
							logSkip(opt, dst, reason)
						} else {
							fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", bak, dst)
						}
					}
				}
			}
			continue
		}

		if job.BeforeRestore != "" {
			if err := runShell(job.BeforeRestore, "beforeRestore", opt); err != nil {
				return err
			}
		}

		// restoreCmd runs as shell side-effect (stdout ignored).
		if job.RestoreCmd != "" {
			if err := runShell(job.RestoreCmd, "restoreCmd", opt); err != nil {
				return err
			}
		}

		// File restore when backup artifact exists (alongside restoreCmd when both set).
		if err := restoreCopyBackup(job, bak, dst, opt); err != nil {
			return err
		}
	}
	return nil
}

func restoreCopyBackup(job Job, bak, dst string, opt Options) error {
	info, err := os.Lstat(bak)
	if err != nil {
		if os.IsNotExist(err) {
			// If only restoreCmd was intended, missing backup is OK.
			if job.RestoreCmd != "" {
				return nil
			}
			fmt.Fprintf(opt.Log, "INFO: skip missing backup %s\n", bak)
			return nil
		}
		return fmt.Errorf("stat backup %s: %w", bak, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		info, err = os.Stat(bak)
		if err != nil {
			if os.IsNotExist(err) {
				if job.RestoreCmd != "" {
					return nil
				}
				fmt.Fprintf(opt.Log, "INFO: skip missing backup %s\n", bak)
				return nil
			}
			return err
		}
	}
	if info.IsDir() {
		return copyDir(bak, dst, job.Excludes, opt, false)
	}
	if skip, reason := shouldSkipPath(dst, job.Excludes, opt); skip {
		logSkip(opt, dst, reason)
		return nil
	}
	return copyFile(bak, dst, opt)
}

// runShell executes a shell snippet via sh -c. Non-zero exit fails the run.
func runShell(script, label string, opt Options) error {
	if opt.Verbose {
		fmt.Fprintf(opt.Log, "run %s: %s\n", label, script)
	}
	c := exec.Command("sh", "-c", script)
	var stderr bytes.Buffer
	c.Stdout = opt.Log
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("%s failed: %w: %s", label, err, msg)
		}
		return fmt.Errorf("%s failed: %w", label, err)
	}
	return nil
}

// runCmdToFile runs a shell command and writes its stdout to dst.
func runCmdToFile(script, dst string, opt Options) error {
	if opt.Verbose {
		fmt.Fprintf(opt.Log, "run cmd: %s\n", script)
	}
	c := exec.Command("sh", "-c", script)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("cmd failed: %w: %s", err, msg)
		}
		return fmt.Errorf("cmd failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir for %s: %w", dst, err)
	}
	tmp := dst + ".tmp." + fmt.Sprintf("%d", os.Getpid())
	if err := os.WriteFile(tmp, stdout.Bytes(), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		os.Remove(tmp)
		return err
	}
	fmt.Fprintf(opt.Log, "wrote %s from cmd stdout\n", dst)
	return nil
}

// copyDir walks srcDir → dstDir applying skip policy.
// backupSide: true evaluates home-rel from source paths; false from dest paths.
func copyDir(srcDir, dstDir string, excludes []string, opt Options, backupSide bool) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		if rel == "." {
			checkPath := path
			if !backupSide {
				checkPath = dstDir
			}
			if skip, reason := shouldSkipPath(checkPath, excludes, opt); skip {
				logSkip(opt, checkPath, reason)
				return filepath.SkipDir
			}
			if opt.DryRun {
				fmt.Fprintf(opt.Log, "dry-run: would ensure dir %s\n", dstDir)
				return nil
			}
			return os.MkdirAll(dstDir, 0o755)
		}

		checkPath := path
		if !backupSide {
			checkPath = dst
		}
		if skip, reason := shouldSkipPath(checkPath, excludes, opt); skip {
			logSkip(opt, checkPath, reason)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			if opt.DryRun {
				fmt.Fprintf(opt.Log, "dry-run: would ensure dir %s\n", dst)
				return nil
			}
			return os.MkdirAll(dst, info.Mode().Perm())
		}
		return copyFile(path, dst, opt)
	})
}

func copyFile(src, dst string, opt Options) error {
	if opt.DryRun {
		fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", src, dst)
		return nil
	}

	// Skip rewrite when content is already identical.
	if sameContent(src, dst) {
		if opt.Verbose {
			fmt.Fprintf(opt.Log, "INFO: identical, skip %s\n", dst)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir for %s: %w", dst, err)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Write via temp then rename for atomicity.
	tmp := dst + ".tmp." + fmt.Sprintf("%d", os.Getpid())
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		os.Remove(tmp)
		return closeErr
	}
	if err := os.Rename(tmp, dst); err != nil {
		os.Remove(tmp)
		return err
	}
	fmt.Fprintf(opt.Log, "copied %s -> %s\n", src, dst)
	return nil
}

func sameContent(a, b string) bool {
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false
	}
	if ai.Size() != bi.Size() {
		return false
	}
	ab, err := os.ReadFile(a)
	if err != nil {
		return false
	}
	bb, err := os.ReadFile(b)
	if err != nil {
		return false
	}
	return string(ab) == string(bb)
}

func excluded(base string, patterns []string) bool {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		// Exact basename match.
		if p == base {
			return true
		}
		// Glob match (e.g. *.tmp).
		if ok, _ := filepath.Match(p, base); ok {
			return true
		}
		// Also try path.Match-style with slash-normalized pattern on base only.
		if strings.Contains(p, "/") {
			if ok, _ := filepath.Match(filepath.Base(p), base); ok && filepath.Base(p) == p {
				return true
			}
		}
	}
	return false
}

// shouldSkipPath applies the home walk skip policy when absPath is under $HOME;
// otherwise only basename excludes apply.
//
// Policy (home-relative):
//  1. force-include (keep; exclude still wins if both)
//  2. force-exclude → skip
//  3. pathflag DefaultSkipMask → skip
//  4. basename/glob excludes → skip
func shouldSkipPath(absPath string, basenameExcludes []string, opt Options) (bool, string) {
	home := os.Getenv("HOME")
	homeRel, underHome := relUnderHome(absPath, home)
	if underHome {
		inc := forcePrefixMatch(homeRel, opt.DotInclude)
		exc := forcePrefixMatch(homeRel, opt.DotExclude)
		if exc {
			return true, "excluded"
		}
		if !inc {
			res, err := pathflag.Classify(homeRel)
			if err == nil && res.Flags&pathflag.DefaultSkipMask != 0 {
				reason := res.Reason
				if reason == "" {
					reason = res.Flags.String()
				}
				return true, reason
			}
		}
	}
	base := filepath.Base(absPath)
	if excluded(base, basenameExcludes) {
		return true, "excluded"
	}
	return false, ""
}

func relUnderHome(absPath, home string) (string, bool) {
	if home == "" {
		return "", false
	}
	absPath = filepath.Clean(absPath)
	home = filepath.Clean(home)
	rel, err := filepath.Rel(home, absPath)
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", false
	}
	if rel == "." {
		return "", false
	}
	return rel, true
}

func forcePrefixMatch(homeRel string, patterns []string) bool {
	homeRel = strings.TrimPrefix(filepath.ToSlash(homeRel), "./")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(filepath.ToSlash(p), "./")
		p = strings.TrimPrefix(p, "/")
		if p == "" {
			continue
		}
		if homeRel == p || strings.HasPrefix(homeRel, p+"/") {
			return true
		}
	}
	return false
}

func logSkip(opt Options, path, reason string) {
	// Prefer home-relative path in the message when under HOME.
	display := path
	if home := os.Getenv("HOME"); home != "" {
		if rel, ok := relUnderHome(path, home); ok {
			display = rel
		}
	}
	if reason != "" {
		fmt.Fprintf(opt.Log, "INFO: skip %s (%s)\n", display, reason)
	} else {
		fmt.Fprintf(opt.Log, "INFO: skip %s\n", display)
	}
}
