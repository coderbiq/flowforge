---
name: flowforge-plan
description: Use ONLY when the user has confirmed the alignment consensus ("looks good", "let's plan") and wants to decompose the proposal into actionable, test-driven slices. Do NOT use during initial exploration or TDD implementation.
---

# flowforge-plan (Tracer Bullet Decomposition)

Decompose settled alignment consensus into 3-6 minimal, end-to-end **Tracer Bullet slices** bound to automated tests.

## ⛔ Prerequisites & Anti-Patterns

- **PREREQUISITE**: `flowforge-align` MUST be completed with user confirmation. If open questions or boundary ambiguities remain, return to `flowforge-align`.
- **NO Pseudo-Code Overload**: Do not write pages of hypothetical function signatures or pseudo-code in slices.
- **NO Big-Bang Steps**: Every slice must be executable in 15-30 minutes and independently verifiable.

## Slicing Principles (Ousterhout & Kent Beck)

1. **Deep Modules**: Design minimal interfaces that hide internal complexity.
2. **Tracer Bullets (Vertical Slices)**: Slice through from input to output across layers, not horizontal component-by-component layers.
3. **Mandatory Test Binding**: Every single slice MUST declare a specific, runnable test command.
4. **Fog of War (Wayfinder)**: Only detail immediate next 1-2 slices. Keep downstream slices as high-level titles.

## Output Format (Update in `01-workspace/<proposal_id>/README.md`)

```markdown
## Actionable Slices

- [ ] **Slice 1: <Slice Title>**
  - **Objective**: <One sentence defining observable behavior>
  - **Seams / Code Files**: `<path/to/file.ext:line>`
  - **Verification Command**: `<e.g., pytest tests/test_feature.py -k test_slice1 or go test ./... -run TestSlice1>`

- [ ] **Slice 2: <Slice Title>**
  - **Objective**: ...
  - **Seams / Code Files**: ...
  - **Verification Command**: ...
```
