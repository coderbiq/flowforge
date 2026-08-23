# FlowForge Agent Configuration

## FlowForge Methodology & Agent Skills

This project adopts the **mattpocock/skills** engineering methodology, backed by the **FlowForge Local-First DAG Tracker**.

### 1. File & Storage Principles (Zero-Friction)
- All specifications (`spec.md`), tickets (`issues/NN-slug.md`), and decision maps (`map.md`) are directly created and edited as Markdown files under `.scratch/`.
- Domain terms are maintained in `CONTEXT.md` (root glossary); architectural decisions are stored in `docs/adr/NNNN-title.md`.
- **NEVER** use CLI commands to pass long-form content or multiline text. Agent uses file tools (`write`/`edit`) directly.

### 2. Graph & Frontier Engine (High-Determinism)
- Slices/Tickets declare explicit dependency edges in frontmatter: `Blocked by: 01, 02`.
- Run `flowforge frontier` to get the instant, unblocked, ready-to-execute task queue.
- Run `flowforge check` to validate dependency DAG health (cycle detection & deadlocks).

### 3. Skill Workflow Routing

When working on features, bugs, refactors, or architecture redesigns, ALWAYS invoke the matching skill for each phase:

| Phase / Intent | Invoke Skill | When & Why |
|:---|:---|:---|
| **Route & Clarify** | `/ask-matt` | Unsure which skill to use, or need meta-guidance on the workflow |
| **Triage** | `/triage` | Categorize incoming requests/bugs, check out-of-scope, create crisp brief |
| **Align & Requirements** | `/grill-with-docs` | Relentless frontier grilling; inline sync with `CONTEXT.md` & `docs/adr/` |
| **Spec Synthesis** | `/to-spec` | Synthesize consensus into unambiguous specification (`.scratch/<feature>/spec.md`) |
| **Plan & Slicing** | `/to-tickets` | Vertical tracer-bullet slicing with explicit DAG blocking edges (`issues/`) |
| **Implement & TDD** | `/implement` | TDD delivery on pre-agreed seams; close out with dual-axis code review |
| **Wayfinding** | `/wayfinder` | Fog-of-war decision mapping (`map.md`) for high-uncertainty efforts |
| **Dual-Axis Review** | `/code-review` | Dual-axis (Standards vs Spec) parallel sub-agent code inspection |
| **Session Handoff** | `/handoff` | Compact session memory into cross-agent handoff artifact |
| **Architecture Probe** | `/codebase-design` | Deep module design scan and architectural surface analysis |
| **Bug Diagnosis** | `/diagnosing-bugs` | Structured hypothesis-driven bug investigation |
