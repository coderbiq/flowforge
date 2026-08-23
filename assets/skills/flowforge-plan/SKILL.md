---
name: flowforge-plan
description: Use ONLY when the user has confirmed alignment consensus ("looks good", "let's plan") and wants to decompose the proposal into actionable, test-driven slices. Adapts slicing mode to task type (Vertical Tracer Bullets for features vs Expand-Contract Batches for wide refactorings). Do NOT use during initial exploration or TDD implementation.
---

# flowforge-plan (Polymorphic Decomposition & Slicing)

Decompose settled alignment consensus into minimal, independently verifiable work slices bound to automated tests.

## ⛔ Prerequisites & Anti-Patterns

- **PREREQUISITE**: `flowforge-align` MUST be confirmed. If open questions or boundary ambiguities remain, return to `flowforge-align`.
- **NO Big-Bang Steps**: Every slice must be executable in 15-30 minutes and independently verifiable.
- **NO Missing Seams**: Slices must explicitly list the physical files/interfaces/classes they modify.

## Dual Slicing Models

Choose slicing model based on task complexity:

### Pattern A: Vertical Tracer Bullets (For Level 2 Features)
- Slices penetrate vertically through layers (e.g. Domain -> Service -> API/UI).
- Keep downstream slices at high level (Fog-of-War).
- Format:
```markdown
## Actionable Slices (Tracer Bullets)

- [ ] **Slice 1: <Slice Title>**
  - **Objective**: <One sentence defining observable behavior>
  - **Seams / Code Files**: `<path/to/file.ext:line>`
  - **Verification Command**: `<e.g., pytest tests/... or go test ./... -run ...>`
```

### Pattern B: Expand-Contract Batches (For Level 3 Architecture Refactorings)
Do NOT force wide refactors into vertical feature slices. Use Expand -> Migrate Batches -> Contract:
1. **Batch 1 (Expand)**: Introduce new modules, contracts, interfaces, or package roots while keeping existing code green.
2. **Batch 2..N (Migrate Batches)**: Migrate callers in bounded sub-module batches (each batch runnable and verifiable).
3. **Batch Final (Contract & Lock)**: Delete deprecated code, prune unused dependencies, and enforce visibility rules (`internal`/`private`).
- Format:
```markdown
## Actionable Slices (Expand-Contract Batches)

- [ ] **Slice 1 [Expand]: <Introduce New SPI / Module Skeleton>**
  - **Objective**: <Build contract skeleton and dependencies without breaking existing code>
  - **Seams / Code Files**: `<new module path, build configs, core interfaces>`
  - **Verification Command**: `<build command for new module + full check>`

- [ ] **Slice 2 [Migrate: Batch 1]: <Migrate Subsystem A>**
  - **Objective**: <Reroute Subsystem A callers to new SPI>
  - **Module Spec Reference**: `modules/<subsystem-a>.md`
  - **Verification Command**: `<test command for Subsystem A>`

- [ ] **Slice 3 [Contract]: <Clean Deprecated Implementations>**
  - **Objective**: <Remove legacy adapters and lock package visibility>
  - **Verification Command**: `<full regression test suite>`
```
