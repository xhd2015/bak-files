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
)

// Options control a backup or restore run.
type Options struct {
	DryRun  bool
	Verbose bool
	// Log is used for informational messages (stdout/stderr of the CLI).
	Log io.Writer
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

// ResolveJobs expands config file entries into copy jobs (declaration order).
func ResolveJobs(cfg *config.Config) ([]Job, error) {
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
	return jobs, nil
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
	jobs, err := ResolveJobs(cfg)
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
			_ = info
			if excluded(filepath.Base(src), job.Excludes) {
				fmt.Fprintf(opt.Log, "INFO: skip excluded %s\n", src)
			} else {
				fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", src, dst)
			}
		}
	}
	if job.AfterCopy != "" {
		fmt.Fprintf(opt.Log, "dry-run: would run afterCopy: %s\n", job.AfterCopy)
	}
	return nil
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

	if err := copyDir(src, dst, job.Excludes, opt); err != nil {
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
		return copyDir(src, dst, job.Excludes, opt)
	}
	if excluded(filepath.Base(src), job.Excludes) {
		fmt.Fprintf(opt.Log, "INFO: skip excluded %s\n", src)
		return nil
	}
	return copyFile(src, dst, opt)
}

// Restore: beforeRestore → restoreCmd (if set) and/or file restore.
func Restore(cfg *config.Config, opt Options) error {
	if opt.Log == nil {
		opt.Log = io.Discard
	}
	jobs, err := ResolveJobs(cfg)
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
			// Still log would-copy when backup exists (file restore alongside restoreCmd).
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
				_ = info
				fmt.Fprintf(opt.Log, "dry-run: would copy %s -> %s\n", bak, dst)
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
		return copyDir(bak, dst, job.Excludes, opt)
	}
	if excluded(filepath.Base(bak), job.Excludes) {
		fmt.Fprintf(opt.Log, "INFO: skip excluded %s\n", bak)
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

func copyDir(srcDir, dstDir string, excludes []string, opt Options) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			// Root directory itself.
			if opt.DryRun {
				fmt.Fprintf(opt.Log, "dry-run: would ensure dir %s\n", dstDir)
				return nil
			}
			return os.MkdirAll(dstDir, 0o755)
		}
		base := filepath.Base(path)
		if excluded(base, excludes) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		dst := filepath.Join(dstDir, rel)
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
