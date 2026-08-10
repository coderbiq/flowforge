package orchestration

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed roles/*.md
var roleSources embed.FS

const sharedWorkflow = `## Shared Workflow

- The primary Coordinator is the only role that talks to the user or delegates.
- Start by reading recent Proposal Journal entries and referenced artifacts.
- Proposal, FEATURE, Step, History, and Verification are authoritative; Journal is only a collaboration index.
- Use FlowForge CLI for stages, steps, validation, preflight, risk review, and Journal updates.
- Workers never delegate or ask the user directly. Return a concise result to the Coordinator.
- Never claim verification passed unless it actually ran. Preserve unrelated worktree changes.

## Result Contract

Every result starts with one of: STATUS: COMPLETED, STATUS: BLOCKED, STATUS: DESIGN_GAP, STATUS: SCOPE_EXPANDED, STATUS: PLAN_STALE, STATUS: VERIFICATION_FAILED, or STATUS: USER_DECISION_REQUIRED.

Then report these headings: Summary, Changed Artifacts or Files, Verification, Findings or Blocker, and Next Action. Use None when a section has no entries.
`

func rolePrompt(role Role) (string, error) {
	source, err := roleSources.ReadFile("roles/" + role.ID + ".md")
	if err != nil {
		return "", fmt.Errorf("reading role source %s: %w", role.ID, err)
	}
	contract := ""
	switch role.Kind {
	case RoleKindCoordinator:
		contract = `## Routing Contract

- Delegate design, investigation, architecture, impact analysis, and replanning to flowforge-design-analyst.
- Before implementation run context preflight with explicit implementation intent; delegate only when it returns allow.
- After implementation run context risk-review. Delegate to flowforge-reviewer when it is installed and review_required; otherwise perform the review in the primary session.
- Route bugs, gaps, stale plans, and scope expansion to flowforge-feedback or flowforge-design.
`
	case RoleKindAnalyst:
		contract = `## Execution Contract

1. Load flowforge-design and read project, Proposal, Journal, artifact, knowledge, and code evidence.
2. Write design facts to FlowForge artifacts; never modify product code.
3. Run stage and validation gates, append a concise Journal entry, and return changed artifacts, decisions, gaps, and next action.
`
	case RoleKindExecutor:
		contract = `## Execution Contract

1. Load flowforge-implement and require context preflight decision allow.
2. Load the exact FEATURE Step context and implement only its declared scope.
3. Run actual verification and inspect the diff; update Step, History, and Verification before Journal.
4. Stop with design_gap, scope_expanded, plan_stale, or verification_failed instead of changing the design.
`
	case RoleKindReviewer:
		contract = `## Review Contract

1. Load flowforge-review and read FEATURE, final diff, verification evidence, and Journal.
2. Check conformance without changing product code or requesting preference-only cleanup.
3. Route implementation issues through flowforge-feedback and append the review result to Journal.
`
	}
	return fmt.Sprintf("Active Role: %s\nModel Profile: %s\nDefault Skill: %s\n\n%s\n%s\n%s", role.DisplayName, role.ModelProfile, role.DefaultSkill, sharedWorkflow, source, contract), nil
}

func RenderOpenCode(policy Policy) (map[string][]byte, error) {
	if err := policy.Validate(KnownManagedSkills()); err != nil {
		return nil, err
	}
	files := make(map[string][]byte)
	enabled := make(map[string]bool)
	for _, role := range policy.Roles {
		enabled[role.ID] = role.Enabled
	}
	for _, role := range policy.Roles {
		if !role.Enabled {
			continue
		}
		prompt, err := rolePrompt(role)
		if err != nil {
			return nil, err
		}
		mode, edit, question := "subagent", "deny", "deny"
		taskPermission := "  task: deny\n"
		bashPermission := "  bash:\n    \"*\": deny\n    \"git status*\": allow\n    \"git diff*\": allow\n    \"git log*\": allow\n    \"git show*\": allow\n"
		if role.Kind == RoleKindCoordinator {
			mode, question = "primary", "allow"
			var allow strings.Builder
			allow.WriteString("  task:\n    \"*\": deny\n")
			for _, target := range role.MayDelegate {
				if enabled[target] {
					fmt.Fprintf(&allow, "    flowforge-%s: allow\n", target)
				}
			}
			taskPermission = allow.String()
			bashPermission = ""
		}
		if role.Kind == RoleKindExecutor || role.Kind == RoleKindAnalyst {
			edit = "allow"
		}
		if role.Kind == RoleKindExecutor {
			bashPermission = "  bash:\n    \"*\": ask\n    \"git status*\": allow\n    \"git diff*\": allow\n"
		}
		content := fmt.Sprintf("---\ndescription: %q\nmode: %s\npermission:\n  edit: %s\n%s  question: %s\n  skill: allow\n%s---\n\n# FlowForge %s\n\n%s", "FlowForge "+role.DisplayName+" using Proposal artifacts, deterministic gates, and Journal handoffs.", mode, edit, taskPermission, question, bashPermission, role.DisplayName, prompt)
		files["flowforge-"+role.ID+".md"] = []byte(content)
	}
	return files, nil
}

func RenderCodex(policy Policy) (map[string][]byte, error) {
	if err := policy.Validate(KnownManagedSkills()); err != nil {
		return nil, err
	}
	files := make(map[string][]byte)
	for _, role := range policy.Roles {
		if !role.Enabled {
			continue
		}
		prompt, err := rolePrompt(role)
		if err != nil {
			return nil, err
		}
		sandbox, effort := "read-only", "medium"
		if role.Kind == RoleKindExecutor {
			sandbox = "workspace-write"
		}
		if role.Kind == RoleKindAnalyst || role.Kind == RoleKindReviewer {
			effort = "high"
		}
		prompt = strings.ReplaceAll(prompt, `"""`, `\"\"\"`)
		content := fmt.Sprintf("name = %q\ndescription = %q\nsandbox_mode = %q\nmodel_reasoning_effort = %q\ndeveloper_instructions = \"\"\"\n%s\n\"\"\"\n", "flowforge-"+role.ID, "FlowForge "+role.DisplayName, sandbox, effort, prompt)
		files["flowforge-"+role.ID+".toml"] = []byte(content)
	}
	return files, nil
}
