## Agent skills

When asked to work on a feature, bug, refactor, or complex task in FlowForge, invoke the appropriate skill:

| Phase / Intent | Skill | Role & Responsibility |
|:---|:---|:---|
| **Triage** | `/flowforge-triage` | Categorize incoming requests/bugs, check out-of-scope, create crisp brief |
| **External Material** | `/flowforge-import` | Classify local PRDs, old proposals, briefs, or notes before their facts enter requirement or design authority |
| **Align & Requirements** | `/flowforge-align` | When requirement outcomes, scope, scenarios, constraints, or terms are unsettled, persist accepted facts and hand design decisions to Solution Design |
| **Solution Design** | `/flowforge-solution-design` | Own module responsibilities, interfaces, seams, flows, migration, and verification strategy after requirements settle |
| **Spec Navigation** | `/flowforge-to-spec` | For multi-session/external review or multiple authorities needing one entry point, create optional non-authoritative navigation; skip compact work |
| **Plan & Slicing** | `/flowforge-plan` | Vertical tracer-bullet slicing with explicit DAG blocking edges (`issues/`) |
| **Implement & TDD** | `/flowforge-implement` | TDD delivery on pre-agreed seams; close out with dual-axis code review |
| **Wayfinding** | `/flowforge-wayfinder` | Fog-of-war decision mapping (`map.md`) for high-uncertainty efforts |
| **Dual-Axis Review** | `/flowforge-review` | Dual-axis (Standards vs Spec) review of implementation completeness, correctness, and standards conformance |
| **Session Handoff** | `/flowforge-handoff` | Compact session memory into cross-agent handoff artifact |
| **Architecture Probe** | `/flowforge-codebase-design` | Deep module design scan and architectural surface analysis |
| **Bug Diagnosis** | `/flowforge-diagnose` | Structured hypothesis-driven bug investigation |
| **Deep Refactoring** | `/flowforge-improve-architecture` | Comprehensive codebase scan and progressive architecture refinement |

## Generic capability dispatch

Match work to a role by capability key first; when a flowforge workflow state
owns the next step, consult the process table below (`## Subagent delegation`).

| Capability key (the work at hand) | Role | Suggested tier |
|:---|:---|:---|
| Targeted probing of one bounded question (source structure, environment, existing behavior) to close a fact gap | `flowforge-investigator` | flash pinnable |
| Batch extract + compare + summarize, split in parallel by group | `flowforge-batch-analyst` | flash pinnable |
| Write or backfill a structured document from a given template | `flowforge-scribe` | flash pinnable |
| Mechanically execute a batch of existing commands or tools | `flowforge-executor` | tool-capable |
| Standards-axis findings (build/test/convention conformance, cited) | `flowforge-reviewer-lite` | flash pinnable |
| Decision-material brief (research output contract) | `flowforge-investigator` (brief form; contract in `flowforge-research`) | flash pinnable |

Task-chain protocol (distilled from GIIS): the orchestrating session plans the
task chain, dispatches each task with the structured template below,
batch-dispatches the parallelizable tasks in one wave, Review converges the
outputs, and the session plans the next batch.

```
你是<角色>。项目根：<绝对路径>。
【目标】<一句话可验收目标>
【输入】<精确文件/目录路径与参照物>
【输出】<输出文档路径；每条结论带引用（文件路径+行号或命令输出）>
```

Dispatch boundary: only mechanically completable work goes down to a flash
tier; task-chain planning, cross-group synthesis, decisions, and Review
convergence stay in the orchestrating session.

Host note (pi): dispatched subagents fork and inherit the orchestrating
session's context as read-only reference, so workers do not need to re-read
the background; batch tasks may be dispatched asynchronously in parallel and
collected on completion wake-up; model and thinking settings can be overridden
per name via user-level `agentOverrides`.

Research landing: discussion-stage outputs land in
`<docs_dir>/research/YYYY-MM-DD-<slug>.md` with citations; before creating a
proposal, check `<docs_dir>/research/` for unconsumed notes and route them
through `/flowforge-import`.

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
| Any code review, implementation audit, or completed-work review | `flowforge-reviewer` | `flowforge-review` |
| A bounded research/diagnosis question blocks a decision | `flowforge-investigator` | `flowforge-diagnose` / `flowforge-research` |


## Execution unit policy

**One ticket, one fresh execution context.** When executing frontier tickets in
batches, dispatch each ticket to a fresh subagent or a new session — never
continue a whole proposal inside one long-lived session. Within one batch keep
the same model and tool set for all tickets (cache-friendly). Cross-ticket
state travels via artifacts, not conversation history: the ticket file, the
`STATUS:` result contract, and `flowforge frontier` — a new execution context
loses nothing it needs.

**Entry reading list.** A fresh execution context starts from artifacts, not
free exploration: this AGENTS.md, the ticket, its linked requirement/design
authorities, and the evidence left by prior tickets in the same proposal.

**Test separation.** Acceptance tests are preset by Plan/refine-ticket
(`Expected tests`); a lightweight executor must not modify preset test files
unless a Change explicitly targets them. Host-level enforcement examples:

- opencode agent definition (`.opencode/agent/*.md`):
  `permission: {edit: {"**/*_test.go": "deny", "**/src/test/**": "deny"}}`
- Claude Code: `disallowedTools` restrictions or a PostToolUse/Stop hook that
  runs the preset tests.
