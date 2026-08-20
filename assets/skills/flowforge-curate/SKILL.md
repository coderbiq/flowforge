---
name: flowforge-curate
description: Use ONLY when all implementation slices of a FlowForge proposal are completed and tested, and the proposal is ready to synthesize architectural decisions into ADRs, merge domain facts into Living Docs, and archive. Do NOT use during active coding or planning.
---

# flowforge-curate

Synthesize completed proposal outcomes, merge domain updates into Living Docs, record ADRs, and archive the proposal.

## Curation & Synthesis Heuristics

### 1. ADR 3-Criteria Filter
Only write an Architectural Decision Record (`docs/architecture/decisions/NNNN-xxx.md`) when a decision meets all 3 criteria:
1. **Hard to Reverse**: Significant cost to undo later.
2. **Surprising**: Not the obvious default; someone new to the codebase would ask "why was this done this way?".
3. **Real Trade-off**: Deliberately sacrificed one quality attribute (e.g. latency) to gain another (e.g. strict consistency).
*Do not record trivial implementation details as ADRs.*

### 2. Living Documentation 3-Way Patch
- Extract updated ubiquitous terms and operational realities from the delivered code into `docs/domains/<domain>/README.md`.
- Actively prune obsolete or superseded domain rules to prevent documentation decay.
- Discard transient debugging notes, trial-and-error logs, and temporary slice states.

### 3. Interactive Diff Review
- Always present a clean, concise markdown diff of the proposed Living Doc changes to the user before writing.

## Workflow

1. Read `01-workspace/<proposal_id>/README.md` and review the final state of code and tests.
2. Apply the ADR 3-Criteria filter to decide if a new ADR is required.
3. Identify the target domain file in `docs/domains/<domain>/README.md` and generate a patch reflecting current system truth.
4. Show the diff to the user for confirmation.
5. Apply updates to Living Docs, update `docs/CONTEXT.md`, and move the proposal directory to `archive/`.
