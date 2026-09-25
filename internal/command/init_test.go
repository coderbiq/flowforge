package command

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// wantManagedGitignoreEntries mirrors the per-machine deploy artifact
// entries pinned by the deploy-artifact-localization design
// (d-artifact-localization): the per-host agent directories, the pi
// extension directory, the skills directory, and the .flowforge/config.yaml
// single-file entry (the rest of .flowforge/ — e.g. subagents/ custom
// sources — stays committable, so a bare `.flowforge/` entry is wrong).
var wantManagedGitignoreEntries = []string{
	".claude/agents/",
	".opencode/agent/",
	".codex/agents/",
	".pi/agents/",
	".pi/extensions/",
	".agents/",
	".flowforge/config.yaml",
}

// gitignoreLineCounts counts exact lines (whitespace-trimmed) of a
// .gitignore body so idempotency can be asserted per entry.
func gitignoreLineCounts(content string) map[string]int {
	counts := make(map[string]int)
	for _, line := range strings.Split(content, "\n") {
		counts[strings.TrimSpace(line)]++
	}
	return counts
}

func runGitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v in %s: %v", args, dir, err)
	}
	return out.String()
}

// TestInitDeploysManagedGitignore verifies init writes the managed deploy
// artifact entries into the project root .gitignore (creating the file when
// absent), pins the single-file config entry against a bare `.flowforge/`
// directory entry, and is idempotent: a second init leaves the file
// byte-identical.
func TestInitDeploysManagedGitignore(t *testing.T) {
	projectRoot := t.TempDir()
	cmd := newInitCmd()
	if err := cmd.RunE(cmd, []string{projectRoot}); err != nil {
		t.Fatal(err)
	}

	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	first, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("init did not manage .gitignore: %v", err)
	}
	content := string(first)

	if !strings.Contains(content, "# flowforge: managed deploy artifacts (per-machine)") {
		t.Fatalf("gitignore block marker missing:\n%s", content)
	}
	counts := gitignoreLineCounts(content)
	for _, entry := range wantManagedGitignoreEntries {
		if counts[entry] != 1 {
			t.Errorf("gitignore entry %q appears %d times, want exactly 1:\n%s", entry, counts[entry], content)
		}
	}
	// The config entry must stay single-file: a bare `.flowforge/` line
	// would ignore the whole directory and keep custom sources out of git.
	if counts[".flowforge/"] != 0 {
		t.Errorf("gitignore must not carry a bare .flowforge/ directory entry:\n%s", content)
	}

	cmd2 := newInitCmd()
	if err := cmd2.RunE(cmd2, []string{projectRoot}); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("second init changed .gitignore:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// TestGitignoreAppendIsIdempotentAndPreservesUserContent verifies the
// idempotency and user-content rules on a pre-existing .gitignore: a user
// entry already naming a managed path counts as satisfied (no duplicate
// line), user lines are preserved verbatim, and a repeat run is a no-op.
func TestGitignoreAppendIsIdempotentAndPreservesUserContent(t *testing.T) {
	projectRoot := t.TempDir()
	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	seed := "# user ignore rules\nnode_modules/\n.claude/agents/" // no trailing newline, one entry pre-satisfied
	if err := os.WriteFile(gitignorePath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ensureDeployArtifactGitignore(projectRoot); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(first)
	if !strings.HasPrefix(content, "# user ignore rules\nnode_modules/\n.claude/agents/\n") {
		t.Fatalf("user content was not preserved verbatim:\n%s", content)
	}
	counts := gitignoreLineCounts(content)
	if counts[".claude/agents/"] != 1 {
		t.Errorf("pre-satisfied entry duplicated (%d occurrences):\n%s", counts[".claude/agents/"], content)
	}
	for _, entry := range wantManagedGitignoreEntries {
		if counts[entry] != 1 {
			t.Errorf("gitignore entry %q appears %d times, want exactly 1:\n%s", entry, counts[entry], content)
		}
	}

	if err := ensureDeployArtifactGitignore(projectRoot); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("repeat ensure call changed .gitignore:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// TestUpgradeSyncManagesGitignore verifies the upgrade sync channel keeps
// the gitignore localization conservative: syncProjectAssets writes the
// managed entries to the project root even when invoked from a nested
// working directory.
func TestUpgradeSyncManagesGitignore(t *testing.T) {
	projectRoot := t.TempDir()
	configDir := filepath.Join(projectRoot, ".flowforge")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("version: 5.0.0\ndocs_dir: wiki\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(projectRoot, "src")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, nested)

	cmd := newUpgradeCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	syncProjectAssets(cmd, "synced")

	data, err := os.ReadFile(filepath.Join(projectRoot, ".gitignore"))
	if err != nil {
		t.Fatalf("upgrade sync did not manage .gitignore: %v", err)
	}
	counts := gitignoreLineCounts(string(data))
	for _, entry := range wantManagedGitignoreEntries {
		if counts[entry] != 1 {
			t.Errorf("upgrade-synced gitignore entry %q appears %d times, want exactly 1:\n%s", entry, counts[entry], data)
		}
	}
}

// TestGitignoreTrackedArtifactGuidance verifies the pre-deploy detection:
// a tracked managed path yields copy-paste `git rm --cached -r <path>`
// guidance plus the per-machine reason on the writer, untracked managed
// paths stay silent, the git index is never modified, and a directory
// without a git repository is skipped silently.
func TestGitignoreTrackedArtifactGuidance(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in test environment")
	}

	projectRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectRoot, ".claude", "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, ".claude", "agents", "a.md"), []byte("agent"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, projectRoot, "init", "-q", ".")
	runGitIn(t, projectRoot, "add", "-A")

	indexBefore := runGitIn(t, projectRoot, "ls-files", "--stage")

	var buf bytes.Buffer
	reportTrackedDeployArtifacts(&buf, projectRoot)
	guidance := buf.String()

	if !strings.Contains(guidance, "git rm --cached -r .claude/agents/") {
		t.Fatalf("tracked managed path missing copy-paste untrack guidance:\n%s", guidance)
	}
	if !strings.Contains(guidance, "per-machine") {
		t.Fatalf("guidance missing per-machine reason sentence:\n%s", guidance)
	}
	if strings.Contains(guidance, ".pi/agents/") {
		t.Fatalf("untracked managed path must stay silent:\n%s", guidance)
	}

	indexAfter := runGitIn(t, projectRoot, "ls-files", "--stage")
	if indexBefore != indexAfter {
		t.Fatalf("detection modified the git index:\nbefore:\n%s\nafter:\n%s", indexBefore, indexAfter)
	}

	// Non-git directory: silent skip, no error.
	plain := t.TempDir()
	var silent bytes.Buffer
	reportTrackedDeployArtifacts(&silent, plain)
	if silent.Len() != 0 {
		t.Fatalf("non-git environment must be skipped silently, got:\n%s", silent.String())
	}
}
