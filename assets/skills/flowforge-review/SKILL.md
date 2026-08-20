---
name: flowforge-review
description: Use ONLY to perform non-blocking adversarial code review, architectural drift analysis, and cognitive load simplification before merging.
---

# FlowForge Review (Adversarial Review & Cognitive Load Reduction)

Advisory red-teaming and simplification review. **Non-blocking (never blocks state or commits).**

## Review Protocol

1. **Architecture Drift**: Check changes against `docs/CONTEXT.md` constraints and domain boundaries.
2. **Cognitive Load Reduction**:
   - **Dead Code Hunting**: Identify unused parameters, variables, and abstractions.
   - **Inline Shallow Code**: Suggest inlining single-use shallow helpers.
   - **Flatten Branches**: Suggest early returns / guard clauses.
3. **Security & Concurrency**: Check race conditions, transaction leaks, and boundary validations.
