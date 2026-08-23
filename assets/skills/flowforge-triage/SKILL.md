---
name: flowforge-triage
description: Use when receiving any new engineering task, user request, bug report, refactoring idea, or greenfield feature in FlowForge. Evaluates complexity, architectural footprint, and fog level to route to the correct workflow (Direct Patch, Tracer Bullet Feature, Hierarchical Refactor, or Wayfinder Exploration).
---

# flowforge-triage (Complexity Triage & Workflow Routing)

Assess the complexity, architectural blast radius, and uncertainty of an incoming task to choose the right engineering workflow.

## ⛔ Absolute Rules

1. **NEVER start planning or coding before classifying complexity.**
2. **NEVER force wide multi-module refactorings or greenfield systems into a flat 1-file proposal.**
3. **Be fast & decisive**: Triage in 1 turn, present the classification card with recommendation, and proceed.

## Complexity Classification Matrix

| Level | Category | Characteristics | Recommended Workflow | Working Memory Shape |
| :--- | :--- | :--- | :--- | :--- |
| **L1** | **Direct Patch / Bug** | Single file/function, existing tests fail/cover, zero API/SPI breaking changes. | `flowforge-diagnose` -> `flowforge-implement` (Red-Green TDD) | None / In-situ test |
| **L2** | **Vertical Feature** | Vertical slice across existing layers (DB -> Service -> API/UI), clear boundaries. | `flowforge-align (Flat)` -> `flowforge-plan (Tracer Bullets)` -> `flowforge-implement` | Flat Proposal (`README.md`) |
| **L3** | **Architecture Refactor** | Crosses $\ge 2$ modules/packages, class migrations, SPI changes, dependency topology shifts. | `flowforge-align (Hierarchical)` -> `flowforge-plan (Expand-Contract)` -> `flowforge-implement` | Hierarchical Proposal (`README.md` + `modules/*.md`) |
| **L4** | **Greenfield / High Fog** | New subsystem/module, unknown external dependencies, open architectural forks. | `flowforge-wayfinder` -> Spikes / Decision Tickets -> Fall back to L2/L3 | Decision Map (`MAP.md` + tickets) |

## Output Card Format (Present to User)

```markdown
### 🧭 FlowForge Task Triage

- **Task Type**: [Direct Patch | Vertical Feature | Architecture Refactor | Greenfield / High Fog]
- **Complexity Level**: [L1 | L2 | L3 | L4]
- **Estimated Footprint**: [Single File | Single Module | Multi-Module Topology | New Subsystem]
- **Uncertainty / Fog**: [Low | Medium | High]
- **Recommended Workflow**:
  1. `<Skill 1>`: <Purpose>
  2. `<Skill 2>`: <Purpose>
- **Working Memory Strategy**: [Flat Proposal | Hierarchical Mode with `modules/` | Wayfinder Decision Map]

👉 **Next Step**: Shall we proceed with `<Recommended First Skill>`?
```
