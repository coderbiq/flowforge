package command

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func deployManagedAssets(targetDir string, docsRoot string) error {
	assetsDir, cleanup, err := locateAssetsDir()
	if err != nil {
		return err
	}
	defer cleanup()

	if !filepath.IsAbs(docsRoot) {
		docsRoot = filepath.Join(targetDir, docsRoot)
	}

	// Deploy skills into .agents/skills/
	if err := copyDir(filepath.Join(assetsDir, "skills"), filepath.Join(targetDir, ".agents", "skills"), true); err != nil {
		return fmt.Errorf("deploying skills: %w", err)
	}

	// Clean up skill directories that were removed from source (design gap-1).
	// Conservative rule (see ticket 06 Implementation note): delete a target
	// .agents/skills/<name>/ subdir only when (a) name is in the managed
	// flowforge-* namespace AND (b) name is not present as a directory in
	// <assetsDir>/skills/. This preserves source-shipped dirs (e.g. _shared)
	// and user-created non-flowforge-* dirs. Idempotent: a second run finds
	// nothing to delete.
	if err := cleanupRemovedSkillDirs(filepath.Join(assetsDir, "skills"), filepath.Join(targetDir, ".agents", "skills")); err != nil {
		return fmt.Errorf("cleaning removed skill dirs: %w", err)
	}

	// Deploy agent documentation rules into <docsRoot>/agents/.
	// Use overwrite=false so project-customised agent docs (e.g. standards.md,
	// domain.md, issue-tracker.md) are preserved; only new files are written.
	// Conflicts are reported via stderr and assets verify reports them as drifted.
	if err := copyDir(filepath.Join(assetsDir, "agents"), filepath.Join(docsRoot, "agents"), false); err != nil {
		return fmt.Errorf("deploying agent rules: %w", err)
	}

	// Deploy or merge AGENTS.md
	agentRules := filepath.Join(assetsDir, "AGENTS.md")
	if _, err := os.Stat(agentRules); err == nil {
		content, err := os.ReadFile(agentRules)
		if err != nil {
			return fmt.Errorf("reading AGENTS.md asset: %w", err)
		}
		targetPath := filepath.Join(targetDir, "AGENTS.md")
		if err := applyAgentsBlock(targetPath, content); err != nil {
			return fmt.Errorf("applying AGENTS.md: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking asset AGENTS.md: %w", err)
	}

	return nil
}

// managedDeployArtifactEntries lists the deploy artifact paths flowforge
// manages on every machine, pinned by the deploy-artifact-localization
// design (d-artifact-localization): the per-host agent directories, the pi
// extension directory, the skills directory, and the .flowforge/config.yaml
// single-file entry. The config entry stays single-file on purpose: the
// rest of .flowforge/ (e.g. subagents/ custom sources) stays committable.
var managedDeployArtifactEntries = []string{
	".claude/agents/",
	".opencode/agent/",
	".codex/agents/",
	".pi/agents/",
	".pi/extensions/",
	".agents/",
	".flowforge/config.yaml",
}

// gitignoreBlockMarker labels the block ensureDeployArtifactGitignore
// appends so the managed entries stay identifiable as flowforge-owned.
const gitignoreBlockMarker = "# flowforge: managed deploy artifacts (per-machine)"

// ensureDeployArtifactGitignore appends the managed deploy artifact
// entries to the project root .gitignore, creating the file when absent.
// The append is idempotent: an entry already present as a line (any
// surrounding whitespace) counts as satisfied and is never rewritten, and
// user content is preserved verbatim.
func ensureDeployArtifactGitignore(projectRoot string) error {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading .gitignore: %w", err)
	}
	content := string(data)

	present := make(map[string]bool)
	for _, line := range strings.Split(content, "\n") {
		present[strings.TrimSpace(line)] = true
	}
	var missing []string
	for _, entry := range managedDeployArtifactEntries {
		if !present[entry] {
			missing = append(missing, entry)
		}
	}
	if len(missing) == 0 {
		return nil
	}

	var block strings.Builder
	if content != "" {
		block.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			block.WriteString("\n")
		}
	}
	block.WriteString(gitignoreBlockMarker + "\n")
	for _, entry := range missing {
		block.WriteString(entry + "\n")
	}
	if err := os.WriteFile(gitignorePath, []byte(block.String()), 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}
	return nil
}

// reportTrackedDeployArtifacts checks each managed deploy artifact path
// against the local git index via `git ls-files --error-unmatch` (a
// read-only local query; exit 0 = tracked, non-zero = not tracked) and,
// for tracked paths, writes copy-paste ready `git rm --cached -r <path>`
// guidance plus a one-line reason to w. The git index is never modified;
// environments without a git binary or outside a repository are skipped
// silently.
func reportTrackedDeployArtifacts(w io.Writer, projectRoot string) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return
	}
	var tracked []string
	for _, entry := range managedDeployArtifactEntries {
		cmd := exec.Command(gitPath, "ls-files", "--error-unmatch", "--", entry)
		cmd.Dir = projectRoot
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		if err := cmd.Run(); err == nil {
			tracked = append(tracked, entry)
		}
	}
	if len(tracked) == 0 {
		return
	}
	fmt.Fprintln(w, "warning: some flowforge deploy artifacts are tracked by git, but deploy artifacts are per-machine files and should not travel via git")
	fmt.Fprintln(w, "untrack each path (the working-tree copies are kept) with:")
	for _, entry := range tracked {
		fmt.Fprintf(w, "  git rm --cached -r %s\n", entry)
	}
}

func applyAgentsBlock(targetPath string, newBlock []byte) error {
	const startMarker = "<!-- FLOWFORGE:START -->"
	const endMarker = "<!-- FLOWFORGE:END -->"

	blockStr := strings.TrimSpace(string(newBlock))
	wrappedBlock := fmt.Sprintf("%s\n%s\n%s", startMarker, blockStr, endMarker)

	existing, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		return os.WriteFile(targetPath, []byte(wrappedBlock+"\n"), 0644)
	}
	if err != nil {
		return err
	}

	content := string(existing)
	firstStart := strings.Index(content, startMarker)
	lastEnd := strings.LastIndex(content, endMarker)

	if firstStart != -1 && lastEnd != -1 && lastEnd >= firstStart {
		// Cleanly replace entire FLOWFORGE block spanning from first start to last end
		updated := strings.TrimRight(content[:firstStart], "\n") + "\n\n" + wrappedBlock + "\n"
		return os.WriteFile(targetPath, []byte(updated), 0644)
	}

	// Append block to end of file
	separator := "\n\n"
	if strings.HasSuffix(content, "\n\n") {
		separator = ""
	} else if strings.HasSuffix(content, "\n") {
		separator = "\n"
	}
	updated := content + separator + wrappedBlock + "\n"
	return os.WriteFile(targetPath, []byte(updated), 0644)
}

func locateAssetsDir() (string, func(), error) {
	noop := func() {}

	// Try embedded assets first (standalone binary)
	if dir, err := extractEmbeddedAssets(); err == nil {
		return dir, func() { os.RemoveAll(dir) }, nil
	}

	// Fallback: filesystem-based lookup (development)
	var candidates []string

	if executable, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(executable)
		candidates = append(candidates,
			filepath.Join(exeDir, "assets"),
			filepath.Join(exeDir, "..", "assets"),
		)
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
		candidates = append(candidates, filepath.Join(repoRoot, "assets"))
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "assets"))
	}

	for _, candidate := range candidates {
		if isAssetsDir(candidate) {
			return candidate, noop, nil
		}
	}

	return "", noop, fmt.Errorf("flowforge assets not found; expected assets next to the executable or in the source checkout")
}

func extractEmbeddedAssets() (string, error) {
	tmpDir, err := os.MkdirTemp(".tmp", "flowforge-assets-")
	if err != nil {
		tmpDir, err = os.MkdirTemp("", "flowforge-assets-")
		if err != nil {
			return "", fmt.Errorf("creating temp dir for embedded assets: %w", err)
		}
	}

	if err := fs.WalkDir(embeddedAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tmpDir, rel)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(embeddedAssets, path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", targetPath, err)
		}

		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", targetPath, err)
		}

		return nil
	}); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("extracting embedded assets: %w", err)
	}

	return tmpDir, nil
}

func isAssetsDir(path string) bool {
	info, err := os.Stat(filepath.Join(path, "skills"))
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

func isProjectDirectory(dir string) bool {
	indicators := []string{
		".agents/skills",
		"docs/agents",
		"AGENTS.md",
		".flowforge",
	}
	for _, ind := range indicators {
		if _, err := os.Stat(filepath.Join(dir, ind)); err == nil {
			return true
		}
	}
	return false
}

// cleanupRemovedSkillDirs deletes target .agents/skills/ subdirectories whose
// name is in the managed flowforge-* namespace but no longer exists in the
// source skills directory. Source-shipped dirs (e.g. _shared) and user-created
// non-flowforge-* dirs are preserved. Idempotent: when there is nothing to
// delete (e.g. target already cleaned, or target skills dir does not exist),
// the call is a no-op.
func cleanupRemovedSkillDirs(sourceSkillsDir, targetSkillsDir string) error {
	entries, err := os.ReadDir(targetSkillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading target skills dir %s: %w", targetSkillsDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only touch the managed flowforge-* namespace; preserve everything
		// else (e.g. _shared, user-created non-flowforge-* skill dirs).
		if !strings.HasPrefix(name, "flowforge-") {
			continue
		}
		// Preserve dirs that still ship in source.
		if _, statErr := os.Stat(filepath.Join(sourceSkillsDir, name)); statErr == nil {
			continue
		} else if !os.IsNotExist(statErr) {
			return fmt.Errorf("checking source skill dir %s: %w", filepath.Join(sourceSkillsDir, name), statErr)
		}
		targetPath := filepath.Join(targetSkillsDir, name)
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("removing stale skill dir %s: %w", targetPath, err)
		}
	}
	return nil
}

func copyDir(srcDir, dstDir string, overwrite bool) error {
	info, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking source directory %s: %w", srcDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", srcDir)
	}

	return filepath.WalkDir(srcDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Name() == ".gitkeep" {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, rel)

		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath, overwrite)
	})
}

func copyFile(srcPath, dstPath string, overwrite bool) error {
	if !overwrite {
		if existing, err := os.ReadFile(dstPath); err == nil {
			source, readErr := os.ReadFile(srcPath)
			if readErr != nil {
				return fmt.Errorf("reading source file %s: %w", srcPath, readErr)
			}
			if !bytes.Equal(existing, source) {
				fmt.Fprintf(os.Stderr, "  info: project-customised %s (preserved)\n", dstPath)
			}
			return nil
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("checking target file %s: %w", dstPath, err)
		}
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("opening source file %s: %w", srcPath, err)
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: closing source file %s: %v\n", srcPath, closeErr)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening target file %s: %w", dstPath, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		if closeErr := dst.Close(); closeErr != nil {
			return fmt.Errorf("copying file %s: %w (closing target: %v)", dstPath, err, closeErr)
		}
		return fmt.Errorf("copying file %s: %w", dstPath, err)
	}

	if err := dst.Close(); err != nil {
		return fmt.Errorf("closing target file %s: %w", dstPath, err)
	}

	return nil
}
