## Agent skills

When asked to work on a feature, bug, refactor, or complex task in FlowForge, invoke the appropriate skill:

| Phase / Intent | Skill | Role & Responsibility |
|:---|:---|:---|
| **Route & Guide** | `/flowforge-route` | Unsure which skill to use, or need meta-guidance on the entire workflow |
| **Triage** | `/flowforge-triage` | Categorize incoming requests/bugs, check out-of-scope, create crisp brief |
| **External Material** | `/flowforge-import` | Classify local PRDs, old proposals, briefs, or notes before their facts enter requirement or design authority |
| **Align & Requirements** | `/flowforge-align` | When requirement outcomes, scope, scenarios, constraints, or terms are unsettled, persist accepted facts and hand design decisions to Solution Design |
| **Solution Design** | `/flowforge-solution-design` | Own module responsibilities, interfaces, seams, flows, migration, and verification strategy after requirements settle |
| **Spec Navigation** | `/flowforge-to-spec` | For multi-session/external review or multiple authorities needing one entry point, create optional non-authoritative navigation; skip compact work |
| **Plan & Slicing** | `/flowforge-plan` | Vertical tracer-bullet slicing with explicit DAG blocking edges (`issues/`) |
| **Implement & TDD** | `/flowforge-implement` | TDD delivery on pre-agreed seams; close out with dual-axis code review |
| **Wayfinding** | `/flowforge-wayfinder` | Fog-of-war decision mapping (`map.md`) for high-uncertainty efforts |
| **Dual-Axis Review** | `/flowforge-review` | Dual-axis (Standards vs Spec) parallel sub-agent code inspection |
| **Session Handoff** | `/flowforge-handoff` | Compact session memory into cross-agent handoff artifact |
| **Architecture Probe** | `/flowforge-codebase-design` | Deep module design scan and architectural surface analysis |
| **Bug Diagnosis** | `/flowforge-diagnose` | Structured hypothesis-driven bug investigation |
| **Deep Refactoring** | `/flowforge-improve-architecture` | Comprehensive codebase scan and progressive architecture refinement |

## Subagent delegation

When the current session can delegate (Claude Code Agent tool / `@mention`,
OpenCode Task tool / `@mention`, Codex sub-session), prefer delegating the next
unresolved-owner step to the matching subagent instead of doing the work inline.
Consult `flowforge frontier` before choosing. Subagents do not call each other;
return to this session and re-delegate based on each subagent's `Next Action`.

| Next unresolved owner | Subagent | Bound Skill |
|:---|:---|:---|
| Requirement outcome, scope, scenario, constraint, or term unsettled | `flowforge-analyst` | `flowforge-align` |
| Requirement settled; responsibility, interface, seam, or verification strategy unsettled | `flowforge-architect` | `flowforge-solution-design` |
| Requirement and design settled; needs ticket slicing with DAG edges | `flowforge-planner` | `flowforge-plan` |
| An executable frontier ticket exists | `flowforge-implementer` | `flowforge-implement` |
| A fixed change set needs dual-axis review | `flowforge-reviewer` | `flowforge-review` |
| A bounded research/diagnosis question blocks a decision | `flowforge-investigator` | `flowforge-diagnose` / `flowforge-research` |

