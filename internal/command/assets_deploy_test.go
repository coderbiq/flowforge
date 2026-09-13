package command

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var localMarkdownLink = regexp.MustCompile(`\[[^]]+\]\(([^)]+)\)`)

func readSkillBody(t *testing.T, relPath string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "assets", "skills", relPath))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestPackagedSkillPointersResolve(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	assertSkillPointersResolve(t, filepath.Join(repoRoot, "assets", "skills"))
	assertRequiredArtifactContractPointers(t, filepath.Join(repoRoot, "assets", "skills"))

	target := t.TempDir()
	if err := deployManagedAssets(target, filepath.Join(target, "docs")); err != nil {
		t.Fatal(err)
	}
	assertSkillPointersResolve(t, filepath.Join(target, ".agents", "skills"))
	assertRequiredArtifactContractPointers(t, filepath.Join(target, ".agents", "skills"))

	// Verify AGENTS.md deployment includes Subagent delegation
	deployedAgents, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(deployedAgents)
	if !strings.Contains(content, "<!-- FLOWFORGE:START -->") || !strings.Contains(content, "<!-- FLOWFORGE:END -->") {
		t.Fatal("deployed AGENTS.md missing FLOWFORGE markers")
	}
	blockStart := strings.Index(content, "<!-- FLOWFORGE:START -->")
	blockEnd := strings.Index(content, "<!-- FLOWFORGE:END -->")
	if blockStart < 0 || blockEnd < 0 || blockEnd <= blockStart {
		t.Fatal("deployed AGENTS.md FLOWFORGE markers invalid")
	}
	block := content[blockStart:blockEnd]
	if !strings.Contains(block, "Subagent delegation") {
		t.Fatal("deployed AGENTS.md FLOWFORGE block missing Subagent delegation section")
	}
}

func TestAgentRulesDescribeOptionalSpecNavigation(t *testing.T) {
	agentRules, err := os.ReadFile(filepath.Join("..", "..", "assets", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(agentRules)
	if !strings.Contains(content, "optional non-authoritative navigation; skip compact work") {
		t.Fatal("AGENTS.md must trigger optional spec navigation without making it authority")
	}
	if strings.Contains(content, "Synthesize consensus into unambiguous specification") {
		t.Fatal("stale authoritative To-Spec pointer returned")
	}
}

func TestAgentRulesDescribeSubagentDelegation(t *testing.T) {
	agentRules, err := os.ReadFile(filepath.Join("..", "..", "assets", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(agentRules)
	if !strings.Contains(content, "Subagent delegation") {
		t.Fatal("AGENTS.md missing Subagent delegation section")
	}
	requiredSubagents := []string{
		"flowforge-analyst",
		"flowforge-architect",
		"flowforge-planner",
		"flowforge-implementer",
		"flowforge-reviewer",
		"flowforge-investigator",
	}
	for _, subagent := range requiredSubagents {
		if !strings.Contains(content, subagent) {
			t.Fatalf("AGENTS.md missing subagent %q in delegation table", subagent)
		}
	}
}

func TestDeployManagedAssetsUsesAbsoluteDocsRoot(t *testing.T) {
	projectRoot := t.TempDir()
	docsRoot := filepath.Join(t.TempDir(), "wiki")
	if err := deployManagedAssets(projectRoot, docsRoot); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(docsRoot, "agents", "issue-tracker.md")); err != nil {
		t.Fatalf("agent rules were not deployed to configured docs root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "docs", "agents", "issue-tracker.md")); !os.IsNotExist(err) {
		t.Fatalf("default docs root was unexpectedly populated: %v", err)
	}
}

// TestDeployCleansRemovedSkillDirs verifies that after a source removes a
// skill, the target's corresponding skill directory is cleaned during deploy.
// It exercises the cleanupRemovedSkillDirs helper that deployManagedAssets
// invokes right after the skills copyDir. Because deployManagedAssets locates
// its own source assets dir, this test drives the cleanup helper directly with
// temp dirs standing in for the two-source scenario: source v2 ships
// flowforge-A only; the target carries flowforge-A from a prior v1 deploy
// plus a stale flowforge-B that v2 of the source no longer ships. Asserts
// flowforge-A is preserved (still in source) and flowforge-B is cleaned
// (absent from source), plus idempotency on a second cleanup run.
func TestDeployCleansRemovedSkillDirs(t *testing.T) {
	source := t.TempDir()
	target := t.TempDir()

	// Source v1 ships flowforge-A only.
	sourceSkillA := filepath.Join(source, "flowforge-A")
	if err := os.MkdirAll(sourceSkillA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceSkillA, "SKILL.md"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Target reflects a prior deploy of v1 (skill A present) plus a stale
	// skill B that v2 of the source no longer ships.
	targetSkillA := filepath.Join(target, "flowforge-A")
	if err := os.MkdirAll(targetSkillA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetSkillA, "SKILL.md"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}
	targetSkillB := filepath.Join(target, "flowforge-B")
	if err := os.MkdirAll(targetSkillB, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetSkillB, "SKILL.md"), []byte("B-stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Run the cleanup step that deployManagedAssets runs after copyDir.
	if err := cleanupRemovedSkillDirs(source, target); err != nil {
		t.Fatal(err)
	}

	// Skill A (still in source) must be preserved.
	if _, err := os.Stat(filepath.Join(targetSkillA, "SKILL.md")); err != nil {
		t.Fatalf("skill A still present in source was wrongly deleted: %v", err)
	}
	// Stale skill B (removed from source) must be cleaned.
	if _, err := os.Stat(targetSkillB); !os.IsNotExist(err) {
		t.Fatalf("stale skill B (not in source) was not cleaned: stat err=%v", err)
	}

	// Idempotency: a second cleanup run on the same state must be a no-op
	// and must not error.
	if err := cleanupRemovedSkillDirs(source, target); err != nil {
		t.Fatalf("idempotent second cleanup run failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetSkillA, "SKILL.md")); err != nil {
		t.Fatalf("skill A disappeared after idempotent rerun: %v", err)
	}
	if _, err := os.Stat(targetSkillB); !os.IsNotExist(err) {
		t.Fatalf("stale skill B reappeared after idempotent rerun: stat err=%v", err)
	}
}

// TestDeployPreservesSharedDirs verifies the conservative deletion rule does
// not touch non-flowforge-* directories in target .agents/skills/, regardless
// of whether they exist in source. _shared (source-shipped) and a user-created
// non-flowforge-* skill dir are both preserved; only flowforge-* dirs absent
// from source are deleted.
func TestDeployPreservesSharedDirs(t *testing.T) {
	source := t.TempDir()
	target := t.TempDir()

	// _shared ships in source.
	sourceShared := filepath.Join(source, "_shared")
	if err := os.MkdirAll(sourceShared, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceShared, "ARTIFACT-CONTRACT.md"), []byte("contract"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Target carries: _shared (also in source), a user-created non-flowforge-*
	// dir (not in source), and a stale flowforge-* dir (not in source).
	targetShared := filepath.Join(target, "_shared")
	if err := os.MkdirAll(targetShared, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetShared, "ARTIFACT-CONTRACT.md"), []byte("contract"), 0o644); err != nil {
		t.Fatal(err)
	}
	userCreated := filepath.Join(target, "my-custom-skill")
	if err := os.MkdirAll(userCreated, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(userCreated, "SKILL.md"), []byte("user"), 0o644); err != nil {
		t.Fatal(err)
	}
	staleFlowforge := filepath.Join(target, "flowforge-stale")
	if err := os.MkdirAll(staleFlowforge, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staleFlowforge, "SKILL.md"), []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := cleanupRemovedSkillDirs(source, target); err != nil {
		t.Fatal(err)
	}

	// _shared (in source) preserved.
	if _, err := os.Stat(filepath.Join(targetShared, "ARTIFACT-CONTRACT.md")); err != nil {
		t.Fatalf("_shared (in source) was wrongly deleted: %v", err)
	}
	// User-created non-flowforge-* dir preserved (conservative rule).
	if _, err := os.Stat(filepath.Join(userCreated, "SKILL.md")); err != nil {
		t.Fatalf("user-created non-flowforge-* dir was wrongly deleted: %v", err)
	}
	// Stale flowforge-* dir (not in source) cleaned.
	if _, err := os.Stat(staleFlowforge); !os.IsNotExist(err) {
		t.Fatalf("stale flowforge-* dir (not in source) was not cleaned: stat err=%v", err)
	}
}

func TestRefineTicketSkillIsPackagedAndLinked(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	source := filepath.Join(repoRoot, "assets", "skills", "flowforge-refine-ticket", "SKILL.md")
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("flowforge-refine-ticket SKILL.md not found: %v", err)
	}
	body := readSkillBody(t, filepath.Join("flowforge-refine-ticket", "SKILL.md"))
	for _, anchor := range []string{"hand-offs", "information-value"} {
		needle := "../_shared/ARTIFACT-CONTRACT.md#" + anchor
		if !strings.Contains(body, needle) {
			t.Errorf("flowforge-refine-ticket missing contract pointer %s", needle)
		}
	}
	if !strings.Contains(body, "Verified contracts") ||
		!strings.Contains(body, "Execution scenarios") ||
		!strings.Contains(body, "Expected tests") ||
		!strings.Contains(body, "Generated artifacts") ||
		!strings.Contains(body, "Conventions") {
		t.Error("flowforge-refine-ticket must name all five execution-contract sections")
	}

	target := t.TempDir()
	if err := deployManagedAssets(target, filepath.Join(target, "docs")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, ".agents", "skills", "flowforge-refine-ticket", "SKILL.md")); err != nil {
		t.Fatalf("refine-ticket skill not deployed: %v", err)
	}
}

func TestPlanPublishesExecutionDetailSkeleton(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-plan", "SKILL.md"))
	if !strings.Contains(body, "Execution detail") {
		t.Fatal("Plan must publish an Execution detail skeleton")
	}
	if !strings.Contains(body, "refine-ticket") {
		t.Fatal("Plan must delegate fact-backed execution-contract content to flowforge-refine-ticket")
	}
}

func TestImplementSkillCarriesWeakExecutorContract(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-implement", "SKILL.md"))
	for _, needle := range []string{
		"**Mode:** lightweight",
		"STATUS: BLOCKED",
		"Restate before executing",
		"at most 2 times",
		"5 failed repair rounds",
		"exit: 0",
		"same edit",
		"verbatim",
		"Execution scenario (both Success and Failure)",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("flowforge-implement weak-executor contract missing %q", needle)
		}
	}
}

func TestImplementPreflightRejectsIncompleteExecutionContract(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-implement", "SKILL.md"))
	if !strings.Contains(body, "execution-contract-incomplete") {
		t.Fatal("Implement preflight must reject execution-contract-incomplete")
	}
	if !strings.Contains(body, "include-gaps") && !strings.Contains(body, "include_gaps") {
		t.Fatal("Implement preflight must not treat include-gaps as permission to bypass")
	}
}

func TestArtifactContractDocumentsExecutionContractSections(t *testing.T) {
	body := readSkillBody(t, filepath.Join("_shared", "ARTIFACT-CONTRACT.md"))
	for _, section := range []string{"Verified contracts", "Execution scenarios", "Expected tests", "Generated artifacts", "Conventions"} {
		if !strings.Contains(body, section) {
			t.Errorf("ARTIFACT-CONTRACT.md must document execution-contract section %q", section)
		}
	}
}

func TestAgentsBlockContainsExecutionUnitPolicy(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	data, err := os.ReadFile(filepath.Join(repoRoot, "assets", "AGENTS.md"))
	if err != nil {
		t.Fatalf("assets/AGENTS.md not found: %v", err)
	}
	body := string(data)
	for _, needle := range []string{
		"## Execution unit policy",
		"One ticket, one fresh execution context",
		"artifacts, not conversation history",
		"permission",
		"disallowedTools",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("assets/AGENTS.md execution unit policy missing %q", needle)
		}
	}
}

func TestPlanSkillDeclaresModeLine(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-plan", "SKILL.md"))
	if !strings.Contains(body, "**Mode:** lightweight") {
		t.Fatal("Plan ticket template must declare **Mode:** lightweight for mechanical tickets")
	}
	if !strings.Contains(body, "never self-selects") {
		t.Fatal("Plan must state the mode is declared, not self-selected by the implementer")
	}
	refine := readSkillBody(t, filepath.Join("flowforge-refine-ticket", "SKILL.md"))
	if !strings.Contains(refine, "**Mode:** lightweight") {
		t.Fatal("refine-ticket must backfill **Mode:** lightweight when Plan omitted it")
	}
}

func TestReviewSkillCarriesRound0Audit(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-review", "SKILL.md"))
	for _, needle := range []string{
		"Round 0",
		"converge audit",
		"missing",
		"partial",
		"contradicts",
		"unrequested",
		"file:line",
		"zero Round 0 gaps",
		"affect an acceptance term",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("flowforge-review Round 0 audit missing %q", needle)
		}
	}
}

func TestReviewSkillDescribesRepairEscalation(t *testing.T) {
	body := readSkillBody(t, filepath.Join("flowforge-review", "SKILL.md"))
	if !strings.Contains(body, "repair") {
		t.Fatal("Review skill must describe repair-ticket escalation for substantive findings")
	}
	if !strings.Contains(body, "needs-repair") {
		t.Fatal("Review skill must reference needs-repair status for original ticket")
	}
	if !strings.Contains(body, "Repair of:") {
		t.Fatal("Review skill must reference Repair of: provenance header")
	}
}

func TestArtifactContractDescribesRepairLoop(t *testing.T) {
	body := readSkillBody(t, filepath.Join("_shared", "ARTIFACT-CONTRACT.md"))
	if !strings.Contains(body, "repair") {
		t.Fatal("ARTIFACT-CONTRACT.md must describe the repair escalation path")
	}
}

func TestIssueTrackerDescribesRepairWorkflow(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "agents", "issue-tracker.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "repair") {
		t.Fatal("issue-tracker.md must describe the repair workflow")
	}
}

func assertRequiredArtifactContractPointers(t *testing.T, root string) {
	t.Helper()
	required := map[string][]string{
		"flowforge-import":          {"source-intake-and-semantic-rewrite", "roles-and-authority", "information-value"},
		"flowforge-align":           {"roles-and-authority", "hand-offs", "information-value"},
		"flowforge-to-spec":         {"roles-and-authority", "hand-offs", "information-value"},
		"flowforge-plan":            {"packaging", "hand-offs", "information-value"},
		"flowforge-wayfinder":       {"packaging", "hand-offs"},
		"flowforge-implement":       {"hand-offs", "diagnostics"},
		"flowforge-tdd":             {"hand-offs", "diagnostics"},
		"flowforge-review":          {"roles-and-authority", "information-value"},
		"flowforge-refine-ticket":   {"hand-offs", "information-value"},
		"flowforge-handoff":         {"hand-offs"},
		"flowforge-solution-design": {"roles-and-authority", "packaging", "hand-offs", "diagnostics", "information-value"},
	}
	for skill, anchors := range required {
		data, err := os.ReadFile(filepath.Join(root, skill, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		for _, anchor := range anchors {
			needle := "../_shared/ARTIFACT-CONTRACT.md#" + anchor
			if !strings.Contains(string(data), needle) {
				t.Errorf("%s is missing required contract pointer %s", skill, needle)
			}
		}
	}
}

func assertSkillPointersResolve(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range localMarkdownLink.FindAllStringSubmatch(string(data), -1) {
			parts := strings.SplitN(match[1], "#", 2)
			target := parts[0]
			if target == "" || target == "link" || strings.Contains(target, "://") {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("broken Skill pointer %s -> %s", path, match[1])
				continue
			}
			if len(parts) == 2 && parts[1] != "" {
				body, err := os.ReadFile(resolved)
				if err != nil {
					return err
				}
				if !markdownHasAnchor(string(body), parts[1]) {
					t.Errorf("broken Skill anchor %s -> %s", path, match[1])
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func markdownHasAnchor(body, wanted string) bool {
	for _, line := range strings.Split(body, "\n") {
		title := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if title == line {
			continue
		}
		slug := strings.ToLower(title)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = strings.ReplaceAll(slug, "/", "")
		if slug == wanted {
			return true
		}
	}
	return false
}
