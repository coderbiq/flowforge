---
name: flowforge-implement
description: Use ONLY when executing a planned FlowForge Tracer Bullet slice or implementing code and tests for a specific work item. Follows strict TDD (Red-Green-Refactor). Do NOT use for initial requirement exploration, planning, or proposal archiving.
---

# flowforge-implement

Execute a planned Tracer Bullet slice using strict Test-Driven Development (TDD).

## Implementation & TDD Heuristics (Kent Beck)

### 1. The Red-Green-Refactor Loop
- **Red**: Write exactly one focused, failing test that captures the target slice behavior. Run it and confirm it fails for the expected reason (not a syntax error or import crash).
- **Green**: Write the simplest, most direct production code to make the test pass. Avoid speculative generalization.
- **Refactor (Cognitive Load Reduction)**: 
  - **Dead Code Hunting**: Eliminate unused parameters, variables, and empty helper wrappers.
  - **Inline Shallow Code**: Inline single-use shallow wrappers that add indirection without value.
  - **Early Returns**: Flatten nested conditionals using guard clauses.
  - Keep all automated tests 100% green throughout refactoring.

### 2. Test the Interface, Not the Internals
- Test against public boundaries and behavior. Avoid asserting private methods, internal state variables, or over-mocking collaborators.
- If a component requires 5+ mocks to test, the design has a coupling problem—refactor to extract a pure helper or introduce a cleaner seam.

### 3. Diagnostic & Bug Isolation Rule
- When a test fails unexpectedly during implementation:
  - Formulate a testable hypothesis before making code changes.
  - Reproduce the failure with the smallest possible test case.
  - Fix the root cause, not the symptom.

### 4. Verification is the Sole Gate
- A slice is complete **only** when its assigned automated test command runs cleanly and project linters/type-checks pass.

## Workflow

1. Retrieve current slice details from `01-workspace/<proposal_id>/README.md` (or via `flowforge context slice`).
2. Write the failing test and run the test command to verify expected failure (Red).
3. Implement the minimal code and verify the test turns green (Green).
4. Refactor if needed, ensuring tests stay green (Refactor).
5. Mark the slice status to `[x]` in `01-workspace/<proposal_id>/README.md`.
6. Proceed to the next slice or invite the user to review.
