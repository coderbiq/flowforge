package command

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"flowforge/internal/config"
	"flowforge/internal/subagent"

	"gopkg.in/yaml.v3"
)

// Guards delivery integrity; model behavior needs actual controlled execution.
func TestCostAwareRulePackageSourceEmbeddedDeployed(t *testing.T) {
	target := t.TempDir()
	docsRoot := filepath.Join(target, "wiki")
	if err := deployManagedAssets(target, docsRoot); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"skills/_shared/ARTIFACT-CONTRACT.md",
		"skills/flowforge-solution-design/SKILL.md",
		"skills/flowforge-plan/SKILL.md",
		"skills/flowforge-refine-ticket/SKILL.md",
		"skills/flowforge-implement/SKILL.md",
		"skills/flowforge-review/SKILL.md",
		"agents/issue-tracker.md",
		"AGENTS.md",
	} {
		t.Run(rel, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join("..", "..", "assets", filepath.FromSlash(rel)))
			if err != nil {
				t.Fatal(err)
			}
			embedded, err := embeddedAssets.ReadFile("assets/" + rel)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(source, embedded) {
				t.Error("source and embedded rule differ; regenerate embedded assets from source")
			}
			deployedPath := filepath.Join(target, ".agents", filepath.FromSlash(rel))
			want := embedded
			switch rel {
			case "agents/issue-tracker.md":
				deployedPath = filepath.Join(docsRoot, filepath.FromSlash(rel))
			case "AGENTS.md":
				deployedPath = filepath.Join(target, rel)
				want = []byte("<!-- FLOWFORGE:START -->\n" + strings.TrimSpace(string(embedded)) + "\n<!-- FLOWFORGE:END -->\n")
			}
			deployed, err := os.ReadFile(deployedPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(want, deployed) {
				t.Error("embedded and freshly deployed rule differ")
			}
		})
	}
}

// These are clause-presence guards, not a natural-language proof. Whitespace,
// case and alternative terms are accepted; independent review checks meaning,
// contradictions and ordering, and ticket 02 observes actual role behavior.
func assertCostAwareClauses(t *testing.T, body, where string, clauses []string) {
	t.Helper()
	body = strings.Join(strings.Fields(body), " ")
	for _, clause := range clauses {
		if !regexp.MustCompile("(?i)" + clause).MatchString(body) {
			t.Errorf("%s: missing contract clause /%s/", where, clause)
		}
	}
}

func TestCostAwareContinuousDeliveryContracts(t *testing.T) {
	for _, tc := range []struct {
		path    string
		clauses []string
	}{
		{"skills/flowforge-plan/SKILL.md", []string{
			`(default|normally).{0,160}(one|single).{0,80}(group|continuous)`,
			`(split|slicing).{0,200}(reason|justify|justification)`,
			`### Verification schedule`,
			`(new ticket|new plan).{0,180}(schedule|arrangement)`,
		}},
		{"skills/_shared/ARTIFACT-CONTRACT.md", []string{
			`(default|normally).{0,160}(one|single).{0,80}(group|continuous)`,
			`(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)`,
			`(no schedule|without a schedule).{0,180}(previous|old|per.change)`,
			`(incomplete|contradictory).{0,180}(stop|return)`,
			`Fix verification`, `cmd.{0,30}exit.{0,30}output.{0,30}artifact`,
			`(stage ticket|stage proof)`, `independent review`,
		}},
		{"skills/flowforge-solution-design/SKILL.md", []string{
			`(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)`,
			`verification-ownership-and-reuse`,
		}},
		{"skills/flowforge-refine-ticket/SKILL.md", []string{
			`(unchanged|still.valid).{0,180}(refer|cit|link)`,
			`Verified contracts`, `Execution scenarios`, `Expected tests`, `Generated artifacts`, `Conventions`,
			`(fixed point|version)`, `no guesses`,
		}},
		{"skills/flowforge-implement/SKILL.md", []string{
			`(one|single|once).{0,100}(execution confirmation|execution check-in)`,
			`(group|arrangement).{0,180}(check|command)`,
			`(unchanged|still.valid).{0,180}(refer|cit|link)`,
			`(test.{0,60}(permission|author|protect)|(?:preset|protected).{0,80}test)`,
			`Fix verification`, `execution-contract-incomplete`, `at most 2 times`, `5 failed repair rounds`,
		}},
		{"skills/flowforge-review/SKILL.md", []string{
			`(main|orchestrat\w*) session.{0,80}(owner|adjudicat|writ)`,
			`(read.only|leaf).{0,160}(only|assigned|dispatched).{0,100}(audit|axis|scope|finding)`,
			`(must not|never|do not|does not).{0,80}(spawn|delegat)`,
			`(same|fixed).{0,100}(point|diff).{0,180}(batch|together)`,
			`zero Round 0 gaps`, `Fix verification`, `needs-repair`,
		}},
		{"AGENTS.md", []string{
			`(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)`,
			`(mechanical|short|small).{0,100}(record|write|backfill).{0,160}(direct|scribe|delegat)`,
			`One ticket, one fresh execution context`,
		}},
		{"agents/issue-tracker.md", []string{
			`(default|new ticket).{0,180}(schedule|group)`,
			`(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)`,
			`Fix verification`, `Stage proof before dispatch`,
		}},
	} {
		t.Run(tc.path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", "assets", filepath.FromSlash(tc.path)))
			if err != nil {
				t.Fatal(err)
			}
			assertCostAwareClauses(t, string(data), tc.path, tc.clauses)
		})
	}
	// These obsolete instructions directly contradict design rev4, rather than
	// merely using an alternative heading or harmless wording.
	for path, obsolete := range map[string][]string{
		"flowforge-plan/SKILL.md":      {"**Optional `Verification schedule`:**"},
		"flowforge-implement/SKILL.md": {"List every Constraint and every Execution scenario", "for every unchecked Change, one line restating"},
	} {
		body := strings.Join(strings.Fields(readSkillBody(t, filepath.FromSlash(path))), " ")
		for _, old := range obsolete {
			if strings.Contains(body, old) {
				t.Errorf("%s: obsolete instruction %q", path, old)
			}
		}
	}
}

func TestCostAwareRoleContracts(t *testing.T) {
	for _, role := range []string{"planner", "implementer", "reviewer", "reviewer-lite"} {
		t.Run(role, func(t *testing.T) {
			def, err := subagent.Parse(filepath.Join("..", "..", "assets", "subagents", "flowforge-"+role+".md"))
			if err != nil {
				t.Fatal(err)
			}
			switch role {
			case "planner":
				assertCostAwareClauses(t, def.Body, role, []string{`(complete|whole).{0,80}(behavior|logic|delivery)`, `(default|normally).{0,120}(group|concentrat|consolidat)`})
			case "implementer":
				assertCostAwareClauses(t, def.Body, role, []string{`(continuous|consecutive|group)`, `(must not|never|do not|does not).{0,100}(review|close)`})
			default:
				if def.Permission != "read-only" || def.DefaultSkill != "flowforge-review" {
					t.Errorf("%s must retain read-only permission and flowforge-review binding", role)
				}
				assertCostAwareClauses(t, def.Body, role, []string{
					`(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)`,
					`(must not|never|do not|does not).{0,100}(spawn|delegat)`,
					`(must not|never|do not|does not).{0,100}(write|edit|modify)`,
					`(must not|never|do not|does not).{0,140}close`,
					`(only|return).{0,80}(finding|audit|axis)`,
				})
				if role == "reviewer" {
					assertCostAwareClauses(t, def.Body, role, []string{`(default|normally).{0,100}Spec`})
				}
				body := strings.Join(strings.Fields(def.Body), " ")
				for _, old := range []string{"That Skill itself spawns", "with `Fix:` Changes appended to the ticket", "sole closeout authority", "those are adjudication outputs of flowforge-reviewer"} {
					if strings.Contains(body, old) {
						t.Errorf("%s retains obsolete owner instruction %q", role, old)
					}
				}
			}
		})
	}
}

func TestCostAwareRolesSourceEmbeddedDeployed(t *testing.T) {
	target := t.TempDir()
	if err := initializeTestProject(target); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(target)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := deploySubagents(target, cfg, ""); err != nil {
		t.Fatal(err)
	}
	for _, role := range []string{"planner", "implementer", "reviewer", "reviewer-lite"} {
		t.Run(role, func(t *testing.T) {
			name := "flowforge-" + role
			sourcePath := filepath.Join("..", "..", "assets", "subagents", name+".md")
			source, err := os.ReadFile(sourcePath)
			if err != nil {
				t.Fatal(err)
			}
			embedded, err := embeddedAssets.ReadFile("assets/subagents/" + name + ".md")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(source, embedded) {
				t.Error("source and embedded role differ; regenerate embedded assets from source")
			}
			def, err := subagent.Parse(sourcePath)
			if err != nil {
				t.Fatal(err)
			}
			for _, host := range []struct{ dir, ext string }{
				{".claude/agents", ".md"}, {".opencode/agent", ".md"},
				{".codex/agents", ".toml"}, {".pi/agents", ".md"},
			} {
				t.Run(host.dir, func(t *testing.T) {
					deployed, err := os.ReadFile(filepath.Join(target, host.dir, name+host.ext))
					if err != nil {
						t.Fatal(err)
					}
					body := string(deployed)
					// Codex rewrites the Skill invocation; all other sections must
					// preserve the source contract byte for byte on every host.
					for _, section := range strings.Split(def.Body, "\n## ") {
						if strings.HasPrefix(section, "Default Skill") {
							continue
						}
						if !strings.Contains(body, strings.TrimSpace(section)) {
							t.Error("deployed role lost a source contract section")
						}
					}
					if !strings.Contains(body, def.DefaultSkill) {
						t.Error("deployed role lost default skill binding")
					}
					if role != "reviewer" && role != "reviewer-lite" {
						return
					}
					if host.dir == ".codex/agents" && !strings.Contains(body, `sandbox_mode = "read-only"`) {
						t.Error("deployed Codex reviewer must remain read-only")
					}
					if host.dir == ".pi/agents" {
						parts := strings.SplitN(body, "---\n", 3)
						if len(parts) != 3 {
							t.Fatal("PI role missing frontmatter")
						}
						var fm struct {
							Tools []string `yaml:"tools"`
						}
						if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
							t.Fatal(err)
						}
						if !slices.Equal(fm.Tools, []string{"read", "grep", "find", "ls"}) {
							t.Errorf("deployed PI reviewer tools = %v; want read-only allowlist", fm.Tools)
						}
					}
				})
			}
		})
	}
}
