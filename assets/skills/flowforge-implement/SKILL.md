---
name: flowforge-implement
description: Use ONLY when the user asks to execute a ready FlowForge implementation task, or provides a task id and wants code changes for that task. Do NOT use for design, analysis, feedback, archive, or general card lookup.
---

# flowforge-implement

If you are the main thread or Coordinator, run `journal recent` and `context preflight`. On `allow` with required handoff, announce and spawn `flowforge-executor` with the exact FEATURE Step context. Never implement locally. If native subagents are unavailable, return `BLOCKED`.

Only Active Role: Executor executes below.

Start with `context feature --feature <id> --step <n>`, then follow `references/workflow-rules.md`. Implement only the declared Step, run its verification, update Step/History/Verification, validate card state, then append one Journal result. Stop on design gaps, scope expansion, stale plans, or verification failure; never redesign.
