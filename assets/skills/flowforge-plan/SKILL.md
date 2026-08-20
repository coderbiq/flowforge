---
name: flowforge-plan
description: Use ONLY when a FlowForge proposal has reached alignment consensus and needs to be decomposed into executable Tracer Bullet slices and TDD work items. Do NOT use for active code implementation or early-stage exploratory discussions.
---

# flowforge-plan

Decompose aligned consensus into a series of minimal, end-to-end verifiable Tracer Bullet slices.

## Architecture & Decomposition Heuristics

### 1. Deep Modules vs Shallow Wrappers (John Ousterhout)
- Strive for deep modules: simple interface hiding complex internal mechanisms.
- Apply the **Deletion Test**: If removing this slice/module causes its internal complexity to leak across multiple files, it is a true module. If the complexity simply disappears, it is a shallow pass-through wrapper—eliminate or merge it.

### 2. Seam Discipline & Testability (Michael Feathers)
- Position seams at natural boundaries: pure business logic functions that accept explicit parameters and return values, rather than performing implicit I/O side effects.
- One adapter is an assumption; two adapters make a real interface. Avoid premature abstraction.

### 3. Tracer Bullet & Fog of War (Wayfinder Protocol)
- **Vertical Slices**: Each slice must cut end-to-end through necessary layers (15-30 min implementation). Slice 1 establishes the happy baseline path. Later slices add edge cases.
- **Fog of War Management**: Keep future, unstarted slices at low resolution (`Not yet specified`). Do not over-specify future steps until earlier slices prove the architecture in code.
- **Scale Organization**:
  - **Lean (< 6 slices)**: List directly under `## 6. Actionable Slices` in `README.md`.
  - **Multi-Module (> 6 slices)**: Group slices by phase/module (e.g. `### Phase 1: Reader (modules/01-reader.md)`).

### 4. Mandatory Test Binding (No Pseudo-Code Bloat)
- Every slice must define: **Behavior Goal**, **Seams / Files**, and an explicit **Verification Command** (`go test ...` / `pytest ...`).
- Never write rigid pseudo-code in plans—let TDD drive the implementation signatures.

## Workflow

1. Read `01-workspace/<proposal_id>/README.md` to review the objective, consensus, and facts.
2. Apply deep module and deletion tests to carve out clean boundaries.
3. Formulate 3-6 progressive Tracer Bullet slices with explicit test commands.
4. Write the slice list into `01-workspace/<proposal_id>/README.md` under `## 6. Actionable Slices`.
5. Prompt the user for a quick review before handing off to `flowforge-implement`.
