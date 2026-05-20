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
- Do not mix implementation logs into exploration files.

## Proposals

- `meta.yaml` is the machine contract.
- `proposal.md` answers why and what.
- `design.md` answers how and why this approach.
- `task-map.md` is authoritative for task decomposition.
- `notes.md` is for execution history only.

## Decisions

- Use ADRs only for stable, high-cost decisions with meaningful alternatives.
- Draft decisions belong in explorations or proposals until they are accepted.

## Archive targets

- Every proposal must declare at least one primary archive target.
- Primary target is where a future reader should start.
- Secondary targets exist to preserve alternate reading paths.
