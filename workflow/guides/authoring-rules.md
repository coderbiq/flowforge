# Authoring Rules

## General

- Prefer concise, decision-oriented writing over transcript dumps.
- Every durable statement should be attributable to a finding, decision, or explicit assumption.
- Use relative links between artifacts to keep references navigable in Git.

## Explorations

- `index.md` is the reading surface.
- `journal/` preserves chronology.
- `findings/` contains atomic statements worth reusing.
- `decisions/` contains candidate decisions with status.
- Use the archived knowledge base as the first source of truth for exploration research.
- Treat proposals and explorations as delta records against the existing final corpus, not as replacements for it.
- Do not mix implementation logs into exploration files.

## Proposals

- `meta.yaml` is the machine contract.
- `proposal.md` answers why and what.
- `design.md` answers how and why this approach.
- `task-map.md` is authoritative for task decomposition.
- `task-map.md` must follow `task-splitting.md` for deliverable-first decomposition, milestone boundaries, and checkpoint rules.
- `notes.md` is for execution history only.
- Proposals should begin by reviewing the canonical corpus and then describe the delta from that corpus.
- `canonical_corpus` in `meta.yaml` records which final docs were reviewed as the baseline for the proposal.
- Proposal creation may infer `canonical_corpus` from the declared archive targets and same-type final docs in the workspace; explicit overrides are allowed for broader baselines.
- Any manually supplied `canonical_corpus` entry must resolve to an existing document in the corresponding workspace.
- When a proposal changes existing final docs, describe the merge surface explicitly:
  - what section is updated in place
  - what facts are appended as new material
  - what facts are replaced or deprecated
  - what history note or changelog entry preserves the old fact
- For large modules, keep one canonical overview doc as the reader entry point and split dense details into subdocs such as `design.md`, `api.md`, `history.md`, or feature-specific pages.
- Do not duplicate the same fact across multiple final docs unless one of them is clearly marked as a pointer or historical record.
- If a proposal redistributes knowledge across docs, update the linked docs in one archive pass so the reader path stays coherent.

## Decisions

- Use ADRs only for stable, high-cost decisions with meaningful alternatives.
- Draft decisions belong in explorations or proposals until they are accepted.

## Archive targets

- Every proposal must declare at least one primary archive target.
- Primary target is where a future reader should start.
- Secondary targets exist to preserve alternate reading paths.
- The archived target corpus is the default baseline for future explorations, so archive targets should be kept navigable and up to date.
