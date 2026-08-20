---
name: flowforge-explore
description: Use ONLY when an unverified assumption, uncertain legacy code behavior, unknown schema, or third-party library nuance arises during FlowForge design or planning. Researches codebase facts and injects evidence back into the Proposal Scratchpad. Do NOT use for general requirement debate or production code modification.
---

# flowforge-explore

Investigate facts in the codebase to eliminate assumptions before planning or implementation.

## Core Rules

1. **Fact-Finding Only**: Never write or modify production code. Focus on discovering the current state.
2. **Precision & Evidence**: Every finding must cite exact file paths and line numbers (`path/to/file.ext:line`).
3. **Structured & Scaled Storage**:
   - **Brief facts (< 10 lines)**: Backfill directly into `01-workspace/<proposal_id>/README.md` under `## 3. Explored Facts`.
   - **Extensive research (web research, full schemas, third-party specs)**: Write to `01-workspace/<proposal_id>/references/<topic>.md`, and link it with a 1-line summary in `README.md`.

## Workflow

1. Identify the specific unknown or hypothesis (e.g. "How does the staging reader handle merged cells?").
2. Search and inspect the relevant code, configurations, or sample data.
3. If necessary, write a minimal read-only probe script to observe runtime behavior.
4. Record the verified fact and code seam into `01-workspace/<proposal_id>/README.md` under `## 3. Explored Facts`.
