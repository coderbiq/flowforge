---
name: flowforge-align
description: Use ONLY when discussing requirements, exploring new features, challenging assumptions, or resuming a multi-day proposal discussion in FlowForge. Engages in conversational grilling, boundary clarification, and maintains the proposal living scratchpad. Do NOT use for code implementation, test running, or final documentation curation.
---

# flowforge-align

Engage in conversational alignment (Grilling) to deeply understand requirements and capture working memory in the Proposal Scratchpad before planning or coding.

## Decision & Alignment Heuristics

### 1. The Design Tree & Decision Frontier (Grilling Protocol)
- Model the architecture as a decision tree. Ask questions along the active "decision frontier" in focused rounds (1-3 questions per round).
- Format each question with a clear recommendation and tradeoff:
  - `❓ Question`: The core architectural or boundary fork.
  - `➡️ Recommended Answer`: Concrete proposal with trade-offs ("We recommend X because Y, accepting tradeoff Z").

### 2. Domain Modeling & Active Challenge (DDD Protocol)
- **Glossary Conflicts**: When terms conflict with `docs/CONTEXT.md`, challenge immediately ("You said account, do you mean Customer or User?").
- **Edge-Case Stress Testing**: Probe domain boundaries with concrete scenarios ("If the target API fails midway during batch row 50, do we rollback the batch or resume?").
- **Code Contradiction Check**: Compare user statements with existing codebase facts. Surface contradictions before agreeing.

### 3. Adaptive Scale & Progressive Disclosure
- **Lean Mode (Default)**: For focused tasks (single domain, <= 5 slices), keep all consensus, terms, and slices inside `01-workspace/<proposal_id>/README.md`.
- **Hierarchical Mode**: If complexity escalates (>= 2 distinct sub-systems or extensive research), automatically split into `modules/<subsystem>.md` and `references/<topic>.md`, keeping `README.md` as the high-level architecture hub.

### 4. Continuous Scratchpad Sync
- When a Proposal is started, ensure `01-workspace/<proposal_id>/README.md` exists.
- Capture `[Ubiquitous Language]`, `[Key Decisions & Consensus]`, and `[Open Questions]` as they crystallize.

## Workflow

1. Check if `docs/CONTEXT.md` exists for project-level constraints and domain concepts.
2. If working on an existing proposal, read `01-workspace/<proposal_id>/README.md` to restore context. If new, create the Proposal Scratchpad.
3. Advance along the decision frontier using structured Grilling rounds.
4. Record confirmed decisions and domain terms directly into the Proposal Scratchpad.
5. When the user confirms readiness ("looks good", "let's plan"), transition to `flowforge-plan`.
