package orchestration

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed roles/*.md
var roleSources embed.FS

const sharedWorkflow = `## Shared Workflow

- The primary Coordinator is the only role that talks to the user or delegates; delegation depth is exactly one.
- Start by reading recent Proposal Journal entries and referenced artifacts.
- Proposal, FEATURE, DEC, FIND, Step, History, and Verification own durable design and execution facts. Journal owns revision and work-item scheduling facts and indexes those artifacts.
- Use FlowForge CLI for stages, steps, validation, preflight, risk review, and Journal updates.
- Workers never delegate or ask the user directly. Return a concise result to the Coordinator.
- External research is denied unless the assigned work item explicitly authorizes external sources. If required access is unavailable, return STATUS: BLOCKED.
- Never claim verification passed unless it actually ran. Preserve unrelated worktree changes.

## Result Contract

Every result starts with exactly one of: STATUS: COMPLETED, STATUS: BLOCKED, STATUS: INCONCLUSIVE, STATUS: EVIDENCE_CONFLICT, STATUS: DESIGN_GAP, STATUS: SCOPE_EXPANDED, STATUS: PLAN_STALE, STATUS: VERIFICATION_FAILED, or STATUS: USER_DECISION_REQUIRED.

Then report these headings: Summary, Changed Artifacts or Files, Verification, Findings or Blocker, and Next Action. Use None when a section has no entries.

INCONCLUSIVE means the authorized budget ended without enough evidence. EVIDENCE_CONFLICT means sources disagree and the role lacks authority to resolve them. USER_DECISION_REQUIRED pauses work and gives the Coordinator the minimum decision to ask; workers never ask the user themselves.
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

- Act as an execution-only scheduler: run deterministic CLI checks, present user-visible actions, dispatch only work already registered by the Design Analyst, and escalate exceptions. Never create investigation questions or synthesize evidence.
- Delegate framing, investigation planning, architecture, impact analysis, synthesis, and replanning to flowforge-design-analyst. Delegate each ready investigation brief directly to flowforge-investigator.
- After a result, record or seal its scheduling state, query ready work and the Analyst re-entry condition again, and invoke the Analyst only when that condition is met. Retry at most one recoverable host failure.
- Before implementation run context preflight with explicit implementation intent; delegate only when it returns allow.
- After implementation run context risk-review. Delegate to flowforge-reviewer when it is installed and review_required; otherwise perform the review in the primary session.
- Route bugs, gaps, stale plans, and scope expansion to flowforge-feedback or flowforge-design.
`
	case RoleKindAnalyst:
		contract = `## Execution Contract

1. Load flowforge-design and read project, Proposal, Journal revision state, artifacts, knowledge, and code evidence.
2. Own framing, complexity, FEATURE decomposition, investigation plans and revisions, evidence acceptance or rejection, synthesis, and stage readiness. Register every follow-up before the Coordinator may dispatch it.
3. Write design facts to Proposal, FEATURE, DEC, or FIND artifacts; never modify product code or delegate.
4. Re-enter after required results return, the budget ends, evidence conflicts, or assumptions become stale. Synthesize the revision, run gates, and return the next registered plan or the minimum user-owned decision.
`
	case RoleKindInvestigator:
		contract = `## Investigation Contract

1. Accept only a registered brief containing Proposal/FEATURE, cycle and revision, work ID, question, scope, sources, evidence requirements, budget, done_when, and one writable FIND.
2. Investigate only that question. Separate observations, inferences, and unknowns; cite reproducible sources.
3. Edit only the assigned FIND Evidence, Source, Impact, and Open Questions fields, then return the scheduling result. Never edit FEATURE, DEC, product code, or the investigation plan.
4. Return BLOCKED for missing authorized access, INCONCLUSIVE when budget expires, EVIDENCE_CONFLICT for unresolved source disagreement, and USER_DECISION_REQUIRED through the Coordinator.
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
	defaultSkill := role.DefaultSkill
	if defaultSkill == "" {
		defaultSkill = "None"
	}
	return fmt.Sprintf("Active Role: %s\nModel Profile: %s\nDefault Skill: %s\n\n%s\n%s\n%s", role.DisplayName, role.ModelProfile, defaultSkill, sharedWorkflow, source, contract), nil
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
		if role.Kind == RoleKindExecutor || role.Kind == RoleKindAnalyst || role.Kind == RoleKindInvestigator {
			edit = "allow"
		}
		if role.Kind == RoleKindAnalyst {
			bashPermission = "  bash:\n    \"*\": deny\n    \"flowforge *\": allow\n    \"./bin/flowforge *\": allow\n    \"git status*\": allow\n    \"git diff*\": allow\n    \"git log*\": allow\n    \"git show*\": allow\n"
		}
		if role.Kind == RoleKindInvestigator {
			bashPermission = "  bash:\n    \"*\": deny\n    \"rg *\": allow\n    \"sed *\": allow\n    \"git status*\": allow\n    \"git diff*\": allow\n    \"git log*\": allow\n    \"git show*\": allow\n"
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
		if role.Kind == RoleKindExecutor || role.Kind == RoleKindAnalyst || role.Kind == RoleKindInvestigator {
			sandbox = "workspace-write"
		}
		if role.Kind == RoleKindAnalyst || role.Kind == RoleKindReviewer {
			effort = "high"
		}
		prompt = strings.ReplaceAll(prompt, `"""`, `\"\"\"`)
		content := fmt.Sprintf("name = %q\ndescription = %q\nsandbox_mode = %q\nmodel_reasoning_effort = %q\ndeveloper_instructions = \"\"\"\n%s\n\"\"\"\n", "flowforge-"+role.ID, codexRoleDescription(role), sandbox, effort, prompt)
		files["flowforge-"+role.ID+".toml"] = []byte(content)
	}
	return files, nil
}

func codexRoleDescription(role Role) string {
	switch role.Kind {
	case RoleKindCoordinator:
		return "Use only when the user explicitly requests the FlowForge Coordinator as a manual fallback for routing ready design, investigation, implementation, and review work."
	case RoleKindAnalyst:
		return "Use for FlowForge Proposal design, FEATURE decomposition, architecture, impact analysis, evidence synthesis, or replanning. Do not use for product implementation."
	case RoleKindInvestigator:
		return "Use only for one ready registered FlowForge investigation brief with one assigned FIND. Do not use for broad design or product implementation."
	case RoleKindExecutor:
		return "Use only to implement one planned FlowForge FEATURE Step after context preflight returns allow and requires handoff. Do not use for design or unplanned work."
	case RoleKindReviewer:
		return "Use only for an independent FlowForge conformance review after implementation and verification are complete. Do not modify product code."
	default:
		return "FlowForge " + role.DisplayName
	}
}

// EnforcementSummary describes host guarantees that cannot be inferred from
// the role prompt alone. Conditional external access is intentionally marked
// unsupported until a host can bind it to an individual work-item brief.
func EnforcementSummary(host string) string {
	switch host {
	case "opencode":
		return "delegation=hard, interaction=hard, write_scope=soft, external_sources=unsupported"
	case "codex":
		return "sandbox=hard, delegation=soft, interaction=soft, write_scope=soft, external_sources=unsupported"
	default:
		return "unsupported"
	}
}
