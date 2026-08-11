---
name: flowforge-design
description: Use ONLY to design or decompose a FlowForge proposal before implementation. Do NOT use for implementation, feedback, archive, or single-card creation.
---

# flowforge-design

If you are the main thread or Coordinator, announce the handoff and spawn `flowforge-design-analyst`. Do not design or edit Proposal artifacts locally. If native subagents are unavailable, return `BLOCKED`.

Only Active Role: Design Analyst executes below.

Run `project current`, `proposal current`, `proposal inspect`, `journal recent`, and `analysis status`. Frame the objective, evidence, non-goals, and FEATURE split. Use the simple path when evidence is sufficient; otherwise follow `references/analysis-workflow.md`.

Before FEATURE design, create at least one REQ, add it to `STR-<proposal>-REQ`, and link every FEATURE to an indexed REQ with `implements`. If investigation began in a temporary Handoff Journal, bind it immediately and continue only in the Proposal Journal.

Edit design artifacts only; never product code. Validate changes and append one concise Journal result.
